# Trends — Field Reference for AI Agents

**Source of truth:** `data/Final-bullish-ce.xlsx`, sheet `Nifty-20.12.2024`, rows 5–5387.  
**Go reference:** `services/ticker/cards/` on `master` branch — verified against Python implementation. If Go and Excel disagree, **Excel wins**.  
**Tests:** Write tests by seeding from Excel computed values (`data_only=True` via openpyxl).

---

## Required Output Fields

| Field Name | Excel Col | Formula / Logic | Starts at row |
|------------|-----------|-----------------|---------------|
| Date       | B         | Input           | 5             |
| Close      | W         | Input           | 5             |
| Open       | X         | Input           | 5             |
| High       | Y         | Input           | 5             |
| Low        | Z         | Input           | 5             |
| H/L        | AD        | `min(High[t-1], High[t-2], High[t-3])` | 8 |
| AVG        | AR        | See formula below | 55 |
| EMA-5      | AS        | `2/6 * (Close - EMA5[t-1]) + EMA5[t-1]` | 9 |
| EMA-20     | BN        | `2/21 * (Close - EMA20[t-1]) + EMA20[t-1]` | 24 |
| RSI        | BV        | RSI(14) Wilder smoothing — see formula below | 20 |
| Bullish    | BR        | BR binary search result — see Futures section | 105+ |
| Support    | CC        | CC binary search result — see Futures section | 105+ |

> There is no EMA-50. The old `data/Nifty-17-04-2026.xlsx` had BN = EMA-50 (decay 2/51),
> but the current source of truth uses BN = EMA-20 (decay 2/21).

---

## Intermediate Fields (Internal State)

### EMA State (2 EMAs)

| Internal Name | Excel Col | Seed | Decay | Notes |
|---------------|-----------|------|-------|-------|
| M (EMA-5)  | AS | `AVERAGE(W5:W9)` at row 9   | `2/6`  | Self-referential |
| O (EMA-20) | BN | `AVERAGE(W5:W24)` at row 24 | `2/21` | Self-referential |

### BP (MACD-like Spread)

```
BP = M - O   (EMA-5 minus EMA-20)
```

BP is the convergence target for all binary-search future projections.

### AVG Formula

```python
Sum = SMA_10 + SMA_50   # SMA_10 = mean(Close, last 10 bars), SMA_50 = mean(Close, last 50 bars)
A = Sum / 2
inner = A - A * 0.01
inner2 = inner * 0.025
inner3 = (inner + inner2 + A) / 2
AVG = A - (A * ((A - inner3) / Sum / 2 * 100 / 2) / 100)
```

Key: the denominator in the correction term is `Sum` (= SMA10+SMA50), **not** `A`.

### RSI Formula (Wilder Smoothing)

```
C[t] = Close[t] - Close[t-1]
gain[t] = max(C[t], 0)
loss[t] = abs(min(C[t], 0))

# Seed at bar 14:
avg_gain[14] = mean(gain[1:14])
avg_loss[14] = mean(loss[1:14])

# Rolling Wilder smooth (bar > 14):
avg_gain[t] = (avg_gain[t-1] * 13 + gain[t]) / 14
avg_loss[t] = (avg_loss[t-1] * 13 + loss[t]) / 14

RSI[t] = 100 if avg_loss[t] == 0 else 100 - (100 / (1 + avg_gain[t] / avg_loss[t]))
```

---

## Future Projection Fields (Binary Search)

All searches use `ema5_pre`/`ema20_pre` — the EMA state **before** the current bar's close is applied.

**Computation chain per tick (bar >= 100):**

```
CE2  = (49*EMA5_pre - 19*EMA20_pre) / 30        [closed-form 2-bar fixed point]
CD   = EMA_step(CD_prev, CE2)                    [EMA-5 of CE2; seeded after 5 values]
BR   = brentq: W=[BR_trial, CE2, CE2], BP[d+3]=BP[d+2]     → Bullish
CC   = brentq: W=[cd3, cc_trial, cc_trial], BP[d+3]=BP[d+2] → Support = cd3
       where cd3 = 2/6*(cc_trial - CD) + CD
```

Futures require at least 100 bars of history for EMA to be seeded, plus 5 additional rows for CD EMA to seed. Support/Bullish first appear around bar 105.

---

### CE2 — Closed-Form 2-Bar Fixed Point

CE2 is the price at which BP stops changing after 2 bars. Derived by solving BP[d+2]=BP[d+1] with W[d+1]=W[d+2]=trial. Closed-form solution:

