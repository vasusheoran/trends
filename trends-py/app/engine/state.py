"""
Per-ticker rolling state manager.
Holds the last 101 bars of OHLCV + all EMA/SMA/RSI state.
"""

from collections import deque
from dataclasses import dataclass, field
from typing import Optional

from app.engine.indicators import EMAState, SMAState, RSIState, calc_hl, calc_avg
from app.engine.futures import compute_futures
from app.models import TickerSnapshot

# Rolling window size — matches Go Capacity = 101
_CAPACITY = 101
# Minimum bars before futures are attempted
_FUTURES_MIN = 100


@dataclass
class Bar:
    date: str
    close: float
    open: float
    high: float
    low: float


@dataclass
class TickerState:
    ticker: str

    # Rolling OHLCV history (source of truth)
    bars: deque = field(default_factory=lambda: deque(maxlen=_CAPACITY))

    # EMA state — M(5), O(20)
    # Source: Final-bullish-ce.xlsx — BN uses 2/21 decay (EMA-20), no EMA-50 in sheet.
    ema5: EMAState = field(default_factory=lambda: EMAState(period=5, decay=2/6))
    ema20: EMAState = field(default_factory=lambda: EMAState(period=20, decay=2/21))

    # SMA state for AVG
    sma10: SMAState = field(default_factory=lambda: SMAState(period=10))
    sma50: SMAState = field(default_factory=lambda: SMAState(period=50))

    # RSI state
    rsi: RSIState = field(default_factory=RSIState)

    # Rolling EMA5 of CE values (CD in Go model); seeded to 0 before first CE
    cd: float = 0.0

    # Last computed High (to gate futures recalculation)
    _last_futures_high: float = 0.0

    def update(self, date: str, close: float, open_: float, high: float, low: float) -> TickerSnapshot:
        bar = Bar(date=date, close=close, open=open_, high=high, low=low)
        self.bars.append(bar)

        # Snapshot EMA state BEFORE updating — futures binary search needs pre-bar state
        # (matches Go: cleanUpEMA reverts to Index before recalculating Index+1..+5)
        ema5_pre = self.ema5.copy()
        ema20_pre = self.ema20.copy()

        # Indicators
        m = self.ema5.update(close)
        o = self.ema20.update(close)
        sma10_val = self.sma10.update(close)
        sma50_val = self.sma50.update(close)
        rsi_val = self.rsi.update(close)

        highs = deque(b.high for b in self.bars)
        hl = calc_hl(highs)
        avg = calc_avg(sma10_val, sma50_val)

        # Futures — only recalculate when High is new (or first time)
        support = None
        bullish = None
        if len(self.bars) >= _FUTURES_MIN and (high > self._last_futures_high or self._last_futures_high == 0.0):
            self._last_futures_high = high
            bp = (m - o) if (m is not None and o is not None) else None
            if bp is not None:
                support, bullish, self.cd = compute_futures(
                    ema5=ema5_pre,
                    ema20=ema20_pre,
                    old_cd=self.cd,
                )

        return TickerSnapshot(
            ticker=self.ticker,
            date=date,
            close=close,
            open=open_,
            high=high,
            low=low,
            hl=hl,
            avg=avg,
            ema5=m,
            ema20=o,
            ema50=None,
            rsi=rsi_val,
            support=support,
            bullish=bullish,
        )
