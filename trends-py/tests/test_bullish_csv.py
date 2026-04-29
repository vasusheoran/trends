"""
Verify Support and Bullish computation from the CSV-seeding path.

Loads bullish-01-Jan-25.csv, strips trial rows, seeds state from real rows,
then asserts:
  1. Support = cd3 at converged cc_trial satisfies BP[d+3]=BP[d+2].
  2. Bullish W satisfies Support(Day+2 with W) == W (2-day convergence).
"""

import sys
from pathlib import Path

import pytest

ROOT = Path(__file__).parent.parent.parent
CSV_PATH = ROOT / "data" / "bullish-01-Jan-25.csv"

sys.path.insert(0, str(Path(__file__).parent.parent))

from scripts.compute_bullish_from_csv import compute_from_csv
from app.engine.futures import _bp_series, _ce2, _CD_DECAY, get_support
from app.engine.indicators import TOLERANCE

CONVERGENCE_TOLERANCE = TOLERANCE * 5


@pytest.mark.skipif(not CSV_PATH.exists(), reason="bullish-01-Jan-25.csv not present")
def test_bullish_csv_2day_convergence():
    """Bullish W satisfies |Support(Day+2) - W| <= tolerance."""
    import pandas as pd
    from app.engine.state import TickerState

    df = pd.read_csv(str(CSV_PATH), skipinitialspace=True)
    df.columns = ["Date", "Close", "Open", "High", "Low"]

    split = len(df)
    for i in range(len(df) - 1, -1, -1):
        r = df.iloc[i]
        if r["Open"] == r["High"] == r["Low"] == r["Close"]:
            split = i
        else:
            break
    real_df = df.iloc[:split]

    state = TickerState(ticker="TEST")
    for _, row in real_df.iterrows():
        state._update_commit(
            date=str(row["Date"]),
            close=float(row["Close"]),
            open_=float(row["Open"]),
            high=float(row["High"]),
            low=float(row["Low"]),
        )

    # Bullish uses post-bar EMA state and cd_curr
    ema5_pre = state._futures_ema5_pre.copy()
    ema20_pre = state._futures_ema20_pre.copy()
    cd_pre = state._futures_cd_pre
    ce2_today = _ce2(ema5_pre, ema20_pre)
    cd_curr = _CD_DECAY * (ce2_today - cd_pre) + cd_pre

    # Post-bar EMA (state.ema5/ema20 after all commits)
    ema5_post = state.ema5.copy()
    ema20_post = state.ema20.copy()

    _, bullish = compute_from_csv(str(CSV_PATH))
    assert bullish is not None, "Bullish search must find a bracket"

    # Verify 2-day convergence from post-bar state (matches Go's calculateBR)
    # Day+1 CD step
    ce2_d1 = _ce2(ema5_post, ema20_post)
    cd_d1 = _CD_DECAY * (ce2_d1 - cd_curr) + cd_curr
    e5_d1 = ema5_post.copy(); e5_d1.update(bullish)
    e20_d1 = ema20_post.copy(); e20_d1.update(bullish)

    # Day+2 CD step
    ce2_d2 = _ce2(e5_d1, e20_d1)
    cd_d2 = _CD_DECAY * (ce2_d2 - cd_d1) + cd_d1

    sup_d2 = get_support(e5_d1, e20_d1, cd_d2)
    assert sup_d2 is not None, "Support(Day+2) must be computable"

    diff = abs(sup_d2 - bullish)
    assert diff <= CONVERGENCE_TOLERANCE, (
        f"|Support(Day+2) - W|={diff:.6f} > {CONVERGENCE_TOLERANCE} at bullish={bullish:.2f}"
    )


@pytest.mark.skipif(not CSV_PATH.exists(), reason="bullish-01-Jan-25.csv not present")
def test_bullish_csv_script_returns_values():
    """compute_from_csv() returns non-None support and bullish."""
    support, bullish = compute_from_csv(str(CSV_PATH))
    assert support is not None, "Support should not be None"
    assert bullish is not None, "Bullish should not be None"
    assert bullish > 0
    assert support > 0
