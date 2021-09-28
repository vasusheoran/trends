package cards

import (
	"math"

	"github.com/go-kit/kit/log"
	"github.com/vsheoran/trends/pkg/contracts"
)

type cards struct {
	logger log.Logger
}

func (r *cards) Get(ts contracts.TickerInfo) contracts.Card {
	futureEMA5 := futureEMA(2, 5, ts.Future.NextCP[2], ts.Future.NextEMACP5[2])
	futureEMA20 := futureEMA(2, 20, ts.Future.NextCP[2], ts.Future.NextEMACP20[2])
	hpEMA5 := futureEMA(1, 5, ts.Future.NextHP[0], ts.Future.NextEMACP5[2])
	hpEMA20 := futureEMA(1, 20, ts.Future.NextHP[0], ts.Future.NextEMACP20[2])

	var card contracts.Card

	card.CJ = cj(ts.Future.NextCP[2], futureEMA5, futureEMA20)
	card.U = u(card.CJ, ts.Future.MinHP2)
	card.BX = bx(ts.Future.NextHP[0], hpEMA5, hpEMA20)
	card.AI = ai(card.U, card.BX, ts.Future.NextHP[0], ts.HP, ts.Future.NextLP[0])
	card.AF = af(card.AI, ts.Future.NextLP[0])
	card.CO1 = co(ts.Future.NextAvCPHP10[1], ts.Future.NextAvCPHP50[1])
	card.CO0 = co(ts.Future.NextAvCPHP10[0], ts.Future.NextAvCPHP50[0])
	card.COLastDay = co(ts.AverageCp10, ts.AverageCp50)
	card.AE1 = ae(ts.Future.NextHP[0], card.CO1, card.CO0, ts.Future.NextEMACPHP5[1], ts.Future.NextEMACPHP5[0])
	card.AE2 = ae(ts.HP, card.CO0, card.COLastDay, ts.Future.NextEMACPHP5[0], ts.EmaCp5)
	card.AO = ao(card.AE1, card.AE2)
	card.BI = bi(card.AI, card.AF, card.AO, ts.Future.NextLP[0])
	card.BK = bk(ts.Future.NextHP[0], card.BI)
	card.BJ = bj(card.BI, card.BK)
	card.AR = ar(ts.MeanCp10[0], ts.MeanCp50[0])
	card.BN = bn(card.AR, ar(ts.AverageCp10, ts.AverageCp50), ts.Future.NextEMACP5[0], ts.EmaCp5)
	card.CR = cr(ts.Future.EmaDiffCpPos, ts.Future.EmaDiffCpNeg)
	card.MinHP3 = ts.MinHP3
	card.MinHP2 = ts.MinHP2
	card.MinHP = math.Min(ts.HP, ts.Future.NextHP[0])

	return card
}

func New(logger log.Logger) *cards {
	return &cards{
		logger: logger,
	}
}
