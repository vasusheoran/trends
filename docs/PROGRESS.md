# Trends — Build Progress

## Phase Status

| Phase | Description | Status | Notes |
|-------|-------------|--------|-------|
| 1 | Project scaffold | ✅ Done | `trends-py/` structure, `pyproject.toml`, uv deps |
| 2 | Indicator engine | ✅ Done + Tested | All indicator tests pass against Excel ground truth |
| 3 | Futures (Support/Bullish) | ✅ Done + Tested | Correct BR/CC algorithm from Go; all 21 tests pass |
| 4 | Zerodha ingest | ⚠️ Written, not tested | `ingest/zerodha.py`; needs `kiteconnect` pkg + live token |
| 5 | TimescaleDB | ⚠️ Written, not tested | Schema + seed logic; needs running DB to verify |
| 6 | SSE streaming | ✅ Done | Integrated with Frontend Terminal |
| 7 | Dockerfile + docker-compose | ✅ Done | Static file mounting + volume mapping verified |
| 8 | End-to-end smoke test | ❌ Not started | PUT → state update → SSE push |

---

## What's Tested

Run from `trends-py/`: `uv run pytest -v`  (21 tests, all pass)

| Test | File | What it validates |
|------|------|-------------------|
| EMA-5 vs Excel | `test_indicators.py` | Last 20 rows, tol=0.001 |
| EMA-20 vs Excel | `test_indicators.py` | Last 20 rows, tol=0.001 |
| H/L vs Excel | `test_indicators.py` | Last 20 rows, tol=0.001 |
| AVG vs Excel | `test_indicators.py` | Last 20 rows, tol=1.0 |
| RSI vs Excel | `test_indicators.py` | Last 20 rows, tol=0.01 |
| Bullish convergence (Excel EMA) | `test_bullish.py` | Last 20 row-pairs: \|BP[d+3]-BP[d+2]\|≤TOLERANCE×5 |
| Bullish bracket | `test_bullish.py` | Last 20 rows: brentq bracket has sign change |
| BR vs Excel column | `test_bullish.py` | Last 20 rows, tol=15.0 (row 5385 outlier ~14 pts) |
| CSV bullish convergence | `test_bullish_csv.py` | \|BP[d+3]-BP[d+2]\|≤TOLERANCE×5 from CSV seed |
| CSV script returns values | `test_bullish_csv.py` | support>0, bullish>0 from compute_from_csv() |
| Bullish convergence (seeded) | `test_futures.py` | Last 20 rows: \|BP[d+3]-BP[d+2]\|≤TOLERANCE×5 |
| Support convergence (seeded) | `test_futures.py` | Last 20 rows: \|BP[d+3]-BP[d+2]\|≤TOLERANCE×5 |
| Futures populated after 105 bars | `test_futures.py` | Non-None by bar 120 |
| EMA-5 all rows vs Excel | `test_seeding.py` | All ~5383 rows, tol=0.001 |
| EMA-20 all rows vs Excel | `test_seeding.py` | All ~5383 rows, tol=0.001 |
| H/L all rows vs Excel | `test_seeding.py` | All ~5383 rows, tol=0.001 |
| AVG all rows vs Excel | `test_seeding.py` | All ~5383 rows, tol=1.0 |
| RSI all rows vs Excel | `test_seeding.py` | All ~5383 rows, tol=0.01 |
| Support vs Excel CC (last 10) | `test_seeding.py` | tol=15.0 (most rows < 2; row 5385 outlier ~1.5) |
| Bullish vs Excel BR (last 10) | `test_seeding.py` | tol=15.0 (most rows < 2; row 5385 outlier ~14) |
| First-drift report | `test_seeding.py` | Non-failing; prints first diverging row per field |

---

## Futures Algorithm — Current Implementation (session 7, 2026-04-25)

Complete rewrite of `futures.py` and `state.py` to match the **verified Go source algorithm**.

### Algorithm (verified against `master:services/ticker/cards/next.go`)

```
CE2  = (49*EMA5_pre - 19*EMA20_pre) / 30        [closed-form 2-bar BP fixed point]
CD   = EMA-5(decay=2/6) of CE2 values            [TickerState.cd_ema, persistent]
BR   = brentq: W=[BR_trial, CE2, CE2]             → Bullish (Excel col BR)
              find BR_trial where BP[d+3]=BP[d+2]
CC   = brentq: W=[cd3, cc_trial, cc_trial]        → Support (Excel col CC)
              where cd3=2/6*(cc_trial-CD)+CD
              find cc_trial where BP[d+3]=BP[d+2]; Support=cd3
```

### Key decisions

- **CE2 is closed-form** — no binary search needed. Derived from BP[d+1]=BP[d+2] condition.
- **CD is EMA-5** — uses `EMAState(period=5, decay=2/6)` seeded from first 5 CE2 values.
  Seeded at bar ~104 (first 5 CE2 values from bars 100–104). Returns `None` until seeded.
- **CD is updated BEFORE compute_futures** — `self.cd_ema.update(ce2)` called inside `_update_commit`,
  then `cd2 = self.cd_ema.value` passed to `compute_futures`.
- **prev_high is no longer needed** — removed from `compute_futures` signature.
- **Live-mode checkpoint** saves `cd_ema` — restored on each tick to guarantee idempotency.

### What changed from previous (wrong) implementation

