
import pandas as pd
import sys
from pathlib import Path
from scipy.optimize import brentq

# Add trends-py to path
sys.path.insert(0, str(Path.cwd() / "trends-py"))

from app.engine.indicators import EMAState, TOLERANCE
from app.engine.futures import _ce2, _search_br, _search_cc, _bp_series
from app.engine.state import TickerState

def get_state_after_31dec():
    df = pd.read_csv("../data/bullish-01-Jan-25.csv", skipinitialspace=True)
    df.columns = ["Date", "Close", "Open", "High", "Low"]
    
    # Real rows are up to 31-Dec-24 (index 5386 if we count from 0, but let's just find it)
    real_df = df[df["Date"] <= "31-Dec-24"]
    
    state = TickerState(ticker="CSV")
    for _, row in real_df.iterrows():
        state._update_commit(
            date=str(row["Date"]),
            close=float(row["Close"]),
            open_=float(row["Open"]),
            high=float(row["High"]),
            low=float(row["Low"]),
        )
    return state

def get_support_for_day(state, price):
    # Apply price to state to get state for "today"
    state_today = TickerState(ticker="Today")
    state_today.ema5 = state.ema5.copy()
    state_today.ema20 = state.ema20.copy()
    state_today.cd_ema = state.cd_ema.copy()
    state_today.bars = state.bars.copy()
    
    state_today._update_commit("01-Jan-25", price, price, price, price)
    
    # Now compute futures for "tomorrow"
    # compute_futures needs ema_pre, ema20_pre, cd2
    # For tomorrow, ema_pre is state_today.ema
    
    ema5_pre = state_today.ema5.copy()
    ema20_pre = state_today.ema20.copy()
    
    # cd2 for tomorrow = update CD with CE2(tomorrow)
    ce2_tomorrow = _ce2(ema5_pre, ema20_pre)
    cd_ema_tomorrow = state_today.cd_ema.copy()
    cd_ema_tomorrow.update(ce2_tomorrow)
    cd2 = cd_ema_tomorrow.value
    
    # Now find support
    cc_trial = brentq(_search_cc, 0, 99999, args=(ema5_pre, ema20_pre, cd2))
    support = (2/6) * (cc_trial - cd2) + cd2
    return support

def main():
    state_after_31dec = get_state_after_31dec()
    
    # Old Bullish result from script was 23687.72
    old_bullish = 23687.72
    print(f"Old Bullish: {old_bullish}")
    
    support_at_old = get_support_for_day(state_after_31dec, old_bullish)
    print(f"Support(tomorrow) if Close={old_bullish}: {support_at_old:.2f}")
    
    # Let's find the price P such that Support(tomorrow) == P
    def objective(p):
        return get_support_for_day(state_after_31dec, p) - p
    
    new_bullish = brentq(objective, 20000, 30000)
    print(f"New Bullish (Support next day == price): {new_bullish:.2f}")
    
    # What if price is 23531?
    support_at_23531 = get_support_for_day(state_after_31dec, 23531)
    print(f"Support(tomorrow) if Close=23531: {support_at_23531:.2f}")

if __name__ == "__main__":
    main()
