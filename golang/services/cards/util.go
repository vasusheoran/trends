package cards

import (
	"math"
)

// CP_CI_CH -> copy of CP except 0 and 1
// diffCP1 -> emaService rolling updated with span:2 for cp difference of current - previous
// emaDiffCP1Pos -> df.at[index, col_name] = (df.at[(index-1), col_name] *13 + df.at[index, col_name])/14
// emaDiffCP1Neg -> df.at[index, col_name] = (df.at[(index-1), col_name] *13 + (-1 * df.at[index, col_name]))/14

//span 5 prevema 15149.292605464798 lastrowcp 15188.5
//span 5 prevema 15171.074491317688
//span 20 prevema 14801.795416480334 lastrowcp 15188.5
//span 20 prevema 14871.94590782177
func futureEMA(iteration, span int, val, result float64) float64 {
	for i := 0; i < iteration; i++ {
		result = ((2 / (float64(span) + 1)) * (val - result)) + result
	}
	return result
}
func ema(span, val, prev_ema float64) float64 {
	return ((2 / (span + 1)) * (val - prev_ema)) + prev_ema
}

func bx(hp, def, ghi float64) float64 {
	return (hp + hp + (((((def) + ((def) + ((ghi)-(def))/2)) / 2) + ((hp + (hp + ((((def)+((def)+((ghi)-(def))/2))/2)-hp)/2)) / 2)) / 2)) / 3
}

//future_ema5 15171.074491317688 future_ema20 14871.94590782177 cpAt0 15188.5 cj 15169.290071967438 cj2 15169.290071967438
func cj(cpAt0, futureEMA5, futureEMA20 float64) float64 {
	return (cpAt0 + cpAt0 + ((((futureEMA5 + (futureEMA5 + (futureEMA20-futureEMA5)/2)) / 2) + ((cpAt0 + (cpAt0 + (((futureEMA5+(futureEMA5+(futureEMA20-futureEMA5)/2))/2)-cpAt0)/2)) / 2)) / 2)) / 3
}

func u(cjVal, minHP3At0 float64) float64 {
	return (2 * minHP3At0) - cjVal
}

func co(avCPCIHP10AtRow, avCPCIHP50AtRow float64) float64 {
	cvAverage := avCPCIHP10AtRow + avCPCIHP50AtRow
	return (cvAverage)/2 - ((cvAverage) / 2 * (((cvAverage)/2 - (((((cvAverage)/2 - ((cvAverage) / 2 * 0.01)) + (((cvAverage)/2 - ((cvAverage) / 2 * 0.01)) * 0.025)) + (cvAverage)/2) / 2)) / (cvAverage) / 2 * 100 / 2) / 100)
}

func ae(hpAtRowPlus1, coAtRow, coAtRowPlus1, emaCPCIHP5AtRow, emaCPCIHP5AtRowPlus1 float64) float64 {
	return hpAtRowPlus1 - ((emaCPCIHP5AtRow - coAtRow) - (emaCPCIHP5AtRowPlus1 - coAtRowPlus1))
}

func ai(uVal, bxVal, hpAt2, hpAt3, lpAt2 float64) float64 {
	return (lpAt2 + (hpAt2+(uVal-hpAt2)/2+math.Min(hpAt3, hpAt2)+((hpAt2*2-bxVal)-math.Min(hpAt3, hpAt2))/2)/2) / 2
}

func af(ai, lpAt2 float64) float64 {
	return lpAt2 + (ai-lpAt2)/2
}

func ao(aeAtRow, aeAtRowPlus1 float64) float64 {
	return aeAtRow - aeAtRowPlus1
}

func bk(hpAt2, prevBI float64) float64 {
	return prevBI + (hpAt2-prevBI)/2 + (hpAt2-(prevBI+(hpAt2-prevBI)/2))/2
}

func br(brsh, bj float64) float64 {
	return brsh - (brsh-bj)/2
}

func barish(lpAt0, br float64) float64 {
	return lpAt0 + ((br - lpAt0) / 2)
}

func bj(curBI, bk float64) float64 {
	return (curBI + bk) / 2
}

func ar(cpMeanRowSpan10, cpMeanRowSpan50 float64) float64 {
	return ((cpMeanRowSpan10 + cpMeanRowSpan50) / 2) - ((cpMeanRowSpan10 + cpMeanRowSpan50) / 2 * (((cpMeanRowSpan10+cpMeanRowSpan50)/2 - (((((cpMeanRowSpan10+cpMeanRowSpan50)/2 - ((cpMeanRowSpan10 + cpMeanRowSpan50) / 2 * 0.01)) + (((cpMeanRowSpan10+cpMeanRowSpan50)/2 - ((cpMeanRowSpan10 + cpMeanRowSpan50) / 2 * 0.01)) * 0.025)) + (cpMeanRowSpan10+cpMeanRowSpan50)/2) / 2)) / (cpMeanRowSpan10 + cpMeanRowSpan50) / 2 * 100 / 2) / 100)
}

//df.at[2, 'ema_diffCP1Pos'] 56.192976714563166 df.at[2, 'ema_diffCP1Neg'] 272.592886852383
func cr(emaDiffCP1PosAt2, emaDiffCP1NegAt2 float64) float64 {
	if emaDiffCP1NegAt2 == 0 {
		return 100.0
	}
	return 100 - (100 / (1 + (emaDiffCP1PosAt2)/emaDiffCP1NegAt2))
}

//bn -985.0807869304917 ar 17022.095212980832 ar3 17202.227584425626 df.at[2,'emaCP5'] 16584.42631675057 df.at[3,'emaCP5'] 17749.639475125856
func bn(arAtRow, arAtRowPlus1, emaCP5AtRow, emaCP5AtRowPlus1 float64) float64 {
	return (emaCP5AtRow - arAtRow) - (emaCP5AtRowPlus1 - arAtRowPlus1)
}

func bi(ai, af, ao, lpAt2 float64) float64 {
	afaiDiff := af - ai
	return math.Max((ai+(af+(afaiDiff)/2))/2, ((lpAt2-ao)+(af+(afaiDiff)/2))/2)
}

func trend(cpAt0, bn float64) float64 {
	return cpAt0 - bn
}
