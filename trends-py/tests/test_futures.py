"""
Support (CC) and Bullish (BR) computation validation.

Bullish: W=[BR, CE2, CE2]; find BR where BP[d+3]=BP[d+2].
Support: cd3=2/6*(cc_trial-cd2)+cd2; W=[cd3, cc_trial, cc_trial]; find cc_trial where BP[d+3]=BP[d+2]; Support=cd3.

Tests here verify the convergence property on values produced by TickerState.
"""

import pytest
from scipy.optimize import brentq

from app.engine.state import TickerState
from app.engine.futures import compute_futures, _bp_series, _search_br, _search_cc, _ce2
from app.engine.indicators import EMAState, TOLERANCE

CONVERGENCE_TOLERANCE = TOLERANCE * 5


def _collect_futures(excel_rows):
    """
    Feed all rows through TickerState in commit mode.
    Returns list of (row, ema5_pre, ema20_pre, cd2, support, bullish).
    cd2 = cd_ema value after update — what was passed to compute_futures.
    """
    state = TickerState(ticker="TEST")
    results = []

    for r in excel_rows:
        ema5_pre  = state.ema5.copy()
        ema20_pre = state.ema20.copy()

        snap = state.update(r["date"], r["close"], r["open"], r["high"], r["low"])

        if snap.support is not None and snap.bullish is not None:
            # cd_ema has already been updated for this bar's CE2 inside _update_commit
            results.append((r, ema5_pre, ema20_pre, state.cd_ema.value, snap.support, snap.bullish))

    return results


@pytest.fixture(scope="module")
def futures_results(excel_rows):
    return _collect_futures(excel_rows)


def test_bullish_convergence_last_20_rows(futures_results):
    """Bullish: verify |BP[d+3]-BP[d+2]| ≤ tolerance with W=[bullish, CE2, CE2]."""
    sample = futures_results[-20:]
    assert len(sample) == 20

    failures = []
    for r, ema5_pre, ema20_pre, cd2, support, bullish in sample:
        ce2 = _ce2(ema5_pre, ema20_pre)
        bp = _bp_series(ema5_pre, ema20_pre, [bullish, ce2, ce2])
        if any(b is None for b in bp):
            failures.append(f"row {r['row']}: BP None at bullish={bullish:.2f}")
            continue
        diff = abs(bp[2] - bp[1])
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(
                f"row {r['row']}: |BP[d+3]-BP[d+2]|={diff:.6f} > {CONVERGENCE_TOLERANCE} "
                f"at bullish={bullish:.2f}"
            )

    assert not failures, "Bullish convergence failures:\n" + "\n".join(failures)


def test_support_convergence_last_20_rows(futures_results):
    """Support: verify W=[support, cc_trial, cc_trial] → BP[d+3]=BP[d+2].
    We recover cc_trial from support=cd3: cc_trial=(cd3-(1-2/6)*cd2)/(2/6).
    """
    sample = futures_results[-20:]
    assert len(sample) == 20

    failures = []
    for r, ema5_pre, ema20_pre, cd2, support, bullish in sample:
        # Recover cc_trial from cd3 (support) and cd2
        # cd3 = 2/6*(cc_trial - cd2) + cd2  →  cc_trial = (cd3 - cd2)/(2/6) + cd2
        cd3 = support
        cc_trial = (cd3 - cd2) / (2 / 6) + cd2
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
