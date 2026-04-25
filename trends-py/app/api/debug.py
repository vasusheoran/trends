"""
GET  /api/debug/{ticker}  — futures trace for the current live state.
POST /api/debug/compute   — manual inputs, returns CSV of all intermediate values.
"""

import io
import csv
from typing import Optional

from fastapi import APIRouter, HTTPException
from fastapi.responses import PlainTextResponse
from pydantic import BaseModel

from app.registry import get_state, _states
from app.engine.futures import _bp_series, _search_support, _search_bullish
from app.engine.indicators import EMAState, TOLERANCE
from scipy.optimize import brentq

router = APIRouter()


class ComputeRequest(BaseModel):
    ema5_pre: float
    ema20_pre: float
    prev_high: float


def _make_ema5(value: float) -> EMAState:
    return EMAState(period=5, decay=2 / 6, value=value, seeded=True)


def _make_ema20(value: float) -> EMAState:
    return EMAState(period=20, decay=2 / 21, value=value, seeded=True)


def _debug_compute(ema5_pre: EMAState, ema20_pre: EMAState, prev_high: float) -> dict:
    """Run full futures chain and return intermediate values."""
    result = {
        "ema5_pre": ema5_pre.value,
        "ema20_pre": ema20_pre.value,
        "prev_high": prev_high,
        "support": None,
        "bullish": None,
        "support_bp": None,
        "bullish_bp": None,
    }

    # Support
    try:
        fl = _search_support(0.0,     ema5_pre, ema20_pre, prev_high)
        fh = _search_support(99999.0, ema5_pre, ema20_pre, prev_high)
        if fl * fh < 0:
            support = brentq(_search_support, 0.0, 99999.0,
                             args=(ema5_pre, ema20_pre, prev_high), xtol=TOLERANCE)
            result["support"] = support
            bp = _bp_series(ema5_pre, ema20_pre, [prev_high, support, support])
            result["support_bp"] = {"BP_d1": bp[0], "BP_d2": bp[1], "BP_d3": bp[2]}
    except ValueError:
        pass

    # Bullish
    try:
        fl = _search_bullish(0.0,     ema5_pre, ema20_pre)
        fh = _search_bullish(99999.0, ema5_pre, ema20_pre)
        if fl * fh < 0:
            bullish = brentq(_search_bullish, 0.0, 99999.0,
                             args=(ema5_pre, ema20_pre), xtol=TOLERANCE)
            result["bullish"] = bullish
            bp = _bp_series(ema5_pre, ema20_pre, [bullish, bullish, bullish, bullish])
            result["bullish_bp"] = {
                "BP_d1": bp[0], "BP_d2": bp[1], "BP_d3": bp[2], "BP_d4": bp[3],
            }
    except ValueError:
        pass

    return result


def _to_csv(d: dict) -> str:
    buf = io.StringIO()
    w = csv.writer(buf)
    w.writerow(["field", "value", "description"])
    w.writerow(["ema5_pre",  d["ema5_pre"],  "Settled EMA5 (after all committed bars)"])
    w.writerow(["ema20_pre", d["ema20_pre"], "Settled EMA20 (after all committed bars)"])
    w.writerow(["prev_high", d["prev_high"], "Previous day's high — d+1 anchor for support"])
    w.writerow([])
    w.writerow(["support", d["support"], "brentq: d+1=prev_high, d+2=d+3=trial, BP[d+3]=BP[d+2]"])
    if d["support_bp"]:
        for k, v in d["support_bp"].items():
            w.writerow([f"  {k}", v, ""])
    w.writerow([])
    w.writerow(["bullish", d["bullish"], "brentq: d+1..d+4=trial, BP[d+4]=BP[d+3]  (self-referential fixed point)"])
    if d["bullish_bp"]:
        for k, v in d["bullish_bp"].items():
            w.writerow([f"  {k}", v, ""])
    return buf.getvalue()


@router.post("/api/debug/compute", response_class=PlainTextResponse)
async def debug_compute(req: ComputeRequest):
    """
    Compute support and bullish from manually supplied settled EMA state.

    - ema5_pre  : EMA5 after all committed bars
    - ema20_pre : EMA20 after all committed bars
    - prev_high : previous day's high (anchor for support search)
    """
    ema5_pre  = _make_ema5(req.ema5_pre)
    ema20_pre = _make_ema20(req.ema20_pre)
    d = _debug_compute(ema5_pre, ema20_pre, req.prev_high)
    return _to_csv(d)


@router.get("/api/debug/{ticker}", response_class=PlainTextResponse)
async def debug_ticker(ticker: str):
    """
    Show the futures computation trace for the current live state of a seeded ticker.
    Uses the checkpoint (settled) EMA state and last committed bar's high.
    """
    ticker = ticker.upper()
    if ticker not in _states:
        raise HTTPException(status_code=404, detail=f"Ticker '{ticker}' not found")

    state = _states[ticker]
    cp = state._checkpoint
    if cp is None:
        raise HTTPException(status_code=400, detail="Ticker is in commit mode — seed it first")

    if not cp.bars:
        raise HTTPException(status_code=400, detail="No committed bars in state")

    prev_high = cp.bars[-1].high
    d = _debug_compute(cp.ema5.copy(), cp.ema20.copy(), prev_high)
    return _to_csv(d)
