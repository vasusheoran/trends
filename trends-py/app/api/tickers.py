from fastapi import APIRouter, HTTPException
from app.registry import delete_state, list_tickers, _states
from app.engine.futures import _ce2, get_support, _CD_DECAY, compute_futures

router = APIRouter()


@router.get("/api/tickers")
async def get_tickers():
    """List all active tickers currently loaded in memory."""
    return {"tickers": list_tickers()}


@router.get("/api/state/{ticker}")
async def get_ticker_state(ticker: str):
    """
    Return current indicator and futures values for a seeded ticker.
    Uses the checkpoint EMA state (end of seeded history) plus the last live bar if present.
    """
    ticker = ticker.lower()
    if ticker not in _states:
        raise HTTPException(status_code=404, detail=f"Ticker '{ticker}' not found")

    state = _states[ticker]
    cp = state._checkpoint

    if cp is None:
        raise HTTPException(status_code=400, detail="Ticker is still in seed/commit mode — call commit first")

    e5 = cp.ema5.copy()
    e20 = cp.ema20.copy()

    result = {
        "ticker": ticker,
        "bars_seeded": len(cp.bars),
        "ema5": round(e5.value, 4) if e5.value else None,
        "ema20": round(e20.value, 4) if e20.value else None,
        "cd": round(cp.cd_ema.value, 4) if cp.cd_ema.seeded else None,
        "support": None,
        "bullish": None,
        "ce2": None,
        "cd_curr": None,
    }

    if cp.cd_ema.seeded and e5.value and e20.value:
        cd_pre = cp.cd_ema.value
        ce2 = _ce2(e5, e20)
        cd_curr = _CD_DECAY * (ce2 - cd_pre) + cd_pre
        result["ce2"] = round(ce2, 4)
        result["cd_curr"] = round(cd_curr, 4)

        support, bullish = compute_futures(e5, e20, cd_pre)
        result["support"] = round(support, 4) if support else None
        result["bullish"] = round(bullish, 4) if bullish else None

    # Include last live values if available
    if state._live_support is not None:
        result["live_support"] = round(state._live_support, 4)
    if state._live_bullish is not None:
        result["live_bullish"] = round(state._live_bullish, 4)

    return result


@router.get("/api/history/{ticker}")
async def get_ticker_history(ticker: str):
    """Return the last 101 bars of history for a ticker."""
    ticker = ticker.lower()
    if ticker not in _states:
        raise HTTPException(status_code=404, detail=f"Ticker '{ticker}' not found")

    state = _states[ticker]
    bars = []
    for b in state.bars:
        bars.append({
            "time": b.date,  # 'time' is used by many charting libs
            "open": round(b.open, 2),
            "high": round(b.high, 2),
            "low": round(b.low, 2),
            "close": round(b.close, 2),
        })
    return {"ticker": ticker.upper(), "history": bars}


@router.delete("/api/tickers/{ticker}")
async def delete_ticker(ticker: str):
    """
    Remove a ticker from memory. All state and history is discarded.
    Returns 404 if the ticker was not loaded.
    """
    ticker = ticker.lower()
    if not delete_state(ticker):
        raise HTTPException(status_code=404, detail=f"Ticker '{ticker}' not found")
    return {"ticker": ticker, "status": "deleted"}
