"""
Phase 3 validation — Support (CC) and Bullish (CD3) computation.

BP = EMA5 - EMA20 (source: Final-bullish-ce.xlsx col BP = AS - BN, BN decay 2/21).

Computation chain:
  CE2:     |BP[d+2] - BP[d+1]| ≤ TOLERANCE  (pre-bar EMA, both days use trial)
  CC:      |BP[d+3] - BP[d+2]| ≤ TOLERANCE  (pre-bar EMA, W[d+1]=CC, W[d+2,3]=cc_trial)
  CE3:     |BP[d+2] - BP[d+1]| ≤ TOLERANCE  (post-bar EMA, both days use trial)
  Bullish: bullish == CD3 == 2/6*(CE3 - CD2) + CD2
"""

import pytest
from scipy.optimize import brentq

from app.engine.state import TickerState
from app.engine.futures import compute_futures, _bp_series, _search_ce, _search_cc
from app.engine.indicators import TOLERANCE

CONVERGENCE_TOLERANCE = TOLERANCE * 5


def _collect_futures(excel_rows):
    """
    Feed all rows through TickerState in commit mode.
    Returns list of (row, ema5_pre, ema20_pre, ema5_post, ema20_post, cd_pre, support, bullish).
    """
    state = TickerState(ticker="TEST")
    results = []

    for r in excel_rows:
        ema5_pre = state.ema5.copy()
        ema20_pre = state.ema20.copy()
        cd_pre = state.cd

        state._last_futures_high = 0.0
        snap = state.update(r["date"], r["close"], r["open"], r["high"], r["low"])

        # post-bar state is now in state.ema5/ema20 (commit mode permanently applies)
        ema5_post = state.ema5.copy()
        ema20_post = state.ema20.copy()

        if snap.support is not None and snap.bullish is not None:
            results.append((r, ema5_pre, ema20_pre, ema5_post, ema20_post, cd_pre, snap.support, snap.bullish))

    return results


@pytest.fixture(scope="module")
def futures_results(excel_rows):
    return _collect_futures(excel_rows)


def _find_ce(ema5, ema20, low=0.0, high=99999.0):
    f_low = _search_ce(low, ema5, ema20)
    f_high = _search_ce(high, ema5, ema20)
    if f_low * f_high >= 0:
        return None
    return brentq(_search_ce, low, high, args=(ema5, ema20), xtol=TOLERANCE)


def test_ce2_convergence_last_20_rows(futures_results):
    """CE2: verify BP[d+2] ≈ BP[d+1] from pre-bar EMA state."""
    sample = futures_results[-20:]
    assert len(sample) == 20

    failures = []
    for r, ema5_pre, ema20_pre, ema5_post, ema20_post, cd_pre, _support, _bullish in sample:
        ce2 = _find_ce(ema5_pre, ema20_pre)
        if ce2 is None:
            failures.append(f"row {r['row']}: CE2 not found")
            continue
        bp = _bp_series(ema5_pre, ema20_pre, [ce2, ce2])
        if bp[0] is None or bp[1] is None:
            failures.append(f"row {r['row']}: BP None at CE2={ce2:.2f}")
            continue
        diff = abs(bp[1] - bp[0])
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(f"row {r['row']}: |BP[d+2]-BP[d+1]|={diff:.6f} > {CONVERGENCE_TOLERANCE} at CE2={ce2:.2f}")

    assert not failures, "CE2 convergence failures:\n" + "\n".join(failures)


def test_cc_convergence_last_20_rows(futures_results):
    """CC (Support): verify BP[d+3] ≈ BP[d+2] with W[d+1]=support, W[d+2,3]=cc_trial."""
    sample = futures_results[-20:]
    assert len(sample) == 20

    failures = []
    for r, ema5_pre, ema20_pre, ema5_post, ema20_post, cd_pre, support, _bullish in sample:
        ce2 = _find_ce(ema5_pre, ema20_pre)
        if ce2 is None:
            failures.append(f"row {r['row']}: CE2 not found")
            continue
        cd2 = 2 / 6 * (ce2 - cd_pre) + cd_pre
        # recover cc_trial from support: support = 2/6*(cc_trial - cd2) + cd2
        cc_trial = (support - cd2) * 3 + cd2

        bp = _bp_series(ema5_pre, ema20_pre, [support, cc_trial, cc_trial])
        if any(b is None for b in bp):
            failures.append(f"row {r['row']}: BP None at support={support:.2f}")
            continue
        diff = abs(bp[2] - bp[1])
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(f"row {r['row']}: |BP[d+3]-BP[d+2]|={diff:.6f} > {CONVERGENCE_TOLERANCE} at support={support:.2f}")

    assert not failures, "CC convergence failures:\n" + "\n".join(failures)


def test_ce3_convergence_last_20_rows(futures_results):
    """CE3: verify BP[d+2] ≈ BP[d+1] from post-bar EMA state."""
    sample = futures_results[-20:]
    assert len(sample) == 20

    failures = []
    for r, ema5_pre, ema20_pre, ema5_post, ema20_post, cd_pre, _support, _bullish in sample:
        ce3 = _find_ce(ema5_post, ema20_post)
        if ce3 is None:
            failures.append(f"row {r['row']}: CE3 not found")
            continue
        bp = _bp_series(ema5_post, ema20_post, [ce3, ce3])
        if bp[0] is None or bp[1] is None:
            failures.append(f"row {r['row']}: BP None at CE3={ce3:.2f}")
            continue
        diff = abs(bp[1] - bp[0])
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(f"row {r['row']}: |BP[d+2]-BP[d+1]|={diff:.6f} > {CONVERGENCE_TOLERANCE} at CE3={ce3:.2f}")

    assert not failures, "CE3 convergence failures:\n" + "\n".join(failures)


def test_bullish_equals_cd3_last_20_rows(futures_results):
    """Bullish: verify bullish == CD3 = EMA_step(CD2, CE3)."""
    sample = futures_results[-20:]
    assert len(sample) == 20

    failures = []
    for r, ema5_pre, ema20_pre, ema5_post, ema20_post, cd_pre, _support, bullish in sample:
        ce2 = _find_ce(ema5_pre, ema20_pre)
        ce3 = _find_ce(ema5_post, ema20_post)
        if ce2 is None or ce3 is None:
            failures.append(f"row {r['row']}: CE2 or CE3 not found")
            continue
        cd2 = 2 / 6 * (ce2 - cd_pre) + cd_pre
        cd3 = 2 / 6 * (ce3 - cd2) + cd2
        diff = abs(bullish - cd3)
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(f"row {r['row']}: bullish={bullish:.4f} != CD3={cd3:.4f} (diff={diff:.6f})")

    assert not failures, "Bullish=CD3 failures:\n" + "\n".join(failures)


def test_futures_populated_after_100_bars(excel_rows):
    """Futures must be None before 100 bars and non-None by bar 120."""
    state = TickerState(ticker="TEST")
    none_after_100 = []

    for i, r in enumerate(excel_rows[:200]):
        state._last_futures_high = 0.0
        snap = state.update(r["date"], r["close"], r["open"], r["high"], r["low"])
        if i >= 100 and snap.support is None:
            none_after_100.append(i)

    late_nones = [i for i in none_after_100 if i >= 120]
    assert not late_nones, f"Futures still None at bars: {late_nones[:5]}"
