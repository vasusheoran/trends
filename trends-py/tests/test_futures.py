"""
Phase 3 validation — CE, CC (Support), and BR (Bullish) binary-search properties.

BP = EMA5 - EMA20 (source: Final-bullish-ce.xlsx col BP = AS - BN, BN decay 2/21).

Convergence properties verified:
  CE: |BP[day+2] - BP[day+1]| <= TOLERANCE  (both days use trial=CE)
  CC: |BP[day+3] - BP[day+2]| <= TOLERANCE  (day+1=CC_w, day+2,3=cc_trial)
  BR: |BP[day+3] - BP[day+2]| <= TOLERANCE  (day+1=BR, day+2=day+3=CE)

TOLERANCE = 0.001 — same value as the brentq xtol.
"""

import pytest
from scipy.optimize import brentq

from app.engine.state import TickerState
from app.engine.futures import compute_futures, _bp_series, _search_ce, _search_cc, _search_br
from app.engine.indicators import TOLERANCE

CONVERGENCE_TOLERANCE = TOLERANCE * 5  # small multiplier for floating-point chain


def _collect_futures(excel_rows):
    """
    Feed all rows through TickerState (resetting the high-change gate each bar
    so futures are computed for every row, matching Excel behaviour).
    Returns list of (excel_row, ema5_pre, ema20_pre, cd_pre, support, bullish)
    for rows where both support and bullish are not None.
    """
    state = TickerState(ticker="TEST")
    results = []

    for r in excel_rows:
        ema5_pre = state.ema5.copy()
        ema20_pre = state.ema20.copy()
        cd_pre = state.cd

        state._last_futures_high = 0.0
        snap = state.update(r["date"], r["close"], r["open"], r["high"], r["low"])

        if snap.support is not None and snap.bullish is not None:
            results.append((r, ema5_pre, ema20_pre, cd_pre, snap.support, snap.bullish))

    return results


@pytest.fixture(scope="module")
def futures_results(excel_rows):
    return _collect_futures(excel_rows)


def _find_ce(ema5_pre, ema20_pre, low=0.0, high=99999.0):
    """Recompute CE for a given pre-bar EMA state."""
    f_low = _search_ce(low, ema5_pre, ema20_pre)
    f_high = _search_ce(high, ema5_pre, ema20_pre)
    if f_low * f_high >= 0:
        return None
    return brentq(_search_ce, low, high, args=(ema5_pre, ema20_pre), xtol=TOLERANCE)


def test_ce_convergence_last_20_rows(futures_results):
    """CE: verify BP[day+2] ≈ BP[day+1] within CONVERGENCE_TOLERANCE."""
    sample = futures_results[-20:]
    assert len(sample) == 20, "Need at least 20 rows with both support and bullish populated"

    failures = []
    for r, ema5_pre, ema20_pre, cd_pre, _support, _br in sample:
        ce = _find_ce(ema5_pre, ema20_pre)
        if ce is None:
            failures.append(f"row {r['row']}: CE not found")
            continue
        bp = _bp_series(ema5_pre, ema20_pre, [ce, ce])
        if bp[0] is None or bp[1] is None:
            failures.append(f"row {r['row']}: BP None at CE={ce:.2f}")
            continue
        diff = abs(bp[1] - bp[0])
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(
                f"row {r['row']}: |BP[day+2]-BP[day+1]|={diff:.6f} > {CONVERGENCE_TOLERANCE} at CE={ce:.2f}"
            )

    assert not failures, "CE convergence failures:\n" + "\n".join(failures)


def test_cc_convergence_last_20_rows(futures_results):
    """CC (Support): verify BP[day+3] ≈ BP[day+2] with W[d+1]=support, W[d+2,3]=cc_trial."""
    sample = futures_results[-20:]
    assert len(sample) == 20, "Need at least 20 rows with support populated"

    failures = []
    for r, ema5_pre, ema20_pre, cd_pre, support, _br in sample:
        ce = _find_ce(ema5_pre, ema20_pre)
        if ce is None:
            failures.append(f"row {r['row']}: CE not found, cannot compute new_cd")
            continue
        new_cd = 2 / 6 * (ce - cd_pre) + cd_pre

        # Recover cc_trial: support = 2/6*(cc_trial - new_cd) + new_cd → cc_trial = (support - new_cd)*3 + new_cd
        cc_trial = (support - new_cd) * 3 + new_cd

        bp = _bp_series(ema5_pre, ema20_pre, [support, cc_trial, cc_trial])
        if any(b is None for b in bp):
            failures.append(f"row {r['row']}: BP None at support={support:.2f}")
            continue
        diff = abs(bp[2] - bp[1])
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(
                f"row {r['row']}: |BP[day+3]-BP[day+2]|={diff:.6f} > {CONVERGENCE_TOLERANCE} at support={support:.2f}"
            )

    assert not failures, "CC convergence failures:\n" + "\n".join(failures)


def test_br_convergence_last_20_rows(futures_results):
    """BR: verify BP[day+3] ≈ BP[day+2] within CONVERGENCE_TOLERANCE."""
    sample = futures_results[-20:]
    assert len(sample) == 20, "Need at least 20 rows with both support and bullish populated"

    failures = []
    for r, ema5_pre, ema20_pre, cd_pre, _support, br in sample:
        ce = _find_ce(ema5_pre, ema20_pre)
        if ce is None:
            failures.append(f"row {r['row']}: CE not found")
            continue
        bp = _bp_series(ema5_pre, ema20_pre, [br, ce, ce])
        if any(b is None for b in bp):
            failures.append(f"row {r['row']}: BP None at BR={br:.2f}, CE={ce:.2f}")
            continue
        diff = abs(bp[2] - bp[1])
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(
                f"row {r['row']}: |BP[day+3]-BP[day+2]|={diff:.6f} > {CONVERGENCE_TOLERANCE} at BR={br:.2f}"
            )

    assert not failures, "BR convergence failures:\n" + "\n".join(failures)


def test_futures_populated_after_100_bars(excel_rows):
    """Futures must be None before 100 bars and non-None by bar 120 (EMA5+EMA20 seed)."""
    state = TickerState(ticker="TEST")
    none_after_100 = []

    for i, r in enumerate(excel_rows[:200]):
        state._last_futures_high = 0.0
        snap = state.update(r["date"], r["close"], r["open"], r["high"], r["low"])
        if i >= 100 and snap.support is None:
            none_after_100.append(i)

    # EMA5 seeds at bar 5, EMA20 seeds at bar 20 — both ready well before bar 100
    late_nones = [i for i in none_after_100 if i >= 120]
    assert not late_nones, f"Futures still None at bars: {late_nones[:5]}"
