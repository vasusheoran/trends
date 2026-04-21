"""
Future projection fields — Support (CC) and Bullish (BR).
Stateless: recomputed fresh every tick. No persistent state.

Computation chain (matches Go updateFutureData):
  1. CE   — binary search: find trial where BP[d+1] = BP[d+2], W[d+1]=W[d+2]=trial
  2. CD   — rolling EMA5 of CE values: CD = 2/6*(CE - CD_prev) + CD_prev  (passed in as state)
  3. CC   — binary search: find cc_trial where BP[d+1] = BP[d+2],
              W[d+1] = EMA_step(CD, cc_trial), W[d+2]=W[d+3]=cc_trial
             → Support = CC = EMA_step(CD, cc_trial)  (= W[d+1], not the trial itself)
  4. BR   — binary search: find trial where BP[d+3] = BP[d+2],
              W[d+1]=trial, W[d+2]=W[d+3]=CE

BP = EMA5 - EMA20  (source: Final-bullish-ce.xlsx col BP = AS - BN, decay 2/21)
EMA state must be from BEFORE the current bar's close was applied.
"""

from typing import Optional, Tuple
from scipy.optimize import brentq

from app.engine.indicators import EMAState, TOLERANCE


def _bp_series(ema5: EMAState, ema20: EMAState, closes: list[float]) -> list[Optional[float]]:
    """Apply closes to copied EMA state, return BP = EMA5 - EMA20 at each step."""
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


def _search_ce(trial: float, ema5: EMAState, ema20: EMAState) -> float:
    """Return BP[d+2] - BP[d+1]; find trial where this is 0."""
    bp = _bp_series(ema5, ema20, [trial, trial])
    if bp[0] is None or bp[1] is None:
        return 0.0
    return bp[1] - bp[0]


def _search_cc(trial: float, ema5: EMAState, ema20: EMAState, cd: float) -> float:
    """
    Return BP[d+3] - BP[d+2]; find trial where this is 0.
    W[d+1] = EMA_step(CD, trial), W[d+2] = W[d+3] = trial.
    CC (Support) = W[d+1] = 2/6*(trial-CD)+CD, not the trial itself.
    """
    w_d1 = 2 / 6 * (trial - cd) + cd
    bp = _bp_series(ema5, ema20, [w_d1, trial, trial])
    if any(b is None for b in bp):
        return 0.0
    return bp[2] - bp[1]


def _search_br(trial: float, ema5: EMAState, ema20: EMAState, ce: float) -> float:
    """Return BP[d+3] - BP[d+2]; find trial where this is 0. W[d+1]=trial, W[d+2,3]=CE."""
    bp = _bp_series(ema5, ema20, [trial, ce, ce])
    if any(b is None for b in bp):
        return 0.0
    return bp[2] - bp[1]


def compute_futures(
    ema5: EMAState,
    ema20: EMAState,
    old_cd: float,
    search_low: float = 0.0,
    search_high: float = 99999.0,
) -> Tuple[Optional[float], Optional[float], float]:
    """
    Returns (support=CC, bullish=BR, new_cd) via binary search.
    ema5/ema20 must be state BEFORE the current bar's close was applied.
    old_cd is the CD value from the previous tick; this function updates it with today's CE.
    Uses copies internally — originals are never modified.
    """
    ce = None
    support = None
    bullish = None

    # Step 1: CE — find trial where BP[d+1] = BP[d+2]
    try:
        f_low = _search_ce(search_low, ema5, ema20)
        f_high = _search_ce(search_high, ema5, ema20)
        if f_low * f_high < 0:
            ce = brentq(
                _search_ce,
                search_low, search_high,
                args=(ema5, ema20),
                xtol=TOLERANCE,
                full_output=False,
            )
    except ValueError:
        pass

    if ce is None:
        return None, None, old_cd

    # Update CD = rolling EMA5 of CE values (matches Go calculateCD)
    new_cd = 2 / 6 * (ce - old_cd) + old_cd

    # Step 2: CC (Support) — find cc_trial where BP[d+1]=BP[d+2] with W[d+1]=EMA(new_cd,trial)
    try:
        f_low = _search_cc(search_low, ema5, ema20, new_cd)
        f_high = _search_cc(search_high, ema5, ema20, new_cd)
        if f_low * f_high < 0:
            cc_trial = brentq(
                _search_cc,
                search_low, search_high,
                args=(ema5, ema20, new_cd),
                xtol=TOLERANCE,
                full_output=False,
            )
            support = 2 / 6 * (cc_trial - new_cd) + new_cd  # CC = W[d+1] = EMA_step(CD, trial)
    except ValueError:
        pass

    # Step 3: BR (Bullish) — find trial where BP[d+3]=BP[d+2] with day+1=trial, day+2,3=CE
    try:
        f_low = _search_br(search_low, ema5, ema20, ce)
        f_high = _search_br(search_high, ema5, ema20, ce)
        if f_low * f_high < 0:
            bullish = brentq(
                _search_br,
                search_low, search_high,
                args=(ema5, ema20, ce),
                xtol=TOLERANCE,
                full_output=False,
            )
    except ValueError:
        pass

    return support, bullish, new_cd
