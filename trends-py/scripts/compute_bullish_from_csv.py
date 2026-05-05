"""
Compute Support (CC) and Bullish for the current trading day from a CSV file.

The CSV must have columns: Date, Close, Open, High, Low.
Trial rows (Open == High == Low == Close) at the tail are stripped; EMA + CD
state is seeded from all preceding real rows.

Algorithm (matches Go source master:services/ticker/cards/):
  CE2     = (49*EMA5 - 19*EMA20) / 30  [closed-form 2-bar fixed point]
  CD      = EMA-5 (decay 2/6) of daily CE2 values
  Support = cc_trial where BP[d+3]==BP[d+2] with W=[cd3, trial, trial];
            Support = cd3 = (2/6)*(trial - CD_curr) + CD_curr
  Bullish = W where Support(Day+2 with W) == W  [2-day convergence]

Usage:
    cd trends-py
    uv run python scripts/compute_bullish_from_csv.py ../data/bullish-01-Jan-25.csv
"""

import sys
from pathlib import Path

import pandas as pd

sys.path.insert(0, str(Path(__file__).parent.parent))

from app.engine.futures import compute_futures, _ce2, get_support, _CD_DECAY
from app.engine.state import TickerState


def _is_trial_row(row: pd.Series) -> bool:
    return row["Open"] == row["High"] == row["Low"] == row["Close"]


def compute_from_csv(csv_path: str):
    df = pd.read_csv(csv_path, skipinitialspace=True)
    df.columns = ["Date", "Close", "Open", "High", "Low"]

    # Split real vs trial rows from the tail
    split = len(df)
    for i in range(len(df) - 1, -1, -1):
        if _is_trial_row(df.iloc[i]):
            split = i
        else:
            break

    real_df = df.iloc[:split]
    trial_df = df.iloc[split:]

    if real_df.empty:
        print("ERROR: No real (non-trial) rows found in CSV.")
        sys.exit(1)

    current_date = trial_df.iloc[0]["Date"] if not trial_df.empty else "(unknown)"

    # Seed EMA + CD state from all real rows
    state = TickerState(ticker="CSV")
    last_snap = None
    for _, row in real_df.iterrows():
        last_snap = state.update(
            date=str(row["Date"]),
            close=float(row["Close"]),
            open_=float(row["Open"]),
            high=float(row["High"]),
            low=float(row["Low"]),
        )

    last_real = real_df.iloc[-1]
    # _pre_bar_state holds the indicator state captured before the last bar was applied
    ps = state._pre_bar_state
    ema5_pre = ps.ema5
    ema20_pre = ps.ema20
    cd_pre = ps.cd_ema.value

    print(f"CSV      : {csv_path}")
    print(f"Seeded   : {len(real_df)} rows  (up to {last_real['Date']})")
    print(f"Day 1 (Today): {current_date}")
    print(f"  Pre-state: EMA5={ema5_pre.value:.2f}, EMA20={ema20_pre.value:.2f}, CD={cd_pre:.2f}")

    ce2 = _ce2(ema5_pre, ema20_pre)
    cd_curr = _CD_DECAY * (ce2 - cd_pre) + cd_pre
    print(f"  CE2={ce2:.2f}, CD_curr={cd_curr:.2f}")

    support = last_snap.support if last_snap else None
    bullish = last_snap.bullish if last_snap else None

    print(f"  Support  = {support:.2f}" if support else "  Support = (no bracket)")
    print(f"  Bullish  = {bullish:.2f}" if bullish is not None else "  Bullish  = (no bracket)")

    # Verification: apply Bullish W for 2 days from post-bar EMA state and check Support(Day+2)
    if bullish is not None:
        ce2_d1 = _ce2(state.ema5, state.ema20)
        cd_d1 = _CD_DECAY * (ce2_d1 - cd_curr) + cd_curr
        e5_d1 = state.ema5.copy(); e5_d1.update(bullish)
        e20_d1 = state.ema20.copy(); e20_d1.update(bullish)
        ce2_d2 = _ce2(e5_d1, e20_d1)
        cd_d2 = _CD_DECAY * (ce2_d2 - cd_d1) + cd_d1
        sup_d2 = get_support(e5_d1, e20_d1, cd_d2)
        print()
        print(f"Verify: W={bullish:.2f} for 2 days → Support(Day+2)={sup_d2:.2f}  (should ≈ W)")

    return support, bullish


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <path/to/bullish-DD-Mon-YY.csv>")
        sys.exit(1)
    compute_from_csv(sys.argv[1])
