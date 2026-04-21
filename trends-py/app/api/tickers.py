from fastapi import APIRouter, HTTPException
from app.registry import delete_state, list_tickers

router = APIRouter()


@router.get("/api/tickers")
async def get_tickers():
    """List all active tickers currently loaded in memory."""
    return {"tickers": list_tickers()}


@router.delete("/api/tickers/{ticker}")
async def delete_ticker(ticker: str):
    """
    Remove a ticker from memory. All state and history is discarded.
    Any open SSE streams for this ticker will receive no further events.
    Returns 404 if the ticker was not loaded.
    """
    if not delete_state(ticker):
        raise HTTPException(status_code=404, detail=f"Ticker '{ticker}' not found")
    return {"ticker": ticker, "status": "deleted"}
