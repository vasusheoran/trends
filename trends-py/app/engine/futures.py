"""
Future projections — Support (CC) and Bullish (BR).

Both searches start from the settled EMA state (after all committed historical bars).
CE2 is the closed-form 2-bar fixed point: price at which BP stops changing after 2 bars.

Bullish (BR):
  d+1 = BR trial, d+2 = d+3 = CE2
  Find BR where BP[d+3] = BP[d+2]

Support (CC):
  cd3 = 2/6 * (trial - cd2) + cd2
  d+1 = cd3, d+2 = d+3 = trial
  Find trial where BP[d+3] = BP[d+2]
  Support = cd3

CD state: EMA-5 (decay 2/6) of daily CE2 values, accumulated from bar 100 onwards.

BP = EMA5 - EMA20  (EMA5 decay 2/6, EMA20 decay 2/21)
"""

from typing import Optional, Tuple
from scipy.optimize import brentq

from app.engine.indicators import EMAState, TOLERANCE


def _bp_series(ema5: EMAState, ema20: EMAState, closes: list[float]) -> list[Optional[float]]:
    """Apply closes sequentially to copied EMA state; return BP = EMA5 - EMA20 at each step."""
    m = ema5.copy()
    o = ema20.copy()
    results = []
    for c in closes:
        mv = m.update(c)
        ov = o.update(c)
        if mv is not None and ov is not None:
            results.append(mv - ov)
        else:
            results.append(None)
    return results


def _ce2(ema5: EMAState, ema20: EMAState) -> float:
    """Closed-form 2-bar fixed point: price where BP stops changing after 2 bars."""
    return (49 * ema5.value - 19 * ema20.value) / 30


def _search_br(trial: float, ema5: EMAState, ema20: EMAState, ce2: float) -> float:
    """
    W = [trial, CE2, CE2]; return BP[d+3] - BP[d+2].
    Find trial (BR/Bullish) where this is 0.
    """
    bp = _bp_series(ema5, ema20, [trial, ce2, ce2])
    if any(b is None for b in bp):
        return 0.0
    return bp[2] - bp[1]


def _search_cc(trial: float, ema5: EMAState, ema20: EMAState, cd2: float) -> float:
    """
    cd3 = 2/6*(trial - cd2) + cd2; W = [cd3, trial, trial].
    Return BP[d+3] - BP[d+2]; find trial where this is 0.
    Support = cd3 at the converged trial.
    """
    cd3 = (2 / 6) * (trial - cd2) + cd2
    bp = _bp_series(ema5, ema20, [cd3, trial, trial])
    if any(b is None for b in bp):
        return 0.0
    return bp[2] - bp[1]


def compute_futures(
    ema5_pre: EMAState,
    ema20_pre: EMAState,
    cd2: float,
    search_low: float = 0.0,
    search_high: float = 99999.0,
) -> Tuple[Optional[float], Optional[float]]:
    """
    Returns (support, bullish).

    ema5_pre / ema20_pre : settled EMA state (before current bar's close).
    cd2                  : CD EMA value after updating with today's CE2.
    """
    ce2 = _ce2(ema5_pre, ema20_pre)
    support = None
    bullish = None

    try:
        f_low  = _search_br(search_low,  ema5_pre, ema20_pre, ce2)
        f_high = _search_br(search_high, ema5_pre, ema20_pre, ce2)
        if f_low * f_high < 0:
            bullish = brentq(
                _search_br, search_low, search_high,
                args=(ema5_pre, ema20_pre, ce2),
                xtol=TOLERANCE, full_output=False,
            )
    except ValueError:
        pass

    try:
        f_low  = _search_cc(search_low,  ema5_pre, ema20_pre, cd2)
        f_high = _search_cc(search_high, ema5_pre, ema20_pre, cd2)
        if f_low * f_high < 0:
            cc_trial = brentq(
                _search_cc, search_low, search_high,
                args=(ema5_pre, ema20_pre, cd2),
                xtol=TOLERANCE, full_output=False,
            )
            support = (2 / 6) * (cc_trial - cd2) + cd2
    except ValueError:
        pass

    return support, bullish
