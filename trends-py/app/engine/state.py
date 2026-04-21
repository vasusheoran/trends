"""
Per-ticker rolling state manager.

Two modes:
  Commit mode  — no checkpoint set; update() permanently appends each bar.
                 Used during seeding. Indicators and futures computed normally.
  Live mode    — checkpoint set via commit(); update() restores from checkpoint
                 before applying the bar, making repeated PUTs idempotent.
                 When the date changes, previous live bar is promoted to a new
                 checkpoint automatically.
"""

from collections import deque
from dataclasses import dataclass, field
from typing import Optional

from app.engine.indicators import EMAState, SMAState, RSIState, calc_hl, calc_avg
from app.engine.futures import compute_futures
from app.models import TickerSnapshot

_CAPACITY = 101
_FUTURES_MIN = 100


@dataclass
class Bar:
    date: str
    close: float
    open: float
    high: float
    low: float


@dataclass
class _Checkpoint:
    """Frozen indicator state after all committed (historical) bars."""
    ema5: EMAState
    ema20: EMAState
    sma10: SMAState
    sma50: SMAState
    rsi: RSIState
    bars: deque
    cd: float


@dataclass
class TickerState:
    ticker: str

    bars: deque = field(default_factory=lambda: deque(maxlen=_CAPACITY))
    ema5: EMAState = field(default_factory=lambda: EMAState(period=5, decay=2/6))
    ema20: EMAState = field(default_factory=lambda: EMAState(period=20, decay=2/21))
    sma10: SMAState = field(default_factory=lambda: SMAState(period=10))
    sma50: SMAState = field(default_factory=lambda: SMAState(period=50))
    rsi: RSIState = field(default_factory=RSIState)
    cd: float = 0.0
    _last_futures_high: float = 0.0

    # Live-mode state (populated after commit())
    _checkpoint: Optional[_Checkpoint] = field(default=None, repr=False)
    _live_date: Optional[str] = field(default=None, repr=False)
    _live_support: Optional[float] = field(default=None, repr=False)
    _live_bullish: Optional[float] = field(default=None, repr=False)
    _live_cd: float = field(default=0.0, repr=False)

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
        self._live_cd = self.cd
        self._last_futures_high = 0.0

    def update(self, date: str, close: float, open_: float, high: float, low: float) -> TickerSnapshot:
        if self._checkpoint is not None:
            return self._update_live(date, close, open_, high, low)
        return self._update_commit(date, close, open_, high, low)

    def _update_commit(self, date, close, open_, high, low) -> TickerSnapshot:
        """Permanently append a bar (commit/seed mode)."""
        bar = Bar(date=date, close=close, open=open_, high=high, low=low)
        self.bars.append(bar)

        ema5_pre = self.ema5.copy()
        ema20_pre = self.ema20.copy()

        m = self.ema5.update(close)
        o = self.ema20.update(close)
        sma10_val = self.sma10.update(close)
        sma50_val = self.sma50.update(close)
        rsi_val = self.rsi.update(close)

        highs = deque(b.high for b in self.bars)
        hl = calc_hl(highs)
        avg = calc_avg(sma10_val, sma50_val)

        support = None
        bullish = None
        if len(self.bars) >= _FUTURES_MIN and (high > self._last_futures_high or self._last_futures_high == 0.0):
            self._last_futures_high = high
            if m is not None and o is not None:
                support, bullish, self.cd = compute_futures(
                    ema5=ema5_pre, ema20=ema20_pre, old_cd=self.cd,
                )

        return TickerSnapshot(
            ticker=self.ticker, date=date, close=close, open=open_, high=high, low=low,
            hl=hl, avg=avg, ema5=m, ema20=o, rsi=rsi_val,
            support=support, bullish=bullish,
        )

    def _update_live(self, date, close, open_, high, low) -> TickerSnapshot:
        """
        Apply bar on top of checkpoint (live mode).
        Idempotent: same inputs always produce the same snapshot.
        When date changes, the previous live bar is promoted to a new checkpoint.
        """
        cp = self._checkpoint
        warning = None

        # Date changed — promote previous live state to new checkpoint
        if self._live_date is not None and date != self._live_date:
            warning = f"Date changed from {self._live_date} to {date} — previous bar committed as history"
            cp = _make_checkpoint(self)
            cp.cd = self._live_cd
            self._checkpoint = cp
            self._live_support = None
            self._live_bullish = None
            self._last_futures_high = 0.0

        # Restore from checkpoint — guarantees idempotency
        self.ema5 = cp.ema5.copy()
        self.ema20 = cp.ema20.copy()
        self.sma10 = cp.sma10.copy()
        self.sma50 = cp.sma50.copy()
        self.rsi = cp.rsi.copy()
        self.bars = deque(cp.bars, maxlen=_CAPACITY)
        self._live_date = date

        # Apply the bar
        bar = Bar(date=date, close=close, open=open_, high=high, low=low)
        self.bars.append(bar)

        ema5_pre = self.ema5.copy()
        ema20_pre = self.ema20.copy()

        m = self.ema5.update(close)
        o = self.ema20.update(close)
        sma10_val = self.sma10.update(close)
        sma50_val = self.sma50.update(close)
        rsi_val = self.rsi.update(close)

        highs = deque(b.high for b in self.bars)
        hl = calc_hl(highs)
        avg = calc_avg(sma10_val, sma50_val)

        # Futures — recompute only when high changes
        if len(self.bars) >= _FUTURES_MIN and high != self._last_futures_high:
            self._last_futures_high = high
            if m is not None and o is not None:
                self._live_support, self._live_bullish, self._live_cd = compute_futures(
                    ema5=ema5_pre, ema20=ema20_pre, old_cd=cp.cd,
                )

        return TickerSnapshot(
            ticker=self.ticker, date=date, close=close, open=open_, high=high, low=low,
            hl=hl, avg=avg, ema5=m, ema20=o, rsi=rsi_val,
            support=self._live_support, bullish=self._live_bullish,
            warning=warning,
        )


def _make_checkpoint(state: TickerState) -> _Checkpoint:
    return _Checkpoint(
        ema5=state.ema5.copy(),
        ema20=state.ema20.copy(),
        sma10=state.sma10.copy(),
        sma50=state.sma50.copy(),
        rsi=state.rsi.copy(),
        bars=deque(state.bars, maxlen=_CAPACITY),
        cd=state.cd,
    )
