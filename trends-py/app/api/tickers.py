from collections import deque
from typing import Optional

from fastapi import APIRouter, HTTPException, Query
from app.registry import delete_state, list_tickers, _states
from app.engine.futures import _ce2, get_support, _CD_DECAY, compute_futures
from app.engine.indicators import calc_hl, calc_avg

_MONTHS = {'Jan':1,'Feb':2,'Mar':3,'Apr':4,'May':5,'Jun':6,
           'Jul':7,'Aug':8,'Sep':9,'Oct':10,'Nov':11,'Dec':12}

def _date_to_iso(d: str) -> Optional[str]:
    """Convert 'DD-Mon-YYYY' (or 'DD-Mon-YY') → 'YYYY-MM-DD'. Returns None if unparseable."""
    try:
        parts = d.split('-')
        if len(parts) == 3 and parts[1] in _MONTHS:
            yr = int(parts[2])
            if yr < 100:                          # 2-digit year: 03 → 2003
                yr += 2000 if yr <= 30 else 1900
            return f"{yr:04d}-{_MONTHS[parts[1]]:02d}-{int(parts[0]):02d}"
    except Exception:
        pass
    return None

router = APIRouter()

_USER_PREFS = {"theme": "dark"}


@router.get("/api/preferences")
async def get_preferences():
    """Return stored user preferences."""
    return _USER_PREFS


@router.post("/api/preferences")
async def update_preferences(prefs: dict):
    """Update user preferences (in-memory)."""
    _USER_PREFS.update(prefs)
    return _USER_PREFS


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
async def get_ticker_history(ticker: str, year: Optional[int] = Query(default=None)):
    """
    Return enriched bar history for a ticker.
    Optional ?year=YYYY filters to Jan–Dec of that year ±3 months (Oct y-1 to Mar y+1).
    Always returns the full list of available years for the picker.
    """
    ticker = ticker.lower()
    if ticker not in _states:
        raise HTTPException(status_code=404, detail=f"Ticker '{ticker}' not found")

    state = _states[ticker]

    bars = []
    years_set: set[int] = set()

    for b in state.history:
        iso = _date_to_iso(b.date)
        if iso:
            y, mo = int(iso[:4]), int(iso[5:7])
            years_set.add(y)
            if year is not None:
                in_range = (
                    (y == year - 1 and mo >= 10) or
                    y == year or
                    (y == year + 1 and mo <= 3)
                )
                if not in_range:
                    continue

        bars.append({
            "date":    b.date,
            "time":    iso,
            "open":    round(b.open, 2),
            "high":    round(b.high, 2),
            "low":     round(b.low, 2),
            "close":   round(b.close, 2),
            "hl":      round(b.hl, 2)      if b.hl      is not None else None,
            "avg":     round(b.avg, 2)     if b.avg     is not None else None,
            "ema5":    round(b.ema5, 2)    if b.ema5    is not None else None,
            "ema20":   round(b.ema20, 2)   if b.ema20   is not None else None,
            "rsi":     round(b.rsi, 1)     if b.rsi     is not None else None,
            "support": round(b.support, 2) if b.support is not None else None,
            "bullish": round(b.bullish, 2) if b.bullish is not None else None,
        })

    years = sorted(years_set, reverse=True)
    return {"ticker": ticker.upper(), "history": bars, "years": years}


@router.get("/api/intraday/{ticker}")
async def get_ticker_intraday(ticker: str, tf: str = "1m"):
    """Return intraday bars (1m or 5m) for the current day."""
    ticker = ticker.lower()
    if ticker not in _states:
        raise HTTPException(status_code=404, detail=f"Ticker '{ticker}' not found")
    
    state = _states[ticker]
    bars = state.bars_1m if tf == "1m" else state.bars_5m
    
    return [
        {
            "time": b.timestamp,
            "open": b.open,
            "high": b.high,
            "low": b.low,
            "close": b.close,
        } for b in bars
    ]


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
