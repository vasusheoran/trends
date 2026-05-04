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
                ema50   DOUBLE PRECISION,
                hold    DOUBLE PRECISION,
                PRIMARY KEY (ticker, ts)
            );
        """)
        # Add columns that may be missing in existing tables
        for col, typedef in [("ema50", "DOUBLE PRECISION"), ("hold", "DOUBLE PRECISION")]:
            try:
                await conn.execute(
                    f"ALTER TABLE ticker_ticks ADD COLUMN IF NOT EXISTS {col} {typedef};"
                )
            except Exception:
                pass
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


async def get_all_tickers() -> list[str]:
    """Return all distinct ticker names stored in DB."""
    async with _pool.acquire() as conn:
        rows = await conn.fetch(
            "SELECT DISTINCT ticker FROM ticker_ticks ORDER BY ticker"
        )
    return [r["ticker"] for r in rows]


async def get_row_count(ticker: str) -> int:
    async with _pool.acquire() as conn:
        row = await conn.fetchrow(
            "SELECT COUNT(*) as cnt FROM ticker_ticks WHERE ticker = $1", ticker
        )
        return row["cnt"]


async def upsert_tick(snapshot) -> None:
    """Upsert a TickerSnapshot into ticker_ticks."""
    from datetime import datetime, timezone
    ts = datetime.fromtimestamp(snapshot.timestamp, timezone.utc) if snapshot.timestamp else datetime.now(timezone.utc)
    async with _pool.acquire() as conn:
        await conn.execute("""
            INSERT INTO ticker_ticks
                (ticker, ts, date, close, open, high, low, ema5, ema20, hl, avg, support, bullish, rsi, ema50, hold)
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
            ON CONFLICT (ticker, ts) DO UPDATE SET
                close=EXCLUDED.close, open=EXCLUDED.open,
                high=EXCLUDED.high,  low=EXCLUDED.low,
                ema5=EXCLUDED.ema5,  ema20=EXCLUDED.ema20,
                hl=EXCLUDED.hl,      avg=EXCLUDED.avg,
                support=EXCLUDED.support, bullish=EXCLUDED.bullish, rsi=EXCLUDED.rsi,
                ema50=EXCLUDED.ema50, hold=EXCLUDED.hold
        """,
            snapshot.ticker,
            ts,
            snapshot.date,
            snapshot.close, snapshot.open, snapshot.high, snapshot.low,
            snapshot.ema5, snapshot.ema20,
            snapshot.hl, snapshot.avg,
            snapshot.support, snapshot.bullish, snapshot.rsi,
            snapshot.ema50, snapshot.hold,
        )


async def load_bars(ticker: str, limit: int = 101) -> list[dict]:
    """Load most recent distinct daily OHLC bars for EMA seeding on startup.

    Aggregates per-second ticks into one row per date so that seeding always
    uses daily bars regardless of whether the table contains tick or EOD data.
    """
    async with _pool.acquire() as conn:
        rows = await conn.fetch("""
            SELECT
                date,
                (array_agg(open  ORDER BY ts ASC))[1]  AS open,
                MAX(high)                               AS high,
                MIN(low)                                AS low,
                (array_agg(close ORDER BY ts DESC))[1]  AS close
            FROM ticker_ticks
            WHERE ticker = $1
            GROUP BY date
            ORDER BY MAX(ts) ASC
            LIMIT $2
        """, ticker, limit)
        return [dict(r) for r in rows]


async def get_available_days(ticker: str) -> list[str]:
    """Return distinct trading days (DD-Mon-YYYY) that have tick data, newest first."""
    from datetime import datetime
    async with _pool.acquire() as conn:
        rows = await conn.fetch(
            "SELECT DISTINCT date FROM ticker_ticks WHERE ticker = $1", ticker
        )
    days = [r["date"] for r in rows]

    def _parse(d: str):
        try:
            return datetime.strptime(d, "%d-%b-%Y")
        except ValueError:
            return datetime.min

    return sorted(days, key=_parse, reverse=True)


async def get_ticks_for_day(ticker: str, date: str) -> list[dict]:
    """Return all per-second ticks for a given date string (DD-Mon-YYYY), ordered by ts."""
    async with _pool.acquire() as conn:
        rows = await conn.fetch("""
            SELECT EXTRACT(EPOCH FROM ts)::BIGINT AS ts_unix,
                   open, high, low, close,
                   ema5, ema20, ema50, hl, avg, support, rsi
            FROM ticker_ticks
            WHERE ticker = $1 AND date = $2
            ORDER BY ts ASC
        """, ticker, date)
    return [dict(r) for r in rows]


def aggregate_ticks(rows: list[dict], period_sec: int) -> list[dict]:
    """Aggregate per-second tick rows into OHLC bars of period_sec width."""
    buckets: dict[int, dict] = {}
    for r in rows:
        ts = int(r["ts_unix"])
        bucket_ts = (ts // period_sec) * period_sec
        if bucket_ts not in buckets:
            buckets[bucket_ts] = {
                "time": bucket_ts,
                "open": r["open"],
                "high": r["high"],
                "low": r["low"],
                "close": r["close"],
            }
        else:
            b = buckets[bucket_ts]
            b["high"] = max(b["high"], r["high"])
            b["low"] = min(b["low"], r["low"])
            b["close"] = r["close"]
    return [buckets[k] for k in sorted(buckets)]
