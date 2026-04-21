"""
Track 2 — Zerodha KiteTicker WebSocket.
Only active when ZERODHA_ACCESS_TOKEN is set.
On error/disconnect, falls back silently to webhook Track 1.
"""

import asyncio
import logging
from datetime import datetime

log = logging.getLogger(__name__)


async def start_zerodha_feed(api_key: str, access_token: str, instrument_token: int) -> None:
    """
    Connect to Zerodha KiteTicker and pipe ticks into the state registry.
    Runs as a background asyncio task.
    Requires: uv add kiteconnect
    """
    try:
        from kiteconnect import KiteTicker  # type: ignore
    except ImportError:
        log.warning("kiteconnect package not installed — Zerodha feed disabled. Run: uv add kiteconnect")
        return

    from app.registry import get_state, publish

    loop = asyncio.get_event_loop()

    def on_ticks(ws, ticks):
        for tick in ticks:
            if tick["instrument_token"] != instrument_token:
                continue
            ohlc = tick.get("ohlc", {})
            close = tick.get("last_price")
            open_ = ohlc.get("open")
            high  = ohlc.get("high")
            low   = ohlc.get("low")
            if None in (close, open_, high, low):
                continue
            date = datetime.now().strftime("%d-%b-%Y")
            # Schedule the coroutine back on the event loop
            asyncio.run_coroutine_threadsafe(_ingest("nifty", date, close, open_, high, low), loop)

    async def _ingest(ticker, date, close, open_, high, low):
        state = get_state(ticker)
        snapshot = state.update(date=date, close=close, open_=open_, high=high, low=low)
        await publish(ticker, snapshot)

    def on_connect(ws, response):
        ws.subscribe([instrument_token])
        ws.set_mode(ws.MODE_FULL, [instrument_token])
        log.info("Zerodha KiteTicker connected, subscribed token=%d", instrument_token)

    def on_error(ws, code, reason):
        log.error("Zerodha error %s: %s", code, reason)

    def on_close(ws, code, reason):
        log.warning("Zerodha connection closed %s: %s — falling back to webhook", code, reason)

    kws = KiteTicker(api_key, access_token)
    kws.on_ticks = on_ticks
    kws.on_connect = on_connect
    kws.on_error = on_error
    kws.on_close = on_close

    # KiteTicker.connect() is blocking — run in a thread
    await loop.run_in_executor(None, kws.connect)
