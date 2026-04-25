"""
Verify bullish (BR) computation from the CSV-seeding path.

Loads bullish-01-Jan-25.csv, strips trial rows, seeds EMA and CD state from
real rows, then asserts mathematical convergence of the computed BR value.
"""

import sys
from pathlib import Path

import pytest
from scipy.optimize import brentq

ROOT = Path(__file__).parent.parent.parent
CSV_PATH = ROOT / "data" / "bullish-01-Jan-25.csv"

sys.path.insert(0, str(Path(__file__).parent.parent))

from scripts.compute_bullish_from_csv import compute_from_csv
from app.engine.futures import _bp_series, _search_br, _ce2
from app.engine.indicators import TOLERANCE

CONVERGENCE_TOLERANCE = TOLERANCE * 5


@pytest.mark.skipif(not CSV_PATH.exists(), reason="bullish-01-Jan-25.csv not present")
def test_bullish_csv_converges():
    """BR found from CSV seeding satisfies |BP[d+3]-BP[d+2]| <= tolerance."""
    import pandas as pd
    from app.engine.state import TickerState

    df = pd.read_csv(str(CSV_PATH), skipinitialspace=True)
    df.columns = ["Date", "Close", "Open", "High", "Low"]

    def _is_trial(row):
        return row["Open"] == row["High"] == row["Low"] == row["Close"]

    split = len(df)
    for i in range(len(df) - 1, -1, -1):
        if _is_trial(df.iloc[i]):
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

    ema5_pre = state.ema5.copy()
    ema20_pre = state.ema20.copy()
    ce2 = _ce2(ema5_pre, ema20_pre)

    fl = _search_br(0.0, ema5_pre, ema20_pre, ce2)
    fh = _search_br(99999.0, ema5_pre, ema20_pre, ce2)
    assert fl * fh < 0, f"Bullish bracket has no sign change: fl={fl:.4f} fh={fh:.4f}"

    bullish = brentq(_search_br, 0.0, 99999.0, args=(ema5_pre, ema20_pre, ce2), xtol=TOLERANCE)

    bp = _bp_series(ema5_pre, ema20_pre, [bullish, ce2, ce2])
    diff = abs(bp[2] - bp[1])
    assert diff <= CONVERGENCE_TOLERANCE, (
        f"|BP[d+3]-BP[d+2]|={diff:.6f} > {CONVERGENCE_TOLERANCE} at bullish={bullish:.2f}"
    )


@pytest.mark.skipif(not CSV_PATH.exists(), reason="bullish-01-Jan-25.csv not present")
def test_bullish_csv_script_returns_values():
    """compute_from_csv() returns non-None support and bullish."""
    support, bullish = compute_from_csv(str(CSV_PATH))
    assert support is not None, "Support should not be None"
    assert bullish is not None, "Bullish should not be None"
    assert bullish > 0
    assert support > 0
