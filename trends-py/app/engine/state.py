"""
Per-ticker rolling state manager — single unified update path.

Every update() call:
  - Captures pre-bar indicator state on the first tick of each new date.
  - For same-date repeated ticks: restores to that captured state before applying,
    making back-to-back PUTs with different intraday closes fully idempotent for
    futures (support / bullish / hold stay pinned to the day's opening EMA state).
  - Futures are computed once per new date from the correct pre-bar EMA state.

commit() is retained as a no-op for caller compatibility but does nothing.
"""

from collections import deque
from dataclasses import dataclass, field
from typing import Optional

from app.engine.indicators import EMAState, SMAState, RSIState, calc_hl, calc_avg
from app.engine.futures import compute_futures, _ce2
from app.models import TickerSnapshot

_CAPACITY = 101
_FUTURES_MIN = 100
_HISTORY_MAX = 5500


@dataclass
class Bar:
    date: str
    close: float
    open: float
    high: float
    low: float
    timestamp: Optional[int] = None
    hl: Optional[float] = None
    avg: Optional[float] = None
    ema5: Optional[float] = None
    ema20: Optional[float] = None
    ema50: Optional[float] = None
    rsi: Optional[float] = None
    support: Optional[float] = None
    bullish: Optional[float] = None
    hold: Optional[float] = None


@dataclass
class _PreBarState:
    """Indicator state captured before the current date's first bar was applied."""
    ema5: EMAState
    ema20: EMAState
    ema50: EMAState
    sma10: SMAState
    sma50: SMAState
    rsi: RSIState
    cd_ema: EMAState
    bars: deque


@dataclass
class TickerState:
    ticker: str

    bars: deque = field(default_factory=lambda: deque(maxlen=_CAPACITY))
    ema5: EMAState = field(default_factory=lambda: EMAState(period=5, decay=2/6))
    ema20: EMAState = field(default_factory=lambda: EMAState(period=20, decay=2/21))
    ema50: EMAState = field(default_factory=lambda: EMAState(period=50, decay=2/51))
    sma10: SMAState = field(default_factory=lambda: SMAState(period=10))
    sma50: SMAState = field(default_factory=lambda: SMAState(period=50))
    rsi: RSIState = field(default_factory=RSIState)
    cd_ema: EMAState = field(default_factory=lambda: EMAState(period=5, decay=2/6))

    history: deque = field(default_factory=lambda: deque(maxlen=_HISTORY_MAX))
    bars_1m: deque = field(default_factory=lambda: deque(maxlen=1000))
    bars_5m: deque = field(default_factory=lambda: deque(maxlen=1000))

    # State captured before the current date's first bar (used for idempotency and futures)
    _last_date: Optional[str] = field(default=None, repr=False)
    _pre_bar_state: Optional[_PreBarState] = field(default=None, repr=False)

    # Futures pinned for the current day (computed from pre-bar EMA on the first tick)
    _day_support: Optional[float] = field(default=None, repr=False)
    _day_bullish: Optional[float] = field(default=None, repr=False)
    _day_hold: Optional[float] = field(default=None, repr=False)

    def commit(self):
        """No-op — retained for caller compatibility."""

    def update(self, date: str, close: float, open_: float, high: float, low: float,
               force: bool = False, timestamp: Optional[int] = None) -> TickerSnapshot:
        is_new_date = (date != self._last_date)

        if is_new_date:
            # Freeze state before this date's bar is applied
            self._pre_bar_state = _PreBarState(
                ema5=self.ema5.copy(),
                ema20=self.ema20.copy(),
                ema50=self.ema50.copy(),
                sma10=self.sma10.copy(),
                sma50=self.sma50.copy(),
                rsi=self.rsi.copy(),
                cd_ema=self.cd_ema.copy(),
                bars=deque(self.bars, maxlen=_CAPACITY),
            )
            self.bars.append(Bar(date=date, close=close, open=open_, high=high, low=low, timestamp=timestamp))
            self._last_date = date
            self._day_support = None
            self._day_bullish = None
            self._day_hold = None
            self.bars_1m.clear()
            self.bars_5m.clear()
        else:
            # Same date: restore to pre-bar state so indicators are computed from the same base
            ps = self._pre_bar_state
            self.ema5 = ps.ema5.copy()
            self.ema20 = ps.ema20.copy()
            self.ema50 = ps.ema50.copy()
            self.sma10 = ps.sma10.copy()
            self.sma50 = ps.sma50.copy()
            self.rsi = ps.rsi.copy()
            self.cd_ema = ps.cd_ema.copy()
            self.bars = deque(ps.bars, maxlen=_CAPACITY)
            self.bars.append(Bar(date=date, close=close, open=open_, high=high, low=low, timestamp=timestamp))

        # Pre-bar EMA/CD (constant for the day — idempotent across same-date ticks)
        pre_ema5 = self.ema5.copy()
        pre_ema20 = self.ema20.copy()
        pre_cd = self.cd_ema.value

        # Apply indicators
        m = self.ema5.update(close)
        o = self.ema20.update(close)
        e50 = self.ema50.update(close)
        sma10_val = self.sma10.update(close)
        sma50_val = self.sma50.update(close)
        rsi_val = self.rsi.update(close)

        highs = deque(b.high for b in self.bars)
        hl = calc_hl(highs)
        avg = calc_avg(sma10_val, sma50_val)

        # Futures: compute once per new date using pre-bar EMA (constant for day)
        if len(self.bars) >= _FUTURES_MIN and m is not None and o is not None:
            ce2 = _ce2(pre_ema5, pre_ema20)
            self.cd_ema.update(ce2)
            if self.cd_ema.seeded and pre_cd is not None and is_new_date:
                self._day_support, self._day_bullish, self._day_hold = compute_futures(
                    ema5_pre=pre_ema5, ema20_pre=pre_ema20,
                    cd_pre=pre_cd,
                )

        bar = self.bars[-1]
        bar.hl = hl; bar.avg = avg; bar.ema5 = m; bar.ema20 = o; bar.ema50 = e50
        bar.rsi = rsi_val; bar.support = self._day_support
        bar.bullish = self._day_bullish; bar.hold = self._day_hold

        if self.history and self.history[-1].date == date:
            self.history[-1] = bar
        else:
            self.history.append(bar)

        if timestamp:
            self._aggregate_intraday(self.bars_1m, 60, timestamp, close, open_, high, low)
            self._aggregate_intraday(self.bars_5m, 300, timestamp, close, open_, high, low)

        return TickerSnapshot(
            ticker=self.ticker, date=date, close=close, open=open_, high=high, low=low,
            hl=hl, avg=avg, ema5=m, ema20=o, ema50=e50, rsi=rsi_val,
            support=self._day_support, bullish=self._day_bullish, hold=self._day_hold,
            timestamp=timestamp,
        )

    def _aggregate_intraday(self, target_deque: deque, period_sec: int, ts: int,
                             close: float, open_: float, high: float, low: float):
        period_ts = (ts // period_sec) * period_sec
        if target_deque and target_deque[-1].timestamp == period_ts:
            b = target_deque[-1]
            b.close = close
            b.high = max(b.high, high)
            b.low = min(b.low, low)
        else:
            target_deque.append(Bar(
                date="", close=close, open=open_, high=high, low=low, timestamp=period_ts
            ))
