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
    Feed all rows through TickerState.
    Returns list of (row, ema5_pre, ema20_pre, cd_pre, support, bullish).
    ema5_pre/ema20_pre/cd_pre: EMA state BEFORE today's bar (used for all futures).
    """
    state = TickerState(ticker="TEST")
    results = []

    for r in excel_rows:
        snap = state.update(r["date"], r["close"], r["open"], r["high"], r["low"])

        if snap.support is not None and snap.bullish is not None:
            ps = state._pre_bar_state
            results.append((r, ps.ema5.copy(), ps.ema20.copy(), ps.cd_ema.value,
                            snap.support, snap.bullish))

    return results


@pytest.fixture(scope="module")
def futures_results(excel_rows):
    return _collect_futures(excel_rows)


def test_bullish_convergence_last_20_rows(futures_results):
    """Bullish: verify _search_bullish(W, ema5_pre, ema20_pre, cd_pre) ≈ 0."""
    from app.engine.futures import _search_bullish
    sample = futures_results[-20:]
    assert len(sample) == 20

    failures = []
    for r, ema5_pre, ema20_pre, cd_pre, support, bullish in sample:
        err = _search_bullish(bullish, ema5_pre, ema20_pre, cd_pre)
        if abs(err) > CONVERGENCE_TOLERANCE:
            failures.append(
                f"row {r['row']}: |Support(Day+2)-W|={abs(err):.6f} > {CONVERGENCE_TOLERANCE} "
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
    for r, ema5_pre, ema20_pre, cd_pre, support, bullish in sample:
        cd_curr = _CD_DECAY * (_ce2(ema5_pre, ema20_pre) - cd_pre) + cd_pre
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