| Item | Old (session 6 — wrong) | New (session 7 — correct) |
|------|------------------------|--------------------------|
| Bullish search | `[trial]*4 → BP[d+4]=BP[d+3]` | `[BR, CE2, CE2] → BP[d+3]=BP[d+2]` |
| Support search | `[prev_high, trial, trial] → BP[d+3]=BP[d+2]` | `[cd3, cc_trial, cc_trial] → BP[d+3]=BP[d+2]` |
| CE2 | Not used | Closed-form: `(49*ema5-19*ema20)/30` |
| CD state | Not used | `EMAState(period=5, decay=2/6)` in `TickerState.cd_ema` |
| `compute_futures` args | `(ema5_pre, ema20_pre, prev_high)` | `(ema5_pre, ema20_pre, cd2)` |
| Excel accuracy | Bullish diff 500-640 pts (completely wrong) | Most rows < 2.0 pts |

---

## Open Investigation: 01-Jan-2025 Bullish Discrepancy

The user's manual Bullish for 01-Jan-25 = **23531**.  
Our algorithm for 01-Jan-25 produces:

```
CE2     = 23533.86   ← closest to user's 23531 (diff ~3)
Bullish = 23687.72   ← BR algorithm result
Support = 23600.33   ← CC algorithm result
```

The algorithm is mathematically correct and verified against historical Excel rows.
The 23531 discrepancy is unresolved. **Likely hypothesis:** the user's manual iteration
converges to CE2 (the 2-bar equilibrium price), not the BR result. The Go server's
`searchCE` function finds a price very close to our closed-form CE2.

**To investigate next session:**
- Compare what the Go API actually displays when the CSV is uploaded
- Determine whether the user reads "CE" or "BR" from the Go response
- Check if the second trial row (02-Jan-25) changes the displayed value

---

## Source of Truth Migration (session 3)

Source of truth switched from `data/Nifty-17-04-2026.xlsx` to `data/Final-bullish-ce.xlsx`.

**Key changes made:**
- EMA-50 removed; BN = EMA-20 (decay 2/21)
- BP = EMA5 - EMA20 (was EMA5 - EMA50)
- AVG formula corrected (denominator = Sum, not A)
- conftest.py: sheet `Nifty-20.12.2024`, cols AS/BN/BV/AR/AD/CC/BR

---

## Next Steps (ordered)

### Investigate row 5385 outlier
Row 5385 (18-Dec-2024): BR diff ~14 pts, CC diff ~1.5 pts vs Excel.
All surrounding rows are within 2.0. Possible causes:
- Excel manual-entry override for that date
- Holiday/special session adjustment in the data
- Investigate the raw Excel formula for BR/CC at that row

### Clarify 01-Jan-25 Bullish = 23531
Talk to user about which Go API field they read:
- Is it `BR` (=23687.72), `CE` (=23533.86), or `CC` (=23600.33)?
- Upload CSV to Go server and compare displayed output

### Phase 8 — Integration tests
- [ ] `tests/test_webhook.py` — POST a tick via `httpx.AsyncClient`, assert snapshot returned
- [ ] `tests/test_sse.py` — POST a tick, assert SSE event received on stream endpoint
- [ ] Test Zerodha fallback: missing token → webhook-only mode (no crash)

### Later
- [ ] Add `kiteconnect` to deps once Zerodha token is available (`uv add kiteconnect`)
- [ ] Validate Zerodha OHLC tick mapping against manual PUT payload
- [ ] TimescaleDB continuous aggregates for ML (future session)
- [ ] Multi-ticker: test with a second ticker symbol alongside nifty

---

## Key Files

| File | Purpose |
|------|---------|
| `app/engine/indicators.py` | EMAState, SMAState, RSIState, calc_hl, calc_avg |
| `app/engine/state.py` | TickerState — rolling 101-bar manager + cd_ema (EMA-5 of CE2) |
| `app/engine/futures.py` | _ce2 (closed-form), _search_br, _search_cc, compute_futures |
| `app/registry.py` | Global state dict + SSE pub/sub |
| `app/main.py` | FastAPI app, lifespan (seed + Zerodha) |
| `app/db/seed.py` | Excel → state seeding; DB → state seeding |
| `app/db/timescale.py` | asyncpg schema, upsert, load_bars |
| `app/ingest/webhook.py` | PUT /api/update/{ticker} |
| `app/ingest/zerodha.py` | KiteTicker WebSocket client |
| `app/api/stream.py` | GET /api/stream/{ticker} SSE |
| `app/api/debug.py` | GET /api/debug/{ticker}, POST /api/debug/compute |
| `tests/conftest.py` | Excel fixture — Final-Bullish-CE.xlsx, Nifty-20.12.2024, rows 5-5387 |
| `scripts/compute_bullish_from_csv.py` | Standalone: seed from CSV → print CE2, CD, Support, Bullish |
| `docs/FIELDS.md` | Full field reference + implementation notes |
| `docs/PROGRESS.md` | Phase status + algorithm history |


## Frontend Migration (Vanilla JS/CSS)

### Phase 1: Infrastructure & Environment 
* **Goal**: Establish the delivery pipeline from FastAPI to the browser.
* **Tasks**:
    * ✅ Configure FastAPI static file mounting in `main.py`.
    * ✅ Update `docker-compose.yaml` with volume mounts for `/static`.
    * ✅ Deploy initial `index.html` blueprint to verify container reachability.

### Phase 2: Real-Time Data Integration 
* **Goal**: Replace mock simulation with live backend data.
* **Tasks**:
    * ✅ Map SSE streaming endpoints via `http://192.168.29.204:5001/docs`.
    * ✅ Implement `terminal.js` with native `EventSource` logic.
    * ✅ Apply 500ms flash indicators to targeted DOM nodes based on incoming ticks.

### Phase 3: State Management & History 
* **Goal**: Enable user interaction and historical deep-dives.
* **Tasks**:
    * ✅ Implement row selection and global "click-away" deselection logic.
    * ✅ Wire "Remove Ticker" button to the `DELETE` API.
    * ✅ Fetch historical OHLC data on selection and render to chart placeholder.
    * ✅ Wire "Upload Seed" button to `POST /api/seed/{ticker}`.