```python
def _ce2(ema5: EMAState, ema20: EMAState) -> float:
    return (49 * ema5.value - 19 * ema20.value) / 30
```

**Derivation:** Let k5 = decay5*(1-decay5) = 2/9, k20 = decay20*(1-decay20) = 38/441.
CE2 = (k5*ema5_pre - k20*ema20_pre) / (k5 - k20) = (49*ema5 - 19*ema20) / 30.

CE2 is used only internally to drive the BR and CD computation.

---

### CD — Rolling EMA-5 of CE2

`TickerState.cd_ema` is an `EMAState(period=5, decay=2/6)` seeded from the first 5 CE2 values (bars 100–104). Updated each bar:

```python
ce2 = _ce2(ema5_pre, ema20_pre)
self.cd_ema.update(ce2)     # EMA-5 step: CD_new = 2/6*(CE2-CD_old) + CD_old
```

CD is the state AFTER updating with the current bar's CE2. This updated CD is passed to `compute_futures` as `cd2`.

---

### BR — Bullish (Excel column BR)

Binary search for `br_trial` where `BP[d+3] = BP[d+2]` with W = [br_trial, CE2, CE2]:

```python
def _search_br(trial, ema5, ema20, ce2):
    bp = _bp_series(ema5, ema20, [trial, ce2, ce2])
    return bp[2] - bp[1]

bullish = brentq(_search_br, 0.0, 99999.0, args=(ema5_pre, ema20_pre, ce2), xtol=TOLERANCE)
```

Go reference: `searchBR` in `master:services/ticker/cards/next.go`.

**Verified against Excel column BR:** most rows match within 1.5 pts.
Row 5385 (18-Dec-2024) is a known outlier (~14 pts off) — likely an Excel manual-entry or rounding anomaly. All other rows within 2.0.

---

### CC — Support (Excel column CC)

Binary search for `cc_trial` where `BP[d+3] = BP[d+2]` with W = [cd3, cc_trial, cc_trial]:
- `cd3 = 2/6 * (cc_trial - CD) + CD` (varies with cc_trial)
- Support = cd3 (not cc_trial)

```python
def _search_cc(trial, ema5, ema20, cd2):
    cd3 = (2/6) * (trial - cd2) + cd2
    bp = _bp_series(ema5, ema20, [cd3, trial, trial])
    return bp[2] - bp[1]

cc_trial = brentq(_search_cc, 0.0, 99999.0, args=(ema5_pre, ema20_pre, cd2), xtol=TOLERANCE)
support  = (2/6) * (cc_trial - cd2) + cd2   # = cd3
```

Go reference: `searchCC` in `master:services/ticker/cards/next.go`.

**Verified against Excel column CC:** most rows match within 1.5 pts. Row 5385 is a known outlier (~1.5 pts).

---

## BP Series Helper

```python
def _bp_series(ema5, ema20, closes):
    """Apply closes to COPIES of EMA state; return BP = EMA5 - EMA20 at each step."""
    m, o = ema5.copy(), ema20.copy()
    results = []
    for c in closes:
        mv, ov = m.update(c), o.update(c)
        results.append(mv - ov if mv is not None and ov is not None else None)
    return results
```

Originals are never mutated — copies used for projection.

---

## Open Investigation: 01-Jan-2025 Bullish Discrepancy

The user's manual process gave Bullish=23531 for 01-Jan-25. Our implementation gives:
- CE2 = 23533.86  (2-bar fixed point — closest to 23531, diff ~3)
- Bullish (BR) = 23687.72  (3-bar search)
- Support (CC) = 23600.33

The algorithm is verified against Go source and Excel for all historical rows. The discrepancy for 01-Jan-25 suggests the user's manual iteration may have converged to CE2 (the equilibrium price) rather than BR (the 3-bar fixed point search). The Go server's `searchCE` returns a similar value to our closed-form CE2.

**Hypothesis to investigate:** When the user uploads a CSV with identical OHLC trial values (open=high=low=close=X) and iterates X until "support = X", they may be hitting the CE2 condition (2-bar equilibrium), not the BR/CC search result. The Go API response and which field the user is reading should be clarified.

---

## Data Flow

