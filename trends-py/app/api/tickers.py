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
    Reflects the most recent update() call — shows live close with day-pinned futures.
    """
    ticker = ticker.lower()
    if ticker not in _states:
        raise HTTPException(status_code=404, detail=f"Ticker '{ticker}' not found")

    state = _states[ticker]

    if not state.history:
        raise HTTPException(status_code=400, detail="Ticker has no data — send at least one tick first")

    live = state.history[-1]

    result = {
        "ticker": ticker,
        "date":  live.date,
        "close": live.close,
        "open":  live.open,
        "high":  live.high,
        "low":   live.low,
        "hl":    round(live.hl, 4)    if live.hl    is not None else None,
        "avg":   round(live.avg, 4)   if live.avg   is not None else None,
        "ema5":  round(live.ema5, 4)  if live.ema5  is not None else None,
        "ema20": round(live.ema20, 4) if live.ema20 is not None else None,
        "ema50": round(live.ema50, 4) if live.ema50 is not None else None,
        "rsi":   round(live.rsi, 2)   if live.rsi   is not None else None,
        "support": round(live.support, 4) if live.support is not None else None,
        "bullish": round(live.bullish, 4) if live.bullish is not None else None,
        "hold":    round(live.hold, 4)    if live.hold    is not None else None,
        "bars_seeded": len(state.bars),
        "cd": round(state.cd_ema.value, 4) if state.cd_ema.seeded else None,
        "ce2": None,
        "cd_curr": None,
    }

    if state.cd_ema.seeded and state.ema5.value and state.ema20.value:
        e5 = state.ema5.copy()
        e20 = state.ema20.copy()
        ce2 = _ce2(e5, e20)
        cd_curr = _CD_DECAY * (ce2 - state.cd_ema.value) + state.cd_ema.value
        result["ce2"] = round(ce2, 4)
        result["cd_curr"] = round(cd_curr, 4)

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
            "ema50":   round(b.ema50, 2)   if b.ema50   is not None else None,
            "rsi":     round(b.rsi, 1)     if b.rsi     is not None else None,
            "hold":    round(b.hold, 2)    if b.hold    is not None else None,
            "support": round(b.support, 2) if b.support is not None else None,
            "bullish": round(b.bullish, 2) if b.bullish is not None else None,
        })

    years = sorted(years_set, reverse=True)
    return {"ticker": ticker.upper(), "history": bars, "years": years}


@router.get("/api/intraday/{ticker}")
async def get_ticker_intraday(ticker: str, tf: str = "1m", date: Optional[str] = None, raw: bool = False):
    """Return intraday OHLC bars (1m or 5m) plus list of available days.

    If `date` is provided (DD-Mon-YYYY), bars come from DB ticks for that day.
    If `raw=true`, returns unaggregated per-second ticks with all indicator columns.
    Otherwise returns aggregated candles (1m or 5m). Live/today bars come from memory.
    """
    ticker = ticker.lower()
    if ticker not in _states:
        raise HTTPException(status_code=404, detail=f"Ticker '{ticker}' not found")

    from app.db import timescale as ts_db

    days: list[str] = []
    try:
        days = await ts_db.get_available_days(ticker)
    except Exception:
        pass

    period_sec = 60 if tf == "1m" else 300

    def _to_raw_bar(r: dict) -> dict:
        return {
            "time":    r["ts_unix"],
            "open":    r["open"],    "high":    r["high"],
            "low":     r["low"],     "close":   r["close"],
            "ema5":    r.get("ema5"),  "ema20": r.get("ema20"),
            "ema50":   r.get("ema50"), "hl":    r.get("hl"),
            "avg":     r.get("avg"),   "support": r.get("support"),
            "rsi":     r.get("rsi"),
        }

    fetch_date = date or (days[0] if raw and days else None)

    if fetch_date:
        try:
            tick_rows = await ts_db.get_ticks_for_day(ticker, fetch_date)
            if raw:
                bars = [_to_raw_bar(r) for r in tick_rows]
            else:
                bars = ts_db.aggregate_ticks(tick_rows, period_sec)
        except Exception:
            bars = []
    else:
        state = _states[ticker]
        mem_bars = state.bars_1m if tf == "1m" else state.bars_5m
        bars = [
            {"time": b.timestamp, "open": b.open, "high": b.high, "low": b.low, "close": b.close}
            for b in mem_bars
        ]

    return {"bars": bars, "days": days}


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
