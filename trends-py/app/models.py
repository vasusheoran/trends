from pydantic import BaseModel
from typing import Optional


class TickerPayload(BaseModel):
    """Incoming webhook payload — matches VB PostTrend shape."""
    date: str
    close: float
    open: float
    high: float
    low: float
    timestamp: Optional[int] = None # Unix timestamp (seconds)


class TickerSnapshot(BaseModel):
    """Computed state emitted via SSE and stored in TimescaleDB."""
    ticker: str
    date: str
    timestamp: Optional[int] = None

    close: float
    open: float
    high: float
    low: float

    # Indicators
    hl: Optional[float] = None        # H/L  — min of prev 3 highs
    avg: Optional[float] = None       # AVG  — (SMA10 + SMA50) / 2 with correction
    ema5: Optional[float] = None      # EMA-5  (decay 2/6)
    ema20: Optional[float] = None     # EMA-20 (decay 2/21)
    ema50: Optional[float] = None     # EMA-50 (decay 2/51)
    rsi: Optional[float] = None       # RSI(14) Wilder

    # Futures (binary-search projections)
    support: Optional[float] = None   # CC — lower support target
    bullish: Optional[float] = None   # BR — upper bullish target
    hold: Optional[float] = None      # Hold — min D+1 close to preserve today's Bullish

    # Set when the PUT date differs from the current live date
    warning: Optional[str] = None
