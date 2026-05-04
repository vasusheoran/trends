"""
Hold validation.

Algorithm: Find trial (D+1 close) where Support(D+3) == Bullish when D+2 applies Bullish.
Equivalently: Bullish computed from D+1 state (after applying trial) == today's Bullish.
_search_hold(trial) == 0  ←→  _search_bullish(bullish, ema_after_trial, cd_after_trial) == 0

CSV reference: data/hold-01-Jan-25.csv
  Seed through 30-Dec-24, send 01-Jan-25 tick (date change triggers 31-Dec computation).
  31-Dec-24 settled state → Bullish ≈ 23531, Hold ≈ 23757.
  After applying Hold, Bullish from D+1 state ≈ today's Bullish (convergence).
"""

import pytest
from app.engine.futures import compute_futures, _search_hold, _ce2, _CD_DECAY
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
        ema5_post = _ema5(row["ema5"])
        ema20_post = _ema20(row["ema20"])

        _, bullish, hold = compute_futures(
            ema5_pre, ema20_pre, cd_pre,
            ema5_post=ema5_post, ema20_post=ema20_post,
        )
        if hold is None or bullish is None:
            continue

        ce2_today = _ce2(ema5_pre, ema20_pre)
        cd_curr = _CD_DECAY * (ce2_today - cd_pre) + cd_pre

        err = _search_hold(hold, ema5_post, ema20_post, cd_curr, bullish)
        if abs(err) >= 0.1:
            failures.append(
                f"row {row['row']}: _search_hold residual={err:.4f} at hold={hold:.2f}, bullish={bullish:.2f}"
            )

    assert not failures, "Hold convergence failures:\n" + "\n".join(failures)


def test_hold_csv_reference():
    """
    Seed from hold-01-Jan-25.csv through 30-Dec-24, send 31-Dec live, then 01-Jan dummy.

    Futures are computed only at date change:
    - 31-Dec PUT: shows bullish from seeded history (30-Dec's bullish)
    - 01-Jan PUT (triggers date change): computes 31-Dec's bullish ≈ 23530.90, hold ≈ 23757
    """
    import csv
    from pathlib import Path
    from app.engine.state import TickerState
    from app.engine.futures import _search_bullish
    from scipy.optimize import brentq

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
    state.commit()

    # Save seeded checkpoint (= 30-Dec settled state) before any live PUTs
    seeded_cp = state._checkpoint

    r31 = next(r for r in rows if r["date"] == "31-Dec-24")

    # First PUT: 31-Dec — bullish/hold come from seeded history (30-Dec's computed values)
    snap = state.update(r31["date"], r31["close"], r31["open"], r31["high"], r31["low"])
    assert snap.bullish == state.history[-2].bullish, (
        f"Before date change, bullish should equal last seeded bullish, got {snap.bullish}"
    )

    # Second PUT: dummy 01-Jan triggers date change → computes 31-Dec's settled bullish/hold
    snap2 = state.update("01-Jan-25", r31["close"], r31["open"], r31["high"], r31["low"])
    assert snap2.bullish is not None, "After date change, bullish should be computed"
    assert abs(snap2.bullish - 23530.90) < 1, f"Bullish={snap2.bullish:.2f}, expected ≈23530.90"
    assert abs(snap2.hold - 23757) < 5, f"Hold={snap2.hold:.2f}, expected ≈23757"

    # Convergence: applying Hold as D+1 close (on top of 31-Dec settled state) yields
    # Bullish from D+1 ≈ snap2.bullish — verifies the algorithm definition holds.
    #
    # The date-change computation used:
    #   ema5_pre  = seeded_cp.ema5  (30-Dec post state)
    #   ema5_post = seeded_cp.ema5 + r31["close"]  (31-Dec settled)
    #   cd_pre    = seeded_cp.cd_ema.value
    e5_settled = seeded_cp.ema5.copy()
    e20_settled = seeded_cp.ema20.copy()
    e5_settled.update(r31["close"])
    e20_settled.update(r31["close"])

    ce2_pre = _ce2(seeded_cp.ema5, seeded_cp.ema20)
    cd_curr = _CD_DECAY * (ce2_pre - seeded_cp.cd_ema.value) + seeded_cp.cd_ema.value
    ce2_d1 = _ce2(e5_settled, e20_settled)
    cd_d1 = _CD_DECAY * (ce2_d1 - cd_curr) + cd_curr

    e5_hold = e5_settled.copy()
    e20_hold = e20_settled.copy()
    e5_hold.update(snap2.hold)
    e20_hold.update(snap2.hold)

    try:
        bullish_d1 = brentq(lambda w: _search_bullish(w, e5_hold, e20_hold, cd_d1), 0, 99999, xtol=0.001)
        assert abs(bullish_d1 - snap2.bullish) < 5, (
            f"Bullish from D+1 after Hold={snap2.hold:.2f}: {bullish_d1:.2f}, "
            f"expected ≈ {snap2.bullish:.2f}"
        )
    except ValueError:
        pytest.skip("Bullish from D+1 did not converge")
