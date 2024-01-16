package ticker

import (
	"errors"
	"fmt"
	"github.com/vsheoran/trends/services/cards"
	"github.com/vsheoran/trends/services/history"
	"math"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ma"
)

type Ticker interface {
	Init(key string) (contracts.Summary, error)
	Update(key string, stock contracts.Stock) error
	Get(key string) (contracts.Summary, error)
	Freeze(key string, st contracts.Stock) error
}

type ticker struct {
	logger         log.Logger
	data           map[string]*contracts.TickerInfo
	summary        map[string]*contracts.Summary
	cardService    cards.Cards
	historyService history.History
	emaService     ma.ExponentialMovingAverage
	averageService ma.MovingAverage
	emaPosNegSvc   ma.EMAPosNegService
}

func (s *ticker) Freeze(key string, st contracts.Stock) error {
	if _, ok := s.data[key]; !ok {
		msg := fmt.Sprintf("ticker '%s' has not been initialized", key)
		return errors.New(msg)
	}

	err := s.Update(key, st)
	if err != nil {
		return err
	}

	sum, err := s.Get(key)
	if err != nil {
		return err
	}

	s.data[key].IsBuyFrozen = true
	s.data[key].BI = sum.Buy

	return nil
}

func (s *ticker) Update(key string, stock contracts.Stock) error {
	if _, ok := s.data[key]; !ok {
		msg := fmt.Sprintf("ticker '%s' has not been initialized", key)
		return errors.New(msg)
	}
	s.setNextValues(key, stock.CP, stock.HP, stock.LP)

	card := s.cardService.Get(*s.data[key])
	s.updateSummaryMap(key, card)

	return nil
}

func (s *ticker) Get(key string) (contracts.Summary, error) {
	var summary *contracts.Summary
	var ok bool
	if _, ok = s.data[key]; !ok {
		msg := fmt.Sprintf("ticker '%s' has not been initialized", key)
		return contracts.Summary{}, errors.New(msg)
	}
	if summary, ok = s.summary[key]; !ok {
		card := s.cardService.Get(*s.data[key])
		s.updateSummaryMap(key, card)
		summary = s.summary[key]
	}

	return *summary, nil
}

func (s *ticker) Init(key string) (contracts.Summary, error) {
	candles, err := s.historyService.Read(key)
	if err != nil {
		level.Error(s.logger).Log("failed to read history from database", "sasSymbol", key)
		return contracts.Summary{}, errors.New("error fetching data")
	}
	if candles == nil {
		level.Error(s.logger).Log("failed to parse history from database", "sasSymbol", key)
		return contracts.Summary{}, errors.New("error fetching data")
	}

	var st *contracts.TickerInfo
	var ok bool
	if st, ok = s.data[key]; !ok {
		st = &contracts.TickerInfo{}
		s.data[key] = st
	}

	previousCP := 0.0
	nextCP := candles[0].CP
	for _, val := range candles {
		s.averageService.Add(key, constants.KeyCP10, val.CP)
		s.averageService.Add(key, constants.KeyCP50, val.CP)
		s.emaService.Add(key, constants.KeyCP5, val.CP)
		s.emaService.Add(key, constants.KeyCP20, val.CP)
		previousCP = nextCP
		nextCP = val.CP
		diff := nextCP - previousCP

		if diff >= 0.0 {
			s.emaPosNegSvc.Add(constants.KeyDiffCpNeg, 0.0)
			s.emaPosNegSvc.Add(constants.KeyDiffCpPos, diff)
		} else {
			s.emaPosNegSvc.Add(constants.KeyDiffCpNeg, (-1)*diff)
			s.emaPosNegSvc.Add(constants.KeyDiffCpPos, 0.0)
		}
	}

	lastIndex := len(candles) - 1

	s.data[key].EmaCp5 = s.emaService.Value(key, constants.KeyCP5)
	s.data[key].EmaCP20 = s.emaService.Value(key, constants.KeyCP20)
	s.data[key].EmaDiffCpPos = s.emaPosNegSvc.Value(constants.KeyDiffCpPos)
	s.data[key].EmaDiffCpNeg = s.emaPosNegSvc.Value(constants.KeyDiffCpNeg)
	s.data[key].MinHP2 = math.Min(candles[lastIndex].HP, candles[lastIndex-1].HP)
	s.data[key].MinHP3 = math.Min(s.data[key].MinHP2, candles[lastIndex-2].HP)
	s.data[key].MinLP2 = math.Min(candles[lastIndex].LP, candles[lastIndex-1].LP)
	s.data[key].MinLP3 = math.Min(s.data[key].MinLP2, candles[lastIndex-2].LP)
	s.data[key].LowerL = candles[lastIndex].LP
	s.data[key].CP = candles[lastIndex].CP
	s.data[key].HP = candles[lastIndex].HP
	s.data[key].LP = candles[lastIndex].LP

	s.setNextValues(key, s.data[key].CP, s.data[key].HP, s.data[key].LP)
	return s.Get(key)
}

