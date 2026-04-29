"""
Track 1 — Manual PUT /api/update/{ticker}
Matches VB script PostTrend payload.
"""

import time
from fastapi import APIRouter, HTTPException, Query
from app.models import TickerPayload, TickerSnapshot
from app.registry import get_state, publish

router = APIRouter()


@router.put("/api/update/{ticker}")
async def update_ticker(
    ticker: str,
    payload: TickerPayload,
    force: bool = Query(False, description="Force futures recompute even if already computed today"),
) -> TickerSnapshot:
    state = get_state(ticker)
    ts = payload.timestamp or int(time.time())
    snapshot = state.update(
        date=payload.date,
        close=payload.close,
        open_=payload.open,
        high=payload.high,
        low=payload.low,
        force=force,
        timestamp=ts,
    )
    await publish(ticker, snapshot)
    return snapshot
