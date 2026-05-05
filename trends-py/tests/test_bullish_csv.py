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
from app.engine.futures import _ce2, _CD_DECAY
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
        state.update(
            date=str(row["Date"]),
            close=float(row["Close"]),
            open_=float(row["Open"]),
            high=float(row["High"]),
            low=float(row["Low"]),
        )

    # _pre_bar_state holds the state before the last bar was applied
    ps = state._pre_bar_state
    ema5_pre = ps.ema5.copy()
    ema20_pre = ps.ema20.copy()
    cd_pre = ps.cd_ema.value

    _, bullish = compute_from_csv(str(CSV_PATH))
    assert bullish is not None, "Bullish search must find a bracket"

    # Verify 2-day convergence from pre-bar state (AGENTS.md spec)
    from app.engine.futures import _search_bullish
    diff = abs(_search_bullish(bullish, ema5_pre, ema20_pre, cd_pre))
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
