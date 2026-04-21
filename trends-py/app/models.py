from pydantic import BaseModel
from typing import Optional


class TickerPayload(BaseModel):
    """Incoming webhook payload — matches VB PostTrend shape."""
    date: str
    close: float
    open: float
    high: float
    low: float


class TickerSnapshot(BaseModel):
    """Computed state emitted via SSE and stored in TimescaleDB."""
    ticker: str
    date: str

    close: float
    open: float
    high: float
    low: float

    # Indicators
    hl: Optional[float] = None        # H/L  — min of prev 3 highs
    avg: Optional[float] = None       # AVG  — (SMA10 + SMA50) / 2 with correction
    ema5: Optional[float] = None      # EMA-5
    ema20: Optional[float] = None     # EMA-20
    ema50: Optional[float] = None     # EMA-50
    rsi: Optional[float] = None       # RSI(14) Wilder

    # Futures (binary-search projections)
    support: Optional[float] = None   # CE  — lower support target
    bullish: Optional[float] = None   # CC/CH — upper bullish target