func (s *ticker) setNextValues(key string, cp, hp, lp float64) {
	minHP2 := s.data[key].MinHP2

	s.data[key].Future.MinHP2 = math.Min(hp, minHP2)
	s.data[key].Future.MinHP4 = math.Min(hp, s.data[key].MinHP3)
	s.data[key].Future.MinLP3 = math.Min(lp, s.data[key].MinLP2)

	cpNext := []float64{cp, cp, s.data[key].Future.MinHP4}
	hpNext := []float64{hp, hp, s.data[key].Future.MinHP4}
	lpNext := []float64{lp, cp, s.data[key].Future.MinHP4}
	cpHpAv := []float64{s.data[key].HP, hp, hp}

	curCPDiff := cp - s.data[key].CP

	if curCPDiff >= 0 {
		s.data[key].Future.EmaDiffCpPos = (s.data[key].EmaDiffCpPos*13 + curCPDiff) / 14
		s.data[key].Future.EmaDiffCpNeg = (s.data[key].EmaDiffCpNeg*13 + 0) / 14
	} else {
		s.data[key].Future.EmaDiffCpPos = (s.data[key].EmaDiffCpPos*13 + 0) / 14
		s.data[key].Future.EmaDiffCpNeg = (s.data[key].EmaDiffCpNeg*13 + (-1 * curCPDiff)) / 14
	}

	s.data[key].Future.NextCP = cpNext
	s.data[key].Future.NextHP = hpNext
	s.data[key].Future.NextLP = lpNext
	s.data[key].AverageCp10 = s.averageService.Value(key, constants.KeyCP10)
	s.data[key].AverageCp50 = s.averageService.Value(key, constants.KeyCP50)
	s.data[key].Future.NextAvCPHP10 = s.averageService.AddArray(key, constants.KeyCP10, cpHpAv)
	s.data[key].Future.NextAvCPHP50 = s.averageService.AddArray(key, constants.KeyCP50, cpHpAv)
	s.data[key].MeanCp10 = s.averageService.AddArray(key, constants.KeyCP10, []float64{cp})
	s.data[key].MeanCp50 = s.averageService.AddArray(key, constants.KeyCP50, []float64{cp})
	s.data[key].Future.NextEMACP5 = s.emaService.AddArray(key, constants.KeyCP5, cpNext)
	s.data[key].Future.NextEMACP20 = s.emaService.AddArray(key, constants.KeyCP20, cpNext)
	s.data[key].Future.NextEMACPHP5 = s.emaService.AddArray(key, constants.KeyCP5, cpHpAv)
	s.data[key].Future.NextEMACPHP20 = s.emaService.AddArray(key, constants.KeyCP20, cpHpAv)

}

func (s *ticker) updateSummaryMap(key string, card contracts.Card) {
	s.summary[key] = &contracts.Summary{
		Close:       s.data[key].Future.NextCP[0],
		High:        s.data[key].Future.NextHP[0],
		Low:         s.data[key].Future.NextLP[0],
		Average:     card.AR,
		Ema20:       s.data[key].Future.NextEMACP20[0],
		MinLP3:      s.data[key].Future.MinLP3,
		LowerL:      s.data[key].LowerL,
		Ema5:        s.data[key].Future.NextEMACP5[0],
		RSI:         card.CR,
		HL3:         s.data[key].MinHP3,
		Barish:      card.Barish,
		Support:     card.BJ,
		Trend:       card.Trend,
		Bullish:     card.BK,
		Buy:         card.BI,
		PreviousBuy: card.PreviousBI,
	}
}

func NewTicker(logger log.Logger, cardsSvc cards.Cards, hs history.History) Ticker {
	return &ticker{
		logger:         logger,
		data:           map[string]*contracts.TickerInfo{},
		summary:        map[string]*contracts.Summary{},
		cardService:    cardsSvc,
		historyService: hs,
		emaService:     ma.NewExponentialMovingAverage(logger, []string{"CP5", "CP20"}, []int{5, 20}),
		averageService: ma.NewMovingAverage(logger, []string{"CP10", "CP50"}, []int{10, 50}),
		emaPosNegSvc:   ma.NewEMAPosNeg(logger, []string{constants.KeyDiffCpPos, constants.KeyDiffCpNeg}, []int{15, 15}),
	}
}
