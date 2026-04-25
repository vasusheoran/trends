"""
Compute Support (CC) and Bullish (BR) for the current trading day from a CSV file.

The CSV must have columns: Date, Close, Open, High, Low
The last 1-2 rows should be "trial" rows where Open == High == Low == Close
(placeholder values the user iterates manually). These are stripped and
EMA + CD state is seeded from all preceding real rows.

Algorithm (verified against Go source master:services/ticker/cards/):
  CE2     = (49*EMA5_pre - 19*EMA20_pre) / 30  [closed-form 2-bar fixed point]
  CD      = EMA-5 (decay 2/6) of daily CE2 values; update with today's CE2
  Bullish = brentq: W=[BR, CE2, CE2], find BR where BP[d+3]=BP[d+2]
  Support = cd3 from brentq: W=[cd3, trial, trial], find trial where BP[d+3]=BP[d+2]
            where cd3 = 2/6*(trial - CD) + CD

Usage:
    cd trends-py
    uv run python scripts/compute_bullish_from_csv.py ../data/bullish-01-Jan-25.csv
"""

import sys
from pathlib import Path

import pandas as pd

sys.path.insert(0, str(Path(__file__).parent.parent))

from app.engine.futures import compute_futures, _ce2
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
    for _, row in real_df.iterrows():
        state._update_commit(
            date=str(row["Date"]),
            close=float(row["Close"]),
            open_=float(row["Open"]),
            high=float(row["High"]),
            low=float(row["Low"]),
        )

    last_real = real_df.iloc[-1]

    # ema5/ema20 after all real rows = pre-bar EMA for the current (trial) day
    ema5_pre = state.ema5.copy()
    ema20_pre = state.ema20.copy()

    # CE2 for current day, then update CD with it
    ce2 = _ce2(ema5_pre, ema20_pre)
    cd_ema = state.cd_ema.copy()
    cd_ema.update(ce2)

    if not cd_ema.seeded:
        print("ERROR: Not enough history to seed CD EMA (need 5 CE2 values = 105 real rows).")
        sys.exit(1)

    cd2 = cd_ema.value

    support, bullish = compute_futures(
        ema5_pre=ema5_pre,
        ema20_pre=ema20_pre,
        cd2=cd2,
    )

    print(f"CSV      : {csv_path}")
    print(f"Seeded   : {len(real_df)} rows  (up to {last_real['Date']})")
    print(f"Day      : {current_date}   CE2: {ce2:.2f}   CD: {cd2:.2f}")
    print()
    print(f"Support  : {support:.2f}" if support is not None else "Support  : (no bracket)")
    print(f"Bullish  : {bullish:.2f}" if bullish is not None else "Bullish  : (no bracket)")

    return support, bullish


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <path/to/bullish-DD-Mon-YY.csv>")
        sys.exit(1)
    compute_from_csv(sys.argv[1])
