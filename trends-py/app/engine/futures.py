"""
Future projection fields — Support (CC) and Bullish.
Stateless: recomputed fresh every tick. No persistent state.

Computation chain:
  1. CE2  — binary search from pre-bar EMA state:
             find trial where BP[d+1] = BP[d+2], W[d+1]=W[d+2]=trial
  2. CD2  — EMA step: 2/6*(CE2 - old_cd) + old_cd   [persisted in TickerState.cd]
  3. CC   — binary search from pre-bar EMA state using CD2:
             find cc_trial where BP[d+2] = BP[d+3],
             W[d+1] = EMA_step(CD2, cc_trial), W[d+2]=W[d+3]=cc_trial
             → Support = CC = EMA_step(CD2, cc_trial)
  4. CE3  — binary search from POST-bar EMA state (after today's close applied):
             find trial where BP[d+1] = BP[d+2], W[d+1]=W[d+2]=trial
  5. CD3  — EMA step: 2/6*(CE3 - CD2) + CD2
  6. Bullish = CD3  (mathematical fixed point of the W4=W3=CD4 equilibrium)

BP = EMA5 - EMA20  (source: Final-bullish-ce.xlsx col BP = AS - BN, decay 2/21)
Pre-bar EMA state = state BEFORE the current bar's close was applied.
Post-bar EMA state = state AFTER the current bar's close was applied.
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


def compute_futures(
    ema5_pre: EMAState,
    ema20_pre: EMAState,
    ema5_post: EMAState,
    ema20_post: EMAState,
    old_cd: float,
    search_low: float = 0.0,
    search_high: float = 99999.0,
) -> Tuple[Optional[float], Optional[float], float]:
    """
    Returns (support=CC, bullish=CD3, new_cd=CD2).

    ema5_pre/ema20_pre: EMA state BEFORE today's close was applied (for CE2 and CC).
    ema5_post/ema20_post: EMA state AFTER today's close was applied (for CE3).
    old_cd: CD from previous tick; updated to CD2 (returned as new_cd).

    Uses copies internally — originals are never modified.
    """
    support = None
    bullish = None

    # Step 1: CE2 — find trial where BP[d+1] = BP[d+2], using pre-bar EMA state
    ce2 = None
    try:
        f_low = _search_ce(search_low, ema5_pre, ema20_pre)
        f_high = _search_ce(search_high, ema5_pre, ema20_pre)
        if f_low * f_high < 0:
            ce2 = brentq(
                _search_ce,
                search_low, search_high,
                args=(ema5_pre, ema20_pre),
                xtol=TOLERANCE,
                full_output=False,
            )
    except ValueError:
        pass

    if ce2 is None:
        return None, None, old_cd

    # Step 2: CD2 — rolling EMA of CE values, persisted as new_cd
    cd2 = 2 / 6 * (ce2 - old_cd) + old_cd

    # Step 3: CC (Support) — binary search using pre-bar EMA state and CD2
    try:
        f_low = _search_cc(search_low, ema5_pre, ema20_pre, cd2)
        f_high = _search_cc(search_high, ema5_pre, ema20_pre, cd2)
        if f_low * f_high < 0:
            cc_trial = brentq(
                _search_cc,
                search_low, search_high,
                args=(ema5_pre, ema20_pre, cd2),
                xtol=TOLERANCE,
                full_output=False,
            )
            support = 2 / 6 * (cc_trial - cd2) + cd2
    except ValueError:
        pass

    # Step 4: CE3 — same binary search but from POST-bar EMA state
    ce3 = None
    try:
        f_low = _search_ce(search_low, ema5_post, ema20_post)
        f_high = _search_ce(search_high, ema5_post, ema20_post)
        if f_low * f_high < 0:
            ce3 = brentq(
                _search_ce,
                search_low, search_high,
                args=(ema5_post, ema20_post),
                xtol=TOLERANCE,
                full_output=False,
            )
    except ValueError:
        pass

    # Step 5: CD3 = EMA_step(CD2, CE3); Bullish = CD3
    if ce3 is not None:
        bullish = 2 / 6 * (ce3 - cd2) + cd2

    return support, bullish, cd2
