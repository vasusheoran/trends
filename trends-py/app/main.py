import asyncio
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.responses import RedirectResponse
from fastapi.staticfiles import StaticFiles

from app.config import settings
from app.db.timescale import init_pool, close_pool
from app.db.seed import seed_state
from app.registry import get_state
from app.ingest.webhook import router as webhook_router
from app.ingest.seed_upload import router as seed_router
from app.api.stream import router as stream_router
from app.api.health import router as health_router
from app.api.tickers import router as tickers_router
from app.api.debug import router as debug_router

logging.basicConfig(level=logging.INFO)
log = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    await init_pool(settings.database_url)

    # Seed all tickers stored in DB; fall back to Excel for nifty if DB is empty
    from app.db.timescale import get_all_tickers
    db_tickers = await get_all_tickers()
    targets = db_tickers if db_tickers else ["nifty"]
    for ticker_name in targets:
        try:
            state = get_state(ticker_name)
            source = await seed_state(state, ticker_name, settings.excel_seed_path)
            log.info("Startup seed: %s loaded from %s (%d bars)", ticker_name, source, len(state.bars))
        except Exception as e:
            log.warning("Startup seed skipped for %s: %s", ticker_name, e)

    # Zerodha feed (Track 2) — only if token is configured
    zerodha_task = None
    if settings.zerodha_access_token:
        from app.ingest.zerodha import start_zerodha_feed
        zerodha_task = asyncio.create_task(
            start_zerodha_feed(
                api_key=settings.zerodha_api_key,
                access_token=settings.zerodha_access_token,
                instrument_token=settings.zerodha_nifty_token,
            )
        )
        log.info("Zerodha feed started")
    else:
        log.info("No ZERODHA_ACCESS_TOKEN — using webhook track only")

    yield

    # Shutdown
    if zerodha_task:
        zerodha_task.cancel()
    await close_pool()


app = FastAPI(title="Trends", lifespan=lifespan)

# Mount static files
app.mount("/static", StaticFiles(directory="static"), name="static")

app.include_router(webhook_router)
app.include_router(seed_router)
app.include_router(stream_router)
app.include_router(health_router)
app.include_router(tickers_router)
app.include_router(debug_router)


@app.get("/", include_in_schema=False)
async def root():
    return RedirectResponse(url="/static/index.html")
