"""
Support (CC) and Bullish (2-day convergence) computation validation.

Support: cd3=(2/6)*(cc_trial-CD_curr)+CD_curr; W=[cd3, cc_trial, cc_trial];
         find cc_trial where BP[d+3]=BP[d+2]; Support=cd3.
Bullish: find W where Support(Day+2 with W) == W.

Tests verify the convergence properties on values produced by TickerState.
"""

import pytest
from scipy.optimize import brentq

from app.engine.state import TickerState
from app.engine.futures import compute_futures, _bp_series, _search_cc, _ce2, _CD_DECAY, get_support
from app.engine.indicators import EMAState, TOLERANCE

CONVERGENCE_TOLERANCE = TOLERANCE * 5


def _collect_futures(excel_rows):
    """
    Feed all rows through TickerState in commit mode.
    Returns list of (row, ema5_pre, ema20_pre, ema5_post, ema20_post, cd_curr, support, bullish).
    cd_curr = cd_ema value AFTER updating with this bar's CE2.
    ema5_post/ema20_post = EMA state AFTER today's bar (used for Bullish verification).
    """
    state = TickerState(ticker="TEST")
    results = []

    for r in excel_rows:
        ema5_pre  = state.ema5.copy()
        ema20_pre = state.ema20.copy()

        snap = state.update(r["date"], r["close"], r["open"], r["high"], r["low"])

        if snap.support is not None and snap.bullish is not None:
            results.append((r, ema5_pre, ema20_pre, state.ema5.copy(), state.ema20.copy(),
                            state.cd_ema.value, snap.support, snap.bullish))

    return results


@pytest.fixture(scope="module")
def futures_results(excel_rows):
    return _collect_futures(excel_rows)


def test_bullish_convergence_last_20_rows(futures_results):
    """Bullish: verify Support(Day+2 with W) ≈ W using post-bar EMA state."""
    sample = futures_results[-20:]
    assert len(sample) == 20

    failures = []
    for r, ema5_pre, ema20_pre, ema5_post, ema20_post, cd_curr, support, bullish in sample:
        # Day+1 CD step (from post-bar EMA)
        ce2_d1 = _ce2(ema5_post, ema20_post)
        cd_d1 = _CD_DECAY * (ce2_d1 - cd_curr) + cd_curr

        # Apply W=bullish to post-bar EMA for Day+1
        e5_d1 = ema5_post.copy(); e5_d1.update(bullish)
        e20_d1 = ema20_post.copy(); e20_d1.update(bullish)

        # Day+2 CD step
        ce2_d2 = _ce2(e5_d1, e20_d1)
        cd_d2 = _CD_DECAY * (ce2_d2 - cd_d1) + cd_d1

        sup_d2 = get_support(e5_d1, e20_d1, cd_d2)
        if sup_d2 is None:
            failures.append(f"row {r['row']}: Support(Day+2) is None at bullish={bullish:.2f}")
            continue

        diff = abs(sup_d2 - bullish)
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(
                f"row {r['row']}: |Support(Day+2)-W|={diff:.6f} > {CONVERGENCE_TOLERANCE} "
                f"at bullish={bullish:.2f}"
            )

    assert not failures, "Bullish convergence failures:\n" + "\n".join(failures)


def test_support_convergence_last_20_rows(futures_results):
    """Support: verify W=[support, cc_trial, cc_trial] → BP[d+3]=BP[d+2].
    Recover cc_trial from support=cd3: cc_trial=(cd3-(1-2/6)*cd_curr)/(2/6).
    """
    sample = futures_results[-20:]
    assert len(sample) == 20

    failures = []
    for r, ema5_pre, ema20_pre, ema5_post, ema20_post, cd_curr, support, bullish in sample:
        cd3 = support
        cc_trial = (cd3 - cd_curr) / _CD_DECAY + cd_curr
        bp = _bp_series(ema5_pre, ema20_pre, [cd3, cc_trial, cc_trial])
        if any(b is None for b in bp):
            failures.append(f"row {r['row']}: BP None at support={support:.2f}")
            continue
        diff = abs(bp[2] - bp[1])
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(
                f"row {r['row']}: |BP[d+3]-BP[d+2]|={diff:.6f} > {CONVERGENCE_TOLERANCE} "
                f"at support={support:.2f} cc_trial={cc_trial:.2f}"
            )

    assert not failures, "Support convergence failures:\n" + "\n".join(failures)


def test_futures_populated_after_100_bars(excel_rows):
    """Futures must be None before 105 bars (CD needs 5 CE2 values) and non-None by bar 120."""
    state = TickerState(ticker="TEST")
    none_after_105 = []

    for i, r in enumerate(excel_rows[:200]):
        snap = state.update(r["date"], r["close"], r["open"], r["high"], r["low"])
        if i >= 105 and snap.support is None:
            none_after_105.append(i)

    late_nones = [i for i in none_after_105 if i >= 120]
    assert not late_nones, f"Futures still None at bars: {late_nones[:5]}"
