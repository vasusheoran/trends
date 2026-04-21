Act as a Quantitative Software Engineer. I am migrating a complex, duplicated Excel workbook used for share market analysis into a high-performance Python backend. 

**The Architecture Goal:**
* **Backend:** Python (FastAPI) handling live per-second ticker updates via an HTTP POST webhook.
* **Data Processing:** Pandas/TA-Lib for real-time vector math and aggregations (1m, 5m, etc.).
* **Streaming:** Pushing computed updates via Server-Sent Events (SSE) to a custom frontend ("trends").
* **Storage:** PostgreSQL with TimescaleDB for storing raw per-second ticks and continuous aggregates for future ML modeling.

**The Current Data State:**
The historical data and current logic live in Excel. The source of truth is `data/Final-bullish-ce.xlsx` (sheet: `Nifty-20.12.2024`, data rows 5–5387). Column mapping:
* Column B: Date (Format: DD-Mon-YYYY)
* Column W: Close
* Column X: Open
* Column Y: High
* Column Z: Low
* Column AD: H/L
* Column AR: AVG
* Column AS: EMA-5  (decay 2/6)
* Column BN: EMA-20 (decay 2/21)
* Column BV: RSI(14)
* Support and Bullish are binary-search futures — no dedicated column in Nifty sheet; see CE-20.12.2024 and Bullish-20.12.2024 sheets

**Available Assets:**
1. `data/Final-bullish-ce.xlsx` — source of truth. Three sheets:
   - `Nifty-20.12.2024`: all indicator values (rows 5–5387)
   - `CE-20.12.2024`: CE (Support) binary-search working (rows 1–10, row 8 = CE definition)
   - `Bullish-20.12.2024`: BR (Bullish) binary-search working (row 8 = BR definition)
2. Go code on `master` branch (`services/ticker/cards/`) — secondary reference for binary-search futures logic. If Go and Excel disagree, **Excel wins**.

**Tooling:**
* Use `uv` to run all Python commands (e.g., `uv run python ...`, `uv add <pkg>` instead of `pip install`).
* Run tests from repo root: `uv run --project trends-py python -m pytest trends-py/tests/`

**Source of Truth (priority order):**
1. `data/Final-bullish-ce.xlsx` — indicator formulas and computed values
2. Go code (`master` branch, `services/ticker/cards/`) — futures binary-search logic
3. `docs/FIELDS.md` — full field reference, formulas, and implementation notes
4. `docs/PROGRESS.md` — build status, test summary, and session change log

**Data Ingestion — Two Tracks:**
* **Track 1 (Manual):** `PUT /api/update/{ticker}` — matches VB script payload `{date, close, open, high, low}`. Always available.
* **Track 2 (Zerodha):** If `ZERODHA_ACCESS_TOKEN` env var is present, connect to Zerodha KiteTicker WebSocket on startup. On tick: extract OHLC + LTP and call the same internal update function as Track 1. If token missing or WebSocket errors, fall back to Track 1 silently.

**Startup Seeding:**
* On startup: check TimescaleDB row count for the ticker. If ≥ 50 rows exist, seed EMA state from DB. Otherwise, seed from Excel (`data/Final-bullish-ce.xlsx`, sheet `Nifty-20.12.2024`, rows 5–5387).

**Futures (Support / Bullish):**
* Support = CC (NOT CE). Bullish = BR.
* CD = rolling EMA5 of CE values — persistent state in `TickerState.cd`, updated each tick.
* Recomputed fresh on every tick via `compute_futures(ema5_pre, ema20_pre, old_cd)`.
* Use `scipy.optimize.brentq` (not bisect) for binary search — superlinear convergence.
* Only recalculate when `High` changes (gate from Go `recalculateCH` logic) to avoid redundant work on flat ticks.
* Computation chain (matches Go `updateFutureData`): CE → CD update → CC (Support) → BR (Bullish)

**Directives:**
1. **Standardize Calculations:** Use standard Python libraries (Pandas, `scipy`) for all indicators. Do not port custom Go EMA loops — use the Excel formula spec instead.
2. **Tests:** Write tests by loading Excel computed values (`data_only=True` via openpyxl) and asserting within `TOLERANCE=0.001`.
3. **Multi-ticker:** All API endpoints include `{ticker}` in the path. State is keyed by ticker name.
