"""
Bullish (2-day convergence) validation.

Algorithm: Find W such that applying W for Day+1 and Day+2 makes Support(Day+2) == W.
"""

import pytest
from app.engine.futures import compute_futures, _search_bullish
from app.engine.indicators import EMAState, TOLERANCE

def _ema5(v: float) -> EMAState:
    return EMAState(period=5, decay=2 / 6, value=v, seeded=True)

def _ema20(v: float) -> EMAState:
    return EMAState(period=20, decay=2 / 21, value=v, seeded=True)

def test_bullish_convergence_internal(excel_rows):
    """At the computed Bullish price W, _search_bullish(W) should be ≈ 0."""
    sample = [r for r in excel_rows if r["ema5"] and r["ema20"]][-10:]

    failures = []
    for row in sample:
        idx = next(i for i, r in enumerate(excel_rows) if r["row"] == row["row"])
        if idx <= 0:
            continue
        prev = excel_rows[idx - 1]
        ema5_pre = _ema5(prev["ema5"])
        ema20_pre = _ema20(prev["ema20"])
        cd_pre = prev["cd"]

        if cd_pre is None:
            continue

        _, bullish = compute_futures(ema5_pre, ema20_pre, cd_pre)
        if bullish is None:
            continue

        err = _search_bullish(bullish, ema5_pre, ema20_pre, cd_pre)
        if abs(err) >= 0.1:
            failures.append(f"row {row['row']}: _search_bullish residual={err:.4f} at W={bullish:.2f}")

    assert not failures, "Bullish convergence failures:\n" + "\n".join(failures)
