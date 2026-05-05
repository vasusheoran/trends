"""
End-to-end live-mode test — verifies the unified update path.

Two concerns:
  1. Value correctness: for each of the last 50 Excel rows, the snapshot produced by
     state.update() must match the Excel ground-truth values within tolerance.
     (Same assertions as test_seeding.py, but exercising the code path where
     _pre_bar_state has already been set by previous bars.)

  2. Idempotency: sending the same date three times with different close prices must
     produce identical support / bullish / hold across all three ticks
     (futures are pinned to the day's opening EMA state, not the intraday close),
     while indicators (ema5, ema20) DO change with each close.

Tolerances mirror test_seeding.py.
"""

import pytest
from app.engine.state import TickerState

# EMA and RSI tolerances are 1.0 because the updated Excel stores these values as whole
# numbers (rounded), so any computed value within 1 of the stored integer is correct.
EMA_TOL     = 1.0
RSI_TOL     = 1.0
HL_TOL      = 0.001
AVG_TOL     = 1.0
FUTURES_TOL = 15.0
LIVE_COUNT  = 50


# ---------------------------------------------------------------------------
# Fixture: seed all-but-last-50 rows, then process the last 50 normally
# ---------------------------------------------------------------------------

@pytest.fixture(scope="module")
def live_results(excel_rows):
    """
    Returns list of (row_dict, snapshot) for the last LIVE_COUNT Excel rows.
    The state has already processed all prior rows before each of these ticks,
    so each row's computed values must match Excel ground truth.
    """
    state = TickerState(ticker="LIVE_TEST")
    seed_rows = excel_rows[:-LIVE_COUNT]
    live_rows = excel_rows[-LIVE_COUNT:]

    for r in seed_rows:
        state.update(r["date"], r["close"], r["open"], r["high"], r["low"])

    results = []
    for r in live_rows:
        snap = state.update(r["date"], r["close"], r["open"], r["high"], r["low"])
        results.append((r, snap))
    return results


# ---------------------------------------------------------------------------
# Value-match tests
# ---------------------------------------------------------------------------

def test_live_ema5_matches_excel(live_results):
    failures = []
    for r, snap in live_results:
        if r["ema5"] is None or snap.ema5 is None:
            continue
        diff = abs(snap.ema5 - r["ema5"])
        if diff > EMA_TOL:
            failures.append(
                f"row {r['row']} ({r['date']}): ema5={snap.ema5:.4f}  excel={r['ema5']:.4f}  diff={diff:.4f}"
            )
    assert not failures, f"{len(failures)} EMA5 mismatches:\n" + "\n".join(failures)


def test_live_ema20_matches_excel(live_results):
    failures = []
    for r, snap in live_results:
        if r["ema20"] is None or snap.ema20 is None:
            continue
        diff = abs(snap.ema20 - r["ema20"])
        if diff > EMA_TOL:
            failures.append(
                f"row {r['row']} ({r['date']}): ema20={snap.ema20:.4f}  excel={r['ema20']:.4f}  diff={diff:.4f}"
            )
    assert not failures, f"{len(failures)} EMA20 mismatches:\n" + "\n".join(failures)


def test_live_hl_matches_excel(live_results):
    failures = []
    for r, snap in live_results:
        if r["hl"] is None or snap.hl is None:
            continue
        diff = abs(snap.hl - r["hl"])
        if diff > HL_TOL:
            failures.append(
                f"row {r['row']} ({r['date']}): hl={snap.hl:.4f}  excel={r['hl']:.4f}  diff={diff:.4f}"
            )
    assert not failures, f"{len(failures)} H/L mismatches:\n" + "\n".join(failures)


def test_live_avg_matches_excel(live_results):
    failures = []
    for r, snap in live_results:
        if r["avg"] is None or snap.avg is None:
            continue
        diff = abs(snap.avg - r["avg"])
        if diff > AVG_TOL:
            failures.append(
                f"row {r['row']} ({r['date']}): avg={snap.avg:.4f}  excel={r['avg']:.4f}  diff={diff:.4f}"
            )
    assert not failures, f"{len(failures)} AVG mismatches:\n" + "\n".join(failures)


def test_live_rsi_matches_excel(live_results):
    failures = []
    for r, snap in live_results:
        if r["rsi"] is None or snap.rsi is None:
            continue
        diff = abs(snap.rsi - r["rsi"])
        if diff > RSI_TOL:
            failures.append(
                f"row {r['row']} ({r['date']}): rsi={snap.rsi:.4f}  excel={r['rsi']:.4f}  diff={diff:.4f}"
            )
    assert not failures, f"{len(failures)} RSI mismatches:\n" + "\n".join(failures)


def test_live_support_matches_excel(live_results):
    with_support = [(r, snap) for r, snap in live_results
                    if r["support"] is not None and snap.support is not None]
    if not with_support:
        pytest.skip("No rows with both Excel and computed support in live sample")
    failures = []
    for r, snap in with_support:
        diff = abs(snap.support - r["support"])
        if diff > FUTURES_TOL:
            failures.append(
                f"row {r['row']} ({r['date']}): support={snap.support:.2f}  excel={r['support']:.2f}  diff={diff:.2f}"
            )
    assert not failures, f"{len(failures)} Support mismatches:\n" + "\n".join(failures)


def test_live_bullish_sanity(live_results):
    with_bullish = [(r, snap) for r, snap in live_results if snap.bullish is not None]
    if not with_bullish:
        pytest.skip("No rows with computed bullish in live sample")
    for r, snap in with_bullish:
        assert snap.bullish > 0, f"row {r['row']}: bullish must be positive, got {snap.bullish}"
        assert snap.bullish < 200000, f"row {r['row']}: bullish suspiciously large: {snap.bullish}"


# ---------------------------------------------------------------------------
# Idempotency test — same date, different closes, futures stay constant
# ---------------------------------------------------------------------------

def test_same_date_futures_stable(excel_rows):
    """
    Sending three different close prices on the same date must:
    - Keep support / bullish / hold identical across all three ticks
      (futures are pinned to the day's opening EMA state, not the intraday close)
    - Produce different ema5 / ema20 on each tick
      (indicators reflect the latest close)
    """
    state = TickerState("IDEM_TEST")
    for r in excel_rows[:-1]:
        state.update(r["date"], r["close"], r["open"], r["high"], r["low"])

    last = excel_rows[-1]
    date, open_, high, low = last["date"], last["open"], last["high"], last["low"]

    s1 = state.update(date, last["close"],      open_, high, low)
    s2 = state.update(date, last["close"] + 50, open_, high, low)
    s3 = state.update(date, last["close"] - 50, open_, high, low)

    # Futures stable across all same-day ticks
    assert s1.support == s2.support == s3.support, \
        f"Support changed across same-day ticks: {s1.support}, {s2.support}, {s3.support}"
    assert s1.bullish == s2.bullish == s3.bullish, \
        f"Bullish changed across same-day ticks: {s1.bullish}, {s2.bullish}, {s3.bullish}"
    assert s1.hold == s2.hold == s3.hold, \
        f"Hold changed across same-day ticks: {s1.hold}, {s2.hold}, {s3.hold}"

    # Indicators DO change with the close price
    assert s1.ema5 != s2.ema5, "ema5 should differ when close changes"
    assert s2.ema5 != s3.ema5, "ema5 should differ when close changes"
