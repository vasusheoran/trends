"""
SSE endpoint — GET /api/stream/{ticker}
Pushes TickerSnapshot JSON on every tick update.
"""

import asyncio
import json
from fastapi import APIRouter
from fastapi.responses import StreamingResponse
from app.registry import subscribe, unsubscribe

router = APIRouter()


@router.get("/api/stream/{ticker}")
async def stream_ticker(ticker: str):
    queue: asyncio.Queue = asyncio.Queue()
    subscribe(ticker, queue)

    async def event_generator():
        try:
            while True:
                snapshot = await queue.get()
                data = json.dumps(snapshot.model_dump())
                yield f"data: {data}\n\n"
        except asyncio.CancelledError:
            pass
        finally:
            unsubscribe(ticker, queue)

    return StreamingResponse(
        event_generator(),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "X-Accel-Buffering": "no",
        },
    )
