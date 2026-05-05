"""
Hold validation.

Algorithm: Find trial (D+1 close) where Support(D+3) == Bullish when D+2 applies Bullish.
Equivalently: Bullish computed from D+1 state (after applying trial) == today's Bullish.
_search_hold(trial) == 0  ←→  _search_bullish(bullish, ema_after_trial, cd_after_trial) == 0

CSV reference: data/hold-01-Jan-25.csv
  Seed through 30-Dec-24, send 01-Jan-25 tick (date change triggers 31-Dec computation).
  31-Dec-24 settled state → Bullish ≈ 23585, Hold ≈ 23812.
  _search_hold(hold, ema5_pre, ema20_pre, cd_pre, bullish) ≈ 0 (convergence).
"""

import pytest
from app.engine.futures import compute_futures, _search_hold
from app.engine.indicators import EMAState, TOLERANCE


def _ema5(v: float) -> EMAState:
    return EMAState(period=5, decay=2 / 6, value=v, seeded=True)


def _ema20(v: float) -> EMAState:
    return EMAState(period=20, decay=2 / 21, value=v, seeded=True)


def test_hold_convergence(excel_rows):
    """At the computed Hold price, _search_hold residual must be ≈ 0."""
    sample = [r for r in excel_rows if r["ema5"] and r["ema20"]][-10:]

    failures = []
    for row in sample:
        idx = next(i for i, r in enumerate(excel_rows) if r["row"] == row["row"])
        if idx <= 0:
            continue
        prev = excel_rows[idx - 1]
        if prev["cd"] is None:
            continue

        ema5_pre = _ema5(prev["ema5"])
        ema20_pre = _ema20(prev["ema20"])
        cd_pre = prev["cd"]

        _, bullish, hold = compute_futures(ema5_pre, ema20_pre, cd_pre)
        if hold is None or bullish is None:
            continue

        err = _search_hold(hold, ema5_pre, ema20_pre, cd_pre, bullish)
        if abs(err) >= 0.1:
            failures.append(
                f"row {row['row']}: _search_hold residual={err:.4f} at hold={hold:.2f}, bullish={bullish:.2f}"
            )

    assert not failures, "Hold convergence failures:\n" + "\n".join(failures)


def test_hold_csv_reference():
    """
    Seed from hold-01-Jan-25.csv through 30-Dec-24, send 31-Dec as a live update.

    In the unified path, 31-Dec's futures are computed immediately when the 31-Dec bar
    is processed (using pre-31-Dec EMA = post-30-Dec settled state):
      - 31-Dec PUT: bullish ≈ 23530.90, hold ≈ 23757
    """
    import csv
    from pathlib import Path
    from app.engine.state import TickerState

    csv_path = Path(__file__).parent.parent.parent / "data" / "hold-01-Jan-25.csv"
    if not csv_path.exists():
        pytest.skip("hold-01-Jan-25.csv not found")

    rows = []
    with open(csv_path) as f:
        reader = csv.DictReader(f)
        for row in reader:
            rows.append({
                "date": row["Date "].strip(),
                "close": float(row["Close "]),
                "open": float(row["Open "]),
                "high": float(row["High "]),
                "low": float(row["Low "]),
            })

    state = TickerState(ticker="NIFTY")
    for r in rows:
        if r["date"] == "31-Dec-24":
            break
        state.update(r["date"], r["close"], r["open"], r["high"], r["low"])

    r31 = next(r for r in rows if r["date"] == "31-Dec-24")

    # 31-Dec is a new date → futures computed immediately from pre-31-Dec EMA (= post-30-Dec)
    snap = state.update(r31["date"], r31["close"], r31["open"], r31["high"], r31["low"])

    # pre_bar_state captured at the start of 31-Dec holds the post-30-Dec indicator state
    pre31 = state._pre_bar_state

    assert snap.bullish is not None, "31-Dec bullish should be computed"
    assert abs(snap.bullish - 23585.43) < 1, f"Bullish={snap.bullish:.2f}, expected ≈23585.43"
    assert abs(snap.hold - 23812) < 5, f"Hold={snap.hold:.2f}, expected ≈23812"

    # Convergence: _search_hold(hold, ema5_pre, ema20_pre, cd_pre, bullish) ≈ 0
    from app.engine.futures import _search_hold
    err = _search_hold(snap.hold, pre31.ema5, pre31.ema20, pre31.cd_ema.value, snap.bullish)
    assert abs(err) < 1, f"_search_hold residual={err:.4f} at hold={snap.hold:.2f}"
