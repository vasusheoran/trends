"""
TimescaleDB via asyncpg.
Schema: ticker_ticks hypertable, partitioned by ts.
"""

import asyncpg
from typing import Optional

_pool: Optional[asyncpg.Pool] = None


async def init_pool(database_url: str) -> None:
    global _pool
    _pool = await asyncpg.create_pool(database_url)
    await _create_schema()


async def close_pool() -> None:
    if _pool:
        await _pool.close()


async def _create_schema() -> None:
    async with _pool.acquire() as conn:
        await conn.execute("""
            CREATE TABLE IF NOT EXISTS ticker_ticks (
                ticker  TEXT        NOT NULL,
                ts      TIMESTAMPTZ NOT NULL,
                date    TEXT        NOT NULL,
                close   DOUBLE PRECISION,
                open    DOUBLE PRECISION,
                high    DOUBLE PRECISION,
                low     DOUBLE PRECISION,
                ema5    DOUBLE PRECISION,
                ema20   DOUBLE PRECISION,
                hl      DOUBLE PRECISION,
                avg     DOUBLE PRECISION,
                support DOUBLE PRECISION,
                bullish DOUBLE PRECISION,
                rsi     DOUBLE PRECISION,
                PRIMARY KEY (ticker, ts)
            );
        """)
        # Create hypertable if TimescaleDB extension is available
        try:
            await conn.execute("""
                SELECT create_hypertable(
                    'ticker_ticks', 'ts',
                    if_not_exists => TRUE,
                    migrate_data  => TRUE
                );
            """)
        except Exception:
            pass  # TimescaleDB extension not installed — plain Postgres still works


async def get_row_count(ticker: str) -> int:
    async with _pool.acquire() as conn:
        row = await conn.fetchrow(
            "SELECT COUNT(*) as cnt FROM ticker_ticks WHERE ticker = $1", ticker
        )
        return row["cnt"]


async def upsert_tick(snapshot) -> None:
    """Upsert a TickerSnapshot into ticker_ticks."""
    from datetime import datetime, timezone
    async with _pool.acquire() as conn:
        await conn.execute("""
            INSERT INTO ticker_ticks
                (ticker, ts, date, close, open, high, low, ema5, ema20, hl, avg, support, bullish, rsi)
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
            ON CONFLICT (ticker, ts) DO UPDATE SET
                close=EXCLUDED.close, open=EXCLUDED.open,
                high=EXCLUDED.high,  low=EXCLUDED.low,
                ema5=EXCLUDED.ema5,  ema20=EXCLUDED.ema20,
                hl=EXCLUDED.hl,      avg=EXCLUDED.avg,
                support=EXCLUDED.support, bullish=EXCLUDED.bullish, rsi=EXCLUDED.rsi
        """,
            snapshot.ticker,
            datetime.now(timezone.utc),
            snapshot.date,
            snapshot.close, snapshot.open, snapshot.high, snapshot.low,
            snapshot.ema5, snapshot.ema20,
            snapshot.hl, snapshot.avg,
            snapshot.support, snapshot.bullish, snapshot.rsi,
        )


async def load_bars(ticker: str, limit: int = 101) -> list[dict]:
    """Load most recent bars for EMA seeding on startup."""
    async with _pool.acquire() as conn:
        rows = await conn.fetch("""
            SELECT date, close, open, high, low
            FROM ticker_ticks
            WHERE ticker = $1
            ORDER BY ts ASC
            LIMIT $2
        """, ticker, limit)
        return [dict(r) for r in rows]
