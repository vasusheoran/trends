"""
Future projections — Support (CC) and Bullish.

Support (CC):
  Binary search for cc_trial where BP[d+3] == BP[d+2] with W = [cd3, cc_trial, cc_trial]
  where cd3 = (2/6)*(cc_trial - CD_curr) + CD_curr.
  Support = cd3 (W[d+1] in Go's searchCC).
  CD_curr = EMA_step(CD_prev, CE2) using decay 2/6.

Bullish:
  The price W such that applying W for Day+1 and Day+2 makes Support(Day+2) == W.
  Binary search over W using the 2-day convergence algorithm.

CE2 = (49*EMA5 - 19*EMA20) / 30  — closed-form 2-bar fixed point.
CD  = EMA-5 (decay 2/6) of daily CE2 values.
BP  = EMA5 - EMA20.
"""

from typing import Optional, Tuple
from scipy.optimize import brentq

from app.engine.indicators import EMAState, TOLERANCE

_CD_DECAY = 2 / 6


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


def _search_cc(trial: float, ema5_pre: EMAState, ema20_pre: EMAState, cd_curr: float) -> float:
    """
    W = [cd3, trial, trial] where cd3 = (2/6)*(trial - CD_curr) + CD_curr.
    Return BP[d+3] - BP[d+2]; root gives the converged cc_trial.
    Support = cd3 at the converged trial.
    """
    cd3 = _CD_DECAY * (trial - cd_curr) + cd_curr
    bp = _bp_series(ema5_pre, ema20_pre, [cd3, trial, trial])
    if any(b is None for b in bp):
        return 0.0
    return bp[2] - bp[1]


def get_support(
    ema5_pre: EMAState,
    ema20_pre: EMAState,
    cd_curr: float,
    search_low: float = 0.0,
    search_high: float = 99999.0,
) -> Optional[float]:
    """Support (CC) for a bar: cd3 at the converged cc_trial."""
    try:
        f_low = _search_cc(search_low, ema5_pre, ema20_pre, cd_curr)
        f_high = _search_cc(search_high, ema5_pre, ema20_pre, cd_curr)
        if f_low * f_high < 0:
            cc_trial = brentq(
                _search_cc, search_low, search_high,
                args=(ema5_pre, ema20_pre, cd_curr),
                xtol=TOLERANCE, full_output=False,
            )
            return _CD_DECAY * (cc_trial - cd_curr) + cd_curr
    except (ValueError, TypeError):
        pass
    return None


def _search_bullish(
    w_trial: float,
    ema5_pre: EMAState,
    ema20_pre: EMAState,
    cd_pre: float,
) -> float:
    """
    2-day convergence: find W where Support(Day+2) == W.

    Day+1: apply W, update EMA and CD.
    Day+2: apply W again, update CD from Day+1 state.
    Compute Support for Day+2 using Day+1 EMA state and Day+2 CD.
    Return Support - W.
    """
    # Day +1
    ce2_d1 = _ce2(ema5_pre, ema20_pre)
    cd_d1 = _CD_DECAY * (ce2_d1 - cd_pre) + cd_pre

    e5_d1 = ema5_pre.copy()
    e5_d1.update(w_trial)
    e20_d1 = ema20_pre.copy()
    e20_d1.update(w_trial)

    # Day +2
    ce2_d2 = _ce2(e5_d1, e20_d1)
    cd_d2 = _CD_DECAY * (ce2_d2 - cd_d1) + cd_d1

    support_d2 = get_support(e5_d1, e20_d1, cd_d2)
    if support_d2 is None:
        return 99999.0 - w_trial
    return support_d2 - w_trial


def _search_hold(
    trial: float,
    ema5_bull: EMAState,
    ema20_bull: EMAState,
    cd_for_bull: float,
    bullish_target: float,
) -> float:
    """
    Find trial (D+1 close) where Bullish from D+1 state == bullish_target.

    Applies trial for D+1, then delegates to _search_bullish for D+2 onward.
    Root: Support(D+3) - bullish_target == 0 when D+2 applies bullish_target.
    """
    ce2_d1 = _ce2(ema5_bull, ema20_bull)
    cd_d1 = _CD_DECAY * (ce2_d1 - cd_for_bull) + cd_for_bull

    e5_d1 = ema5_bull.copy()
    e5_d1.update(trial)
    e20_d1 = ema20_bull.copy()
    e20_d1.update(trial)

    return _search_bullish(bullish_target, e5_d1, e20_d1, cd_d1)


def compute_futures(
    ema5_pre: EMAState,
    ema20_pre: EMAState,
    cd_pre: float,
    search_low: float = 0.0,
    search_high: float = 99999.0,
) -> Tuple[Optional[float], Optional[float], Optional[float]]:
    """
    Returns (support, bullish, hold).

    All three use ema5_pre/ema20_pre — the EMA state BEFORE today's bar.
    cd_pre: CD EMA value BEFORE updating with today's CE2.

    Hold: the minimum D+1 close such that Bullish from D+1 state == today's Bullish.
    """
    # CD updated with today's CE2 (using pre-bar EMA)
    ce2_today = _ce2(ema5_pre, ema20_pre)
    cd_curr = _CD_DECAY * (ce2_today - cd_pre) + cd_pre

    support = get_support(ema5_pre, ema20_pre, cd_curr, search_low, search_high)

    bullish = None
    try:
        f_low = _search_bullish(search_low, ema5_pre, ema20_pre, cd_pre)
        f_high = _search_bullish(search_high, ema5_pre, ema20_pre, cd_pre)
        if f_low * f_high < 0:
            bullish = brentq(
                _search_bullish, search_low, search_high,
                args=(ema5_pre, ema20_pre, cd_pre),
                xtol=TOLERANCE, full_output=False,
            )
    except (ValueError, TypeError):
        pass

    hold = None
    if bullish is not None:
        try:
            f_low = _search_hold(search_low, ema5_pre, ema20_pre, cd_pre, bullish)
            f_high = _search_hold(search_high, ema5_pre, ema20_pre, cd_pre, bullish)
            if f_low * f_high < 0:
                hold = brentq(
                    _search_hold, search_low, search_high,
                    args=(ema5_pre, ema20_pre, cd_pre, bullish),
                    xtol=TOLERANCE, full_output=False,
                )
        except (ValueError, TypeError):
            pass

    return support, bullish, hold
