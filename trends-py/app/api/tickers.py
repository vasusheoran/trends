from collections import deque

from fastapi import APIRouter, HTTPException
from app.registry import delete_state, list_tickers, _states
from app.engine.futures import _ce2, get_support, _CD_DECAY, compute_futures
from app.engine.indicators import calc_hl, calc_avg

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

    # Last bar OHLC
    last_bar = cp.bars[-1] if cp.bars else None
    highs = deque(b.high for b in cp.bars)
    hl = calc_hl(highs)
    sma10_val = cp.sma10._total / len(cp.sma10._buf) if len(cp.sma10._buf) == cp.sma10.period else None
    sma50_val = cp.sma50._total / len(cp.sma50._buf) if len(cp.sma50._buf) == cp.sma50.period else None
    avg = calc_avg(sma10_val, sma50_val)
    rsi_val = cp.rsi._rsi() if cp.rsi.seeded else None

    result = {
        "ticker": ticker,
        "date": last_bar.date if last_bar else None,
        "close": last_bar.close if last_bar else None,
        "open": last_bar.open if last_bar else None,
        "high": last_bar.high if last_bar else None,
        "low": last_bar.low if last_bar else None,
        "hl": round(hl, 4) if hl is not None else None,
        "avg": round(avg, 4) if avg is not None else None,
        "ema5": round(e5.value, 4) if e5.value else None,
        "ema20": round(e20.value, 4) if e20.value else None,
        "rsi": round(rsi_val, 2) if rsi_val is not None else None,
        "bars_seeded": len(cp.bars),
        "cd": round(cp.cd_ema.value, 4) if cp.cd_ema.seeded else None,
        "support": None,
        "bullish": None,
        "ce2": None,
        "cd_curr": None,
    }

    if cp.cd_ema.seeded and e5.value and e20.value:
        ce2 = _ce2(e5, e20)
        cd_curr = _CD_DECAY * (ce2 - cp.cd_ema.value) + cp.cd_ema.value
        result["ce2"] = round(ce2, 4)
        result["cd_curr"] = round(cd_curr, 4)

        # Support uses pre-last-bar EMA; Bullish uses post-last-bar EMA (matches Go timing)
        if cp.ema5_pre is not None and cp.cd_pre is not None:
            support, bullish = compute_futures(
                cp.ema5_pre.copy(), cp.ema20_pre.copy(), cp.cd_pre,
                ema5_post=e5, ema20_post=e20,
            )
        else:
            support, bullish = compute_futures(e5, e20, cp.cd_ema.value)
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
    bars = [
        {
            "time": b.date,
            "open": round(b.open, 2),
            "high": round(b.high, 2),
            "low": round(b.low, 2),
            "close": round(b.close, 2),
        }
        for b in state.history
    ]
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
