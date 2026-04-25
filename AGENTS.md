Act as a Quantitative Software Engineer. I am migrating a complex, duplicated Excel workbook used for share market analysis into a high-performance Python backend. 

**The Architecture Goal:**
* **Backend:** Python (FastAPI) handling live per-second ticker updates via an HTTP POST webhook.
* **Data Processing:** Pandas/TA-Lib for real-time vector math and aggregations (1m, 5m, etc.).
* **Streaming:** Pushing computed updates via Server-Sent Events (SSE) to a custom frontend ("trends").
* **Storage:** PostgreSQL with TimescaleDB for storing raw per-second ticks and continuous aggregates for future ML modeling.

**The Current Data State:**
The historical data and current logic live in Excel. The source of truth is `data/Final-Bullish-CE.xlsx` (sheet: `Nifty-20.12.2024`, data rows 5–5387). Column mapping:
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
* Support and Bullish are binary-search futures — no dedicated column in Nifty sheet; see `CE-20.12.2024` and `Bullish-20.12.2024` sheets

**Available Assets:**
1. `data/Final-Bullish-CE.xlsx` — source of truth. Key sheets:
   - `Nifty-20.12.2024`: all indicator values (rows 5–5387)
   - `CE-20.12.2024`: CE binary-search working (row 8 note: "When BP3 = BP4 then w3 value is CE")
   - `Bullish-20.12.2024`: Bullish binary-search working (W col and CD col formulas define the search)
2. Go code on `master` branch (`services/ticker/cards/`) — secondary reference. If Go and Excel disagree, **Excel wins**.

**Tooling:**
* Use `uv` to run all Python commands (e.g., `uv run python ...`, `uv add <pkg>` instead of `pip install`).
* Run tests from `trends-py/`: `uv run pytest -v`

**Source of Truth (priority order):**
1. `data/Final-Bullish-CE.xlsx` — indicator formulas and computed values
2. Go code (`master` branch, `services/ticker/cards/`) — futures binary-search logic
3. `docs/FIELDS.md` — full field reference, formulas, and implementation notes
4. `docs/PROGRESS.md` — build status, test summary, and session change log

**Data Ingestion — Two Tracks:**
* **Track 1 (Manual):** `PUT /api/update/{ticker}` — matches VB script payload `{date, close, open, high, low}`. Always available.
* **Track 2 (Zerodha):** If `ZERODHA_ACCESS_TOKEN` env var is present, connect to Zerodha KiteTicker WebSocket on startup. On tick: extract OHLC + LTP and call the same internal update function as Track 1. If token missing or WebSocket errors, fall back to Track 1 silently.

**Startup Seeding:**
* On startup: check TimescaleDB row count for the ticker. If ≥ 50 rows exist, seed EMA state from DB. Otherwise, seed from Excel (`data/Final-Bullish-CE.xlsx`, sheet `Nifty-20.12.2024`, rows 5–5387).

**Futures (Support / Bullish) — Verified Algorithm:**

All futures searches use `ema5_pre`/`ema20_pre` (EMA state **before** the current bar's close).
`ema5_post`/`ema20_post` are **not used**.

Computation chain per tick:
```
CE2           ← brentq(_search_ce,  ema5_pre, ema20_pre)
                  W2=W3=trial, find BP[d+2]=BP[d+1]
CD2           ← 2/6*(CE2 - old_cd) + old_cd             [persisted in TickerState.cd]
CC            ← brentq(_search_cc,  ema5_pre, ema20_pre, CD2)
                  W[d+1]=EMA_step(CD2,trial), W[d+2]=W[d+3]=trial, find BP[d+3]=BP[d+2]
                  → Support = W[d+1] = 2/6*(cc_trial - CD2) + CD2
CE3           ← brentq(_search_ce3, ema5_pre, ema20_pre, close)
                  W2=close (fixed), W3=W4=trial, find BP[d+3]=BP[d+2]
CD3           ← 2/6*(CE3 - CD2) + CD2
bullish_trial ← brentq(_search_bullish, ema5_pre, ema20_pre, CD3)
                  W2=W3=CD4(trial)=2/6*(trial-CD3)+CD3, W4=W5=trial
                  find BP[d+5]=BP[d+4]
Bullish       ← 2/6*(bullish_trial - CD3) + CD3   [= CD4 = W2 = W3, NOT the raw trial]
```

Key notes:
* `close` is **not** in the Bullish W sequence. It feeds the chain only via CE3 → CD3.
* CD4 inside `_search_bullish` **varies with trial** — recomputed on every brentq call.
* Bullish return value = CD4(trial_converged) = W2 = W3 (source: Bullish sheet CB column).
* CE3 uses pre-bar EMA. Applying `close` as W2 projects the EMA one step; then W3=W4=trial.
* Only recalculate futures when `High` changes (gate: `TickerState._last_futures_high`).
* Use `scipy.optimize.brentq` for all binary searches.

**Directives:**
1. **Standardize Calculations:** Use standard Python libraries (`scipy`) for all indicators. Do not port custom Go EMA loops — use the Excel formula spec instead.
2. **Tests:** Write tests by loading Excel computed values (`data_only=True` via openpyxl) and asserting within `TOLERANCE=0.001`. Tests verify **mathematical convergence** (|BP[n+1]-BP[n]| ≤ TOLERANCE), not hard-coded expected values.
3. **Multi-ticker:** All API endpoints include `{ticker}` in the path. State is keyed by ticker name.
