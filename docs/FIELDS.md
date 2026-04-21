# Trends — Field Reference for AI Agents

**Source of truth:** `data/Final-bullish-ce.xlsx`, sheet `Nifty-20.12.2024`, rows 5–5387.  
**Go reference:** `services/ticker/cards/` on `master` branch — use for futures binary-search logic. If Go and Excel disagree, **Excel wins**.  
**Tests:** Write tests by seeding from Excel computed values (`data_only=True` via openpyxl).

---

## Required Output Fields

These are the fields exposed via SSE and stored in TimescaleDB.

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
| Support    | —         | CC binary search result — see Futures section | 100+ |
| Bullish    | —         | BR binary search result — see Futures section | 100+ |

> There is no EMA-50 in the current implementation. The old `data/Nifty-17-04-2026.xlsx`
> had a BN column labeled EMA-50 (decay 2/51), but the new source of truth uses BN = EMA-20 (decay 2/21).

---

## Intermediate Fields (Internal State, Not Exposed)

### EMA State (2 distinct EMAs)

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

Key: the denominator in the correction term is `Sum` (= SMA10+SMA50), **not** `A`. Excel operator
precedence `/ Sum / 2 * 100 / 2` ≠ `/ A * 100 / 2`. Tolerance 1.0 used in tests due to correction formula sensitivity.

### RSI Formula (Wilder Smoothing)

```
C[t] = Close[t] - Close[t-1]
gain[t] = max(C[t], 0)
loss[t] = abs(min(C[t], 0))

# Seed at bar 14 (14 bars):
avg_gain[14] = mean(gain[1:14])
avg_loss[14] = mean(loss[1:14])

# Rolling Wilder smooth (bar > 14):
avg_gain[t] = (avg_gain[t-1] * 13 + gain[t]) / 14
avg_loss[t] = (avg_loss[t-1] * 13 + loss[t]) / 14

RSI[t] = 100 if avg_loss[t] == 0 else 100 - (100 / (1 + avg_gain[t] / avg_loss[t]))
```

Go field: `CW`. Excel col: `BV`.

---

## Future Projection Fields (Binary Search)

**Computation chain per tick (matches Go `updateFutureData`):**
1. CE — binary search (stateless)
2. CD — rolling EMA5 of CE values (persistent `TickerState.cd`)
3. CC → Support — binary search using new CD
4. BR → Bullish — binary search using CE

**Support = CC (not CE). Bullish = BR.**

### CE

Binary search for `trial` where `BP[d+1] = BP[d+2]`, with `W[d+1] = W[d+2] = trial`.

```python
def _search_ce(trial, ema5, ema20):
    bp = _bp_series(ema5, ema20, [trial, trial])
    return bp[1] - bp[0]   # zero when converged

ce = brentq(_search_ce, 0.0, 99999.0, args=(ema5_pre, ema20_pre), xtol=0.001)
```

CE is an intermediate value — not exposed directly, but used to update CD and compute BR.

### CD (persistent state)

Rolling EMA5 of CE values. Updated every tick after CE is computed:

```python
new_cd = 2/6 * (ce - old_cd) + old_cd
```

Stored in `TickerState.cd`. Initialized to 0.0 (seeds naturally as CE accumulates).

### CC (Support)

Binary search for `cc_trial` where `BP[d+2] = BP[d+3]`, with:
- `W[d+1] = EMA_step(CD, cc_trial)` — one EMA step from CD toward cc_trial
- `W[d+2] = W[d+3] = cc_trial`

CC (Support) = `W[d+1]` = `2/6 * (cc_trial - CD) + CD` — **not** the raw cc_trial.

```python
def _search_cc(trial, ema5, ema20, cd):
    w_d1 = 2/6 * (trial - cd) + cd
    bp = _bp_series(ema5, ema20, [w_d1, trial, trial])
    return bp[2] - bp[1]   # zero when converged

cc_trial = brentq(_search_cc, 0.0, 99999.0, args=(ema5_pre, ema20_pre, new_cd), xtol=0.001)
support = 2/6 * (cc_trial - new_cd) + new_cd   # = W[d+1]
```

### BR (Bullish)

Binary search for `trial` where `BP[d+3] = BP[d+2]`, with:
- `W[d+1] = trial`
- `W[d+2] = W[d+3] = CE`

```python
def _search_br(trial, ema5, ema20, ce):
    bp = _bp_series(ema5, ema20, [trial, ce, ce])
    return bp[2] - bp[1]   # zero when converged

bullish = brentq(_search_br, 0.0, 99999.0, args=(ema5_pre, ema20_pre, ce), xtol=0.001)
```

### BP Series Helper

```python
def _bp_series(ema5, ema20, closes):
    """Apply closes to COPIES of EMA state, return BP = EMA5 - EMA20 at each step."""
    m, o = ema5.copy(), ema20.copy()
    return [mv - ov for mv, ov in ((m.update(c), o.update(c)) for c in closes)]
```

Originals are never mutated — copies used for projection.

### Futures Gate

Futures are only recomputed when `High` has changed (matches Go `recalculateCH`). Gate field: `TickerState._last_futures_high`. Reset to `0.0` in tests to force recomputation every bar.

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
  └─► Compute futures (when High is new, needs 100+ bars):
        CE  ← brentq(_search_ce,  ema5_pre, ema20_pre)
        CD  ← 2/6*(CE - old_cd) + old_cd            [persistent state]
        CC  ← brentq(_search_cc,  ema5_pre, ema20_pre, CD)  → Support
        BR  ← brentq(_search_br,  ema5_pre, ema20_pre, CE)  → Bullish
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

Futures require at least 100 bars of history (matches Go `if current.Index < 100: return`).

---

## ⚠️ Why Excel BV/BW Cannot Be Used as Futures Ground Truth

The `BV` (Support) and `BW` (Bullish) columns in the old `data/Nifty-17-04-2026.xlsx` contain
**hard-coded values stored by the Go server** — they are not live Excel formulas. Verified:
opening without `data_only=True` shows numeric values, not formula strings.

The new source of truth `data/Final-bullish-ce.xlsx` does not have pre-computed Support/Bullish
columns in the main sheet. Tests instead verify **mathematical convergence properties**:
- CE: `|BP[d+2] - BP[d+1]| ≤ TOLERANCE` when both days use `trial=CE`
- CC: `|BP[d+3] - BP[d+2]| ≤ TOLERANCE` with `W[d+1]=support, W[d+2,3]=cc_trial`
- BR: `|BP[d+3] - BP[d+2]| ≤ TOLERANCE` with `W[d+1]=BR, W[d+2,3]=CE`

---

## Adding a New Field (for future sessions)

1. Check if formula exists in Excel — if so, Excel is the implementation spec.
2. Check Go `services/ticker/cards/helper.go` for the function.
3. If it's a binary-search future: follow the pattern in `futures.py`, add a new `_search_*` function.
4. Add field to `TickerSnapshot` model and SSE output.
5. Write a test seeding from Excel computed values and assert within `TOLERANCE=0.001`.
