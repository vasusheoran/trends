"""
Indicator calculations — all formulas verified against main Excel.
Source of truth: data/Nifty-17-04-2026.xlsx
See docs/FIELDS.md for full field reference.
"""

from collections import deque
from dataclasses import dataclass, field
from typing import Optional


TOLERANCE = 0.001


@dataclass
class EMAState:
    """Single EMA series state."""
    period: int
    decay: float
    value: float = 0.0
    seeded: bool = False
    _seed_buf: list = field(default_factory=list)

    def update(self, price: float) -> Optional[float]:
        if self.seeded:
            self.value = self.decay * (price - self.value) + self.value
            return self.value
        self._seed_buf.append(price)
        if len(self._seed_buf) == self.period:
            self.value = sum(self._seed_buf) / self.period
            self.seeded = True
            return self.value
        return None

    def copy(self) -> "EMAState":
        s = EMAState(period=self.period, decay=self.decay, value=self.value, seeded=self.seeded)
        s._seed_buf = self._seed_buf.copy()
        return s


@dataclass
class SMAState:
    """Rolling SMA using a fixed-size deque."""
    period: int
    _buf: deque = field(default_factory=deque)
    _total: float = 0.0

    def update(self, price: float) -> Optional[float]:
        self._buf.append(price)
        self._total += price
        if len(self._buf) > self.period:
            self._total -= self._buf.popleft()
        if len(self._buf) == self.period:
            return self._total / self.period
        return None

    def copy(self) -> "SMAState":
        s = SMAState(period=self.period)
        s._buf = self._buf.copy()
        s._total = self._total
        return s


@dataclass
class RSIState:
    """RSI(14) using Wilder smoothing — matches Excel col CB formula."""
    period: int = 14
    avg_gain: float = 0.0
    avg_loss: float = 0.0
    seeded: bool = False
    prev_close: Optional[float] = None
    _seed_gains: list = field(default_factory=list)
    _seed_losses: list = field(default_factory=list)

    def update(self, close: float) -> Optional[float]:
        if self.prev_close is None:
            self.prev_close = close
            return None

        change = close - self.prev_close
        self.prev_close = close
        gain = max(change, 0.0)
        loss = abs(min(change, 0.0))

        if not self.seeded:
            self._seed_gains.append(gain)
            self._seed_losses.append(loss)
            if len(self._seed_gains) == self.period:
                self.avg_gain = sum(self._seed_gains) / self.period
                self.avg_loss = sum(self._seed_losses) / self.period
                self.seeded = True
                return self._rsi()
            return None

        self.avg_gain = (self.avg_gain * (self.period - 1) + gain) / self.period
        self.avg_loss = (self.avg_loss * (self.period - 1) + loss) / self.period
        return self._rsi()

    def _rsi(self) -> float:
        if self.avg_loss == 0:
            return 100.0
        return 100.0 - (100.0 / (1.0 + self.avg_gain / self.avg_loss))

    def copy(self) -> "RSIState":
        s = RSIState(period=self.period, avg_gain=self.avg_gain, avg_loss=self.avg_loss,
                     seeded=self.seeded, prev_close=self.prev_close)
        s._seed_gains = self._seed_gains.copy()
        s._seed_losses = self._seed_losses.copy()
        return s


def calc_hl(highs: deque) -> Optional[float]:
    """H/L = min of previous 3 highs. Needs at least 4 bars (index >= 3)."""
    if len(highs) < 4:
        return None
    h = list(highs)
    return min(h[-2], h[-3], h[-4])


def calc_avg(sma10: Optional[float], sma50: Optional[float]) -> Optional[float]:
    """AVG = (SMA10+SMA50)/2 with correction. Source: Final-bullish-ce.xlsx col AR.
    Denominator in the correction term is Sum (sma10+sma50), not A — matches Excel operator precedence."""
    if sma10 is None or sma50 is None:
        return None
    Sum = sma10 + sma50
    A = Sum / 2
    inner = A - A * 0.01
    inner2 = inner * 0.025
    inner3 = (inner + inner2 + A) / 2
    result = A - (A * ((A - inner3) / Sum / 2 * 100 / 2) / 100)
    return result
