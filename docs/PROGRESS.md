# Trends — Build Progress

## Phase Status

| Phase | Description | Status | Notes |
|-------|-------------|--------|-------|
| 1 | Project scaffold | ✅ Done | `trends-py/` structure, `pyproject.toml`, uv deps |
| 2 | Indicator engine | ✅ Done + Tested | All indicator tests pass against Excel ground truth |
| 3 | Futures (Support/Bullish) | ✅ Done + Tested | CE/CC/BR convergence verified; 9 tests pass |
| 4 | Zerodha ingest | ⚠️ Written, not tested | `ingest/zerodha.py`; needs `kiteconnect` pkg + live token |
| 5 | TimescaleDB | ⚠️ Written, not tested | Schema + seed logic; needs running DB to verify |
| 6 | SSE streaming | ⚠️ Written, not tested | `api/stream.py`; needs integration test |
| 7 | Dockerfile + docker-compose | ❌ Not started | Needs TimescaleDB service + app container |
| 8 | End-to-end smoke test | ❌ Not started | PUT → state update → SSE push |

---

## What's Tested

Run from repo root: `uv run --project trends-py python -m pytest trends-py/tests/`

Currently **9 tests**, all passing.

| Test | File | What it validates |
|------|------|-------------------|
| EMA-5 vs Excel | `test_indicators.py` | Last 20 rows, tol=0.001 |
| EMA-20 vs Excel | `test_indicators.py` | Last 20 rows, tol=0.001 |
| H/L vs Excel | `test_indicators.py` | Last 20 rows, tol=0.001 |
| AVG vs Excel | `test_indicators.py` | Last 20 rows, tol=1.0 |
| RSI vs Excel | `test_indicators.py` | Last 20 rows, tol=0.01 |
| CE convergence | `test_futures.py` | Last 20 rows: \|BP[d+2]-BP[d+1]\| ≤ 0.005 at CE |
| CC convergence | `test_futures.py` | Last 20 rows: \|BP[d+3]-BP[d+2]\| ≤ 0.005 at CC (Support) |
| BR convergence | `test_futures.py` | Last 20 rows: \|BP[d+3]-BP[d+2]\| ≤ 0.005 at BR (Bullish) |
| Futures warm-up | `test_futures.py` | None before bar 100, non-None by bar 120 |

---

## Source of Truth Migration (session 3)

Source of truth switched from `data/Nifty-17-04-2026.xlsx` to `data/Final-bullish-ce.xlsx`.

**Key changes made:**

### EMA-50 removed
The new Excel file does not have an EMA-50 column. BN in the new file is EMA-20 (decay 2/21).
- Removed `EMAState(period=50, decay=2/51)` from `state.py`
- Removed `ema50` from `TickerSnapshot`
- Removed `test_ema50_matches_excel` from tests
- BP formula changed from `EMA5 - EMA50` to `EMA5 - EMA20`

### AVG formula corrected
The AVG correction denominator must be `Sum` (= SMA10+SMA50), **not** `A` (= Sum/2).
Multipliers: `0.01` and `0.025` (not `0.0001`/`0.0001` as in old Excel).
- Fixed `calc_avg` in `indicators.py`

### conftest.py updated
Now loads `Final-bullish-ce.xlsx`, sheet `Nifty-20.12.2024`, rows 5–5387.
Column mapping: `AS`→ema5, `BN`→ema20, `BV`→rsi, `AR`→avg, `AD`→hl.

---

## Futures Overhaul (session 3)

Support is **CC** (not CE as originally implemented). Bullish remains **BR**.

### What changed
- `compute_futures` now takes `old_cd` (previous rolling CD value) instead of `cd`
- CD updated inside `compute_futures`: `new_cd = 2/6*(ce - old_cd) + old_cd`
- CC (Support) binary search uses `new_cd` for the W[d+1] step
- Returns `(support, bullish, new_cd)` — state.py stores `new_cd` back into `self.cd`
- `TickerState` gains `cd: float = 0.0` field (persistent rolling EMA5 of CE)

### Computation chain (matches Go `updateFutureData`)
```
CE  ← brentq(_search_ce,  ema5_pre, ema20_pre)
CD  ← 2/6*(CE - old_cd) + old_cd            [updates TickerState.cd]
CC  ← brentq(_search_cc,  ema5_pre, ema20_pre, CD)  → Support = W[d+1]
BR  ← brentq(_search_br,  ema5_pre, ema20_pre, CE)  → Bullish
```

### Previous bugs fixed across sessions
1. **Wrong EMA starting state**: `compute_futures` must receive EMA copies from **before** the current bar's close is applied. (`ema5_pre`/`ema20_pre` captured before `self.ema5.update(close)`)
2. **Wrong BR convergence**: `_search_br` projects `[trial, CE, CE]`, not `[trial, trial, trial]`
3. **Support was CE, not CC**: Return value corrected; CD rolling state added

---

## Next Steps (ordered)

### Phase 7 — Docker
- [ ] Write `Dockerfile` for the FastAPI app
- [ ] Update `docker-compose.yaml` with `timescaledb/timescaledb` service + app service
- [ ] Verify seed-from-Excel runs cleanly on first boot (no DB rows)
- [ ] Verify seed-from-DB runs on subsequent boots

### Phase 8 — Integration tests
- [ ] `tests/test_webhook.py` — POST a tick via `httpx.AsyncClient`, assert snapshot returned
- [ ] `tests/test_sse.py` — POST a tick, assert SSE event received on stream endpoint
- [ ] Test Zerodha fallback: missing token → webhook-only mode (no crash)

### Later
- [ ] Add `kiteconnect` to deps once Zerodha token is available (`uv add kiteconnect`)
- [ ] Validate Zerodha OHLC tick mapping against manual PUT payload
- [ ] TimescaleDB continuous aggregates for ML (future session)
- [ ] Multi-ticker: test with a second ticker symbol alongside nifty
- [ ] Investigate CH/CI/CJ futures (next bullish levels from Go `next.go`) if needed

---

## Key Files

| File | Purpose |
|------|---------|
| `app/engine/indicators.py` | EMAState, SMAState, RSIState, calc_hl, calc_avg |
| `app/engine/state.py` | TickerState — rolling 101-bar manager + CD rolling state |
| `app/engine/futures.py` | CE/CC/BR binary search via brentq; returns (support, bullish, new_cd) |
| `app/registry.py` | Global state dict + SSE pub/sub |
| `app/main.py` | FastAPI app, lifespan (seed + Zerodha) |
| `app/db/seed.py` | Excel → state seeding; DB → state seeding |
| `app/db/timescale.py` | asyncpg schema, upsert, load_bars |
| `app/ingest/webhook.py` | PUT /api/update/{ticker} |
| `app/ingest/zerodha.py` | KiteTicker WebSocket client |
| `app/api/stream.py` | GET /api/stream/{ticker} SSE |
| `tests/conftest.py` | Excel fixture — Final-bullish-ce.xlsx, Nifty-20.12.2024 sheet |
| `docs/FIELDS.md` | Full field reference + implementation notes |
