package ticker

import (
	"math"

	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/pkg/api"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/ma"
)

type ticker struct {
	logger         log.Logger
	data           contracts.TickerInfo
	emaService     ma.ExponentialMovingAverage
	averageService ma.MovingAverage
	emaPosNegSvc   ma.EMAPosNegService
}

func (s *ticker) Init(candles []contracts.Candle) (contracts.TickerInfo, error) {
	previousCP := 0.0
	nextCP := candles[0].CP
	for _, val := range candles {
		s.averageService.Add(constants.KeyCP10, val.CP)
		s.averageService.Add(constants.KeyCP50, val.CP)
		s.emaService.Add(constants.KeyCP5, val.CP)
		s.emaService.Add(constants.KeyCP20, val.CP)
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

	s.data.EmaCp5 = s.emaService.Value(constants.KeyCP5)
	s.data.EmaCP20 = s.emaService.Value(constants.KeyCP20)
	s.data.EmaDiffCpPos = s.emaPosNegSvc.Value(constants.KeyDiffCpPos)
	s.data.EmaDiffCpNeg = s.emaPosNegSvc.Value(constants.KeyDiffCpNeg)
	s.data.MinHP2 = math.Min(candles[lastIndex].HP, candles[lastIndex-1].HP)
	s.data.MinHP3 = math.Min(s.data.MinHP2, candles[lastIndex-2].HP)
	s.data.CP = candles[lastIndex].CP
	s.data.HP = candles[lastIndex].HP
	s.data.LP = candles[lastIndex].LP

	s.setNextValues(s.data.CP, s.data.HP, s.data.LP)
	return s.data, nil
}

func (s *ticker) Update(stock contracts.Stock) (contracts.TickerInfo, error) {
	s.setNextValues(stock.Close, stock.High, stock.Low)
	return s.data, nil
}

func (s *ticker) setNextValues(cp, hp, lp float64) {
	minHP2 := s.data.MinHP2

	s.data.Future.MinHP2 = math.Min(hp, minHP2)
	s.data.Future.MinHP3 = math.Min(hp, s.data.MinHP3)

	cpNext := []float64{cp, cp, s.data.Future.MinHP3}
	hpNext := []float64{hp, hp, s.data.Future.MinHP3}
	lpNext := []float64{lp, cp, s.data.Future.MinHP3}
	cpHpAv := []float64{s.data.HP, hp, hp}

	curCPDiff := cp - s.data.CP

	if curCPDiff >= 0 {
		s.data.Future.EmaDiffCpPos = (s.data.EmaDiffCpPos*13 + curCPDiff) / 14
		s.data.Future.EmaDiffCpNeg = (s.data.EmaDiffCpNeg*13 + 0) / 14
	} else {
		s.data.Future.EmaDiffCpPos = (s.data.EmaDiffCpPos*13 + 0) / 14
		s.data.Future.EmaDiffCpNeg = (s.data.EmaDiffCpNeg*13 + (-1 * curCPDiff)) / 14
	}

	s.data.Future.NextCP = cpNext
	s.data.Future.NextHP = hpNext
	s.data.Future.NextLP = lpNext
	s.data.AverageCp10 = s.averageService.Get(constants.KeyCP10)
	s.data.AverageCp50 = s.averageService.Get(constants.KeyCP50)
	s.data.Future.NextAvCPHP10 = s.averageService.AddArray(constants.KeyCP10, cpHpAv)
	s.data.Future.NextAvCPHP50 = s.averageService.AddArray(constants.KeyCP50, cpHpAv)
	s.data.MeanCp10 = s.averageService.AddArray(constants.KeyCP10, []float64{cp})
	s.data.MeanCp50 = s.averageService.AddArray(constants.KeyCP50, []float64{cp})
	s.data.Future.NextEMACP5 = s.emaService.AddArrayAndGet(constants.KeyCP5, cpNext)
	s.data.Future.NextEMACP20 = s.emaService.AddArrayAndGet(constants.KeyCP20, cpNext)
	s.data.Future.EmaCp5 = s.emaService.AddArray(constants.KeyCP5, cpNext)
	s.data.Future.EmaCP20 = s.emaService.AddArray(constants.KeyCP20, cpNext)
	s.data.Future.NextEMACPHP5 = s.emaService.AddArrayAndGet(constants.KeyCP5, cpHpAv)
	s.data.Future.NextEMACPHP20 = s.emaService.AddArrayAndGet(constants.KeyCP20, cpHpAv)

}

func NewTicker(logger log.Logger) api.TickerAPI {
	return &ticker{
		logger:         logger,
		data:           contracts.TickerInfo{},
		emaService:     ma.NewExponentialMovingAverage(logger, []string{"CP5", "CP20"}, []int{5, 20}),
		averageService: ma.NewMovingAverage(logger, []string{"CP10", "CP50"}, []int{10, 50}),
		emaPosNegSvc:   ma.NewEMAPosNeg(logger, []string{constants.KeyDiffCpPos, constants.KeyDiffCpNeg}, []int{15, 15}),
	}
}
