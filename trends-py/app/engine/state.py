"""
Per-ticker rolling state manager.

Two modes:
  Commit mode  — no checkpoint set; update() permanently appends each bar.
                 Used during seeding. Indicators and futures computed normally.
  Live mode    — checkpoint set via commit(); update() restores from checkpoint
                 before applying the bar, making repeated PUTs idempotent.
                 When the date changes, previous live bar is promoted to a new
                 checkpoint automatically.

Futures (support / bullish) are computed once per day from the settled EMA state
and previous day's high. Pass force=True to recompute on every tick.
"""

from collections import deque
from dataclasses import dataclass, field
from typing import Optional

from app.engine.indicators import EMAState, SMAState, RSIState, calc_hl, calc_avg
from app.engine.futures import compute_futures, _ce2
from app.models import TickerSnapshot

_CAPACITY = 101
_FUTURES_MIN = 100
_HISTORY_MAX = 5500  # retain full seed history for charting


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
class _Checkpoint:
    """Frozen indicator state after all committed (historical) bars."""
    ema5: EMAState
    ema20: EMAState
    ema50: EMAState
    sma10: SMAState
    sma50: SMAState
    rsi: RSIState
    bars: deque
    cd_ema: EMAState
    # EMA state from before the last committed bar — used by /api/state endpoint.
    ema5_pre: Optional[EMAState]
    ema20_pre: Optional[EMAState]
    cd_pre: Optional[float]


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
    # CD: EMA-5 (decay 2/6) of daily CE2 values, accumulated from bar 100 onward.
    cd_ema: EMAState = field(default_factory=lambda: EMAState(period=5, decay=2/6))

    # Pre-last-bar state saved during _update_commit, used for futures at checkpoint.
    _futures_ema5_pre: Optional[EMAState] = field(default=None, repr=False)
    _futures_ema20_pre: Optional[EMAState] = field(default=None, repr=False)
    _futures_cd_pre: Optional[float] = field(default=None, repr=False)

    # Full bar history for charting (not capped to _CAPACITY)
    history: deque = field(default_factory=lambda: deque(maxlen=_HISTORY_MAX))

    # Intraday bars (last 1000 bars for 1m and 5m)
    bars_1m: deque = field(default_factory=lambda: deque(maxlen=1000))
    bars_5m: deque = field(default_factory=lambda: deque(maxlen=1000))

    # Live-mode state (populated after commit())
    _checkpoint: Optional[_Checkpoint] = field(default=None, repr=False)
    _live_date: Optional[str] = field(default=None, repr=False)
    _live_support: Optional[float] = field(default=None, repr=False)
    _live_bullish: Optional[float] = field(default=None, repr=False)
    _live_hold: Optional[float] = field(default=None, repr=False)

    def commit(self):
        """
        Freeze current state as the end-of-history checkpoint and switch to live mode.
        After this, update() restores from this checkpoint before each bar so that
        repeated calls with the same data produce identical results.
        """
        self._checkpoint = _make_checkpoint(self)
        self._live_date = None
        self._live_support = None
        self._live_bullish = None
        self._live_hold = None

    def update(self, date: str, close: float, open_: float, high: float, low: float,
               force: bool = False, timestamp: Optional[int] = None) -> TickerSnapshot:
        if self._checkpoint is not None:
            return self._update_live(date, close, open_, high, low, force, timestamp)
        return self._update_commit(date, close, open_, high, low, timestamp)

    def _update_commit(self, date, close, open_, high, low, timestamp=None) -> TickerSnapshot:
        """Permanently append a bar (commit/seed mode)."""
        bar = Bar(date=date, close=close, open=open_, high=high, low=low, timestamp=timestamp)
        self.bars.append(bar)

        ema5_pre = self.ema5.copy()
        ema20_pre = self.ema20.copy()
        cd_pre = self.cd_ema.value

        # Save for checkpoint futures computation
        self._futures_ema5_pre = ema5_pre
        self._futures_ema20_pre = ema20_pre
        self._futures_cd_pre = cd_pre

        m = self.ema5.update(close)
        o = self.ema20.update(close)
        e50 = self.ema50.update(close)
        sma10_val = self.sma10.update(close)
        sma50_val = self.sma50.update(close)
        rsi_val = self.rsi.update(close)

        highs = deque(b.high for b in self.bars)
        hl = calc_hl(highs)
        avg = calc_avg(sma10_val, sma50_val)

        support = None
        bullish = None
        hold = None
        if len(self.bars) >= _FUTURES_MIN and m is not None and o is not None:
            ce2 = _ce2(ema5_pre, ema20_pre)
            self.cd_ema.update(ce2)
            if self.cd_ema.seeded and cd_pre is not None:
                support, bullish, hold = compute_futures(
                    ema5_pre=ema5_pre, ema20_pre=ema20_pre,
                    cd_pre=cd_pre,
                    ema5_post=self.ema5, ema20_post=self.ema20,
                )

        # Enrich bar with computed indicators, then add to full history
        bar.hl = hl; bar.avg = avg; bar.ema5 = m; bar.ema20 = o; bar.ema50 = e50
        bar.rsi = rsi_val; bar.support = support; bar.bullish = bullish; bar.hold = hold
        self.history.append(bar)

        return TickerSnapshot(
            ticker=self.ticker, date=date, close=close, open=open_, high=high, low=low,
            hl=hl, avg=avg, ema5=m, ema20=o, ema50=e50, rsi=rsi_val,
            support=support, bullish=bullish, hold=hold,
            timestamp=timestamp,
        )

    def _update_live(self, date, close, open_, high, low, force: bool, timestamp: Optional[int]) -> TickerSnapshot:
        """
        Apply bar on top of checkpoint (live mode).
        Idempotent: same inputs always produce the same snapshot.
        When date changes, the previous live bar is promoted to a new checkpoint.
        """
        cp = self._checkpoint
        warning = None

        # On first ever live tick: inherit support/bullish/hold from the applicable seeded bar.
        # Same-day replay → use history[-2] (D-1 values apply during D's session).
        # New day → use history[-1] (last seeded bar's values apply today).
        if self._live_date is None and self._live_bullish is None and self.history:
            last_seeded_date = self.history[-1].date
            prev = self.history[-2] if date == last_seeded_date and len(self.history) >= 2 else self.history[-1]
            self._live_support = prev.support if prev else None
            self._live_bullish = prev.bullish if prev else None
            self._live_hold    = prev.hold    if prev else None

        # Date changed — compute futures from D's settled state, promote to new checkpoint
        if self._live_date is not None and date != self._live_date:
            warning = f"Date changed from {self._live_date} to {date} — previous bar committed as history"
            old_cp = self._checkpoint
            settled_ema5 = self.ema5.copy()
            settled_ema20 = self.ema20.copy()
            cp = _make_checkpoint(self)
            self._checkpoint = cp
            if len(old_cp.bars) >= _FUTURES_MIN and old_cp.cd_ema.seeded:
                self._live_support, self._live_bullish, self._live_hold = compute_futures(
                    ema5_pre=old_cp.ema5.copy(), ema20_pre=old_cp.ema20.copy(),
                    cd_pre=old_cp.cd_ema.value,
                    ema5_post=settled_ema5, ema20_post=settled_ema20,
                )
            else:
                self._live_support = None
                self._live_bullish = None
                self._live_hold = None
            self.bars_1m.clear()
            self.bars_5m.clear()

        # Restore from checkpoint — guarantees idempotency
        self.ema5 = cp.ema5.copy()
        self.ema20 = cp.ema20.copy()
        self.ema50 = cp.ema50.copy()
        self.sma10 = cp.sma10.copy()
        self.sma50 = cp.sma50.copy()
        self.rsi = cp.rsi.copy()
        self.cd_ema = cp.cd_ema.copy()
        self.bars = deque(cp.bars, maxlen=_CAPACITY)
        self._live_date = date

        # Apply the bar
        bar = Bar(date=date, close=close, open=open_, high=high, low=low, timestamp=timestamp)
        self.bars.append(bar)

        m = self.ema5.update(close)
        o = self.ema20.update(close)
        e50 = self.ema50.update(close)
        sma10_val = self.sma10.update(close)
        sma50_val = self.sma50.update(close)
        rsi_val = self.rsi.update(close)

        highs = deque(b.high for b in self.bars)
        hl = calc_hl(highs)
        avg = calc_avg(sma10_val, sma50_val)

        # Enrich bar and update-or-append to history (one entry per day)
        bar.hl = hl; bar.avg = avg; bar.ema5 = m; bar.ema20 = o; bar.ema50 = e50
        bar.rsi = rsi_val; bar.support = self._live_support; bar.bullish = self._live_bullish
        bar.hold = self._live_hold
        if self.history and self.history[-1].date == date:
            self.history[-1] = bar
        else:
            self.history.append(bar)

        # Intraday aggregation (1m and 5m)
        if timestamp:
            self._aggregate_intraday(self.bars_1m, 60, timestamp, close, open_, high, low)
            self._aggregate_intraday(self.bars_5m, 300, timestamp, close, open_, high, low)

        return TickerSnapshot(
            ticker=self.ticker, date=date, close=close, open=open_, high=high, low=low,
            hl=hl, avg=avg, ema5=m, ema20=o, ema50=e50, rsi=rsi_val,
            support=self._live_support, bullish=self._live_bullish, hold=self._live_hold,
            warning=warning,
            timestamp=timestamp,
        )

    def _aggregate_intraday(self, target_deque: deque, period_sec: int, ts: int, 
                             close: float, open_: float, high: float, low: float):
        """Aggregate per-second ticks into intraday candles."""
        period_ts = (ts // period_sec) * period_sec
        if target_deque and target_deque[-1].timestamp == period_ts:
            # Update current candle
            b = target_deque[-1]
            b.close = close
            b.high = max(b.high, high)
            b.low = min(b.low, low)
        else:
            # Start new candle
            target_deque.append(Bar(
                date="", close=close, open=open_, high=high, low=low, timestamp=period_ts
            ))


def _make_checkpoint(state: TickerState) -> _Checkpoint:
    return _Checkpoint(
        ema5=state.ema5.copy(),
        ema20=state.ema20.copy(),
        ema50=state.ema50.copy(),
        sma10=state.sma10.copy(),
        sma50=state.sma50.copy(),
        rsi=state.rsi.copy(),
        bars=deque(state.bars, maxlen=_CAPACITY),
        cd_ema=state.cd_ema.copy(),
        ema5_pre=state._futures_ema5_pre.copy() if state._futures_ema5_pre else None,
        ema20_pre=state._futures_ema20_pre.copy() if state._futures_ema20_pre else None,
        cd_pre=state._futures_cd_pre,
    )