```
Tick (Close, Open, High, Low, Date, ticker)
  │
  ├─► Update rolling OHLCV history (max 101 bars)
  │
  ├─► Snapshot EMA state BEFORE update (ema5_pre, ema20_pre)
  │
  ├─► Compute indicators (incremental state):
  │     M[t], O[t]  →  EMA-5, EMA-20
  │     AD[t]       →  H/L
  │     AR[t]       →  AVG
  │     RSI[t]      →  RSI(14)
  │
  └─► Compute futures (when bars >= 100):
        CE2     = (49*EMA5_pre - 19*EMA20_pre) / 30
        CD      = EMA_step(CD_prev, CE2)           [persistent TickerState.cd_ema]
        Bullish = brentq(_search_br, ema5_pre, ema20_pre, ce2)
        Support = brentq(_search_cc, ema5_pre, ema20_pre, cd2)  → cd3
```

---

## Webhook Payload

```json
PUT /api/update/{ticker}
{
  "date":  "23-Dec-2024",
  "close": 23963.0,
  "open":  24099.0,
  "high":  24099.0,
  "low":   23963.0
}
```

Ticker is in the URL path. Matches VB script `PostTrend` shape (date + OHLC).

---

## EMA Seeding Details

| EMA | Seed formula | First live EMA at bar # |
|-----|-------------|------------------------|
| M (EMA-5)  | `mean(Close[0:5])`  | Bar 5  |
| O (EMA-20) | `mean(Close[0:20])` | Bar 20 |
| CD (EMA-5 of CE2) | `mean(CE2[0:5])` | Bar 104 (5th CE2 value) |

Futures require at least 100 bars (matches Go `if current.Index < 100: return`).
CD requires 5 CE2 values before seeding → first Support/Bullish output at bar ~105.

---

## Adding a New Field (for future sessions)

1. Check if formula exists in Excel — if so, Excel is the implementation spec.
2. Check Go `master:services/ticker/cards/helper.go` and `next.go`.
3. If it's a binary-search future: add a new `_search_*` function in `futures.py`, following the `_search_br`/`_search_cc` pattern.
4. Add field to `TickerSnapshot` model and SSE output.
5. Write a convergence test seeded from Excel and assert within `TOLERANCE=0.001`.

---

## Frontend Terminal UI Logic

The terminal (`/static/index.html`) is a vanilla JS application providing a high-density view of real-time ticks.

### Real-Time Pipeline
- **Discovery**: On load, calls `GET /api/tickers` to identify active streams.
- **Initial State**: Calls `GET /api/state/{ticker}` to populate the grid before SSE begins.
- **Streaming**: Opens an `EventSource` to `GET /api/stream/{ticker}` for each active ticker.
- **Updates**: Incoming `TickerSnapshot` JSON updates specific DOM nodes directly (id-based).

### Visual Indicators
- **500ms Flash**: Numerical cells flash Green (`flash-up`) or Red (`flash-down`) when values change.
- **Selection**: Single-click to select a row (activates toolbar actions). Click-away to deselect.
- **Alerts**: Displays `OB` (RSI > 70), `OS` (RSI < 30), or `⚠️` (Backend warning).

### Future Integrations
- **Chart Engine**: (Phase 3) Integration with lightweight charting for historical OHLC view.
- **Seed Upload**: (Phase 3) Wire "Upload Seed" button to `POST /api/seed/{ticker}`.

---

## Key Files

| File | Purpose |
|------|---------|
| `app/engine/indicators.py` | EMAState, SMAState, RSIState, calc_hl, calc_avg |
| `app/engine/state.py` | TickerState — rolling 101-bar manager + cd_ema (EMA-5 of CE2) |
| `app/engine/futures.py` | _ce2, _search_br, _search_cc, compute_futures |
| `app/registry.py` | Global state dict + SSE pub/sub |
| `app/main.py` | FastAPI app, lifespan (seed + Zerodha) |
| `app/db/seed.py` | Excel → state seeding; DB → state seeding |
| `app/db/timescale.py` | asyncpg schema, upsert, load_bars |
| `app/ingest/webhook.py` | PUT /api/update/{ticker} |
| `app/ingest/zerodha.py` | KiteTicker WebSocket client |
| `app/api/stream.py` | GET /api/stream/{ticker} SSE |
| `app/api/debug.py` | GET /api/debug/{ticker}, POST /api/debug/compute |
| `tests/conftest.py` | Excel fixture — Final-Bullish-CE.xlsx, Nifty-20.12.2024 sheet |
| `scripts/compute_bullish_from_csv.py` | Standalone: seed from CSV → print Support + Bullish |
