"""
GET  /api/debug/{ticker}  — futures trace for the current live state.
"""

import io
import csv
from typing import Optional

from fastapi import APIRouter, HTTPException
from fastapi.responses import PlainTextResponse

from app.registry import _states
from app.engine.futures import _ce2, get_support, _search_bullish, _CD_DECAY
from app.engine.indicators import EMAState

router = APIRouter()


@router.get("/api/debug/{ticker}", response_class=PlainTextResponse)
async def debug_ticker(ticker: str):
    ticker = ticker.upper()
    if ticker not in _states:
        raise HTTPException(status_code=404, detail=f"Ticker '{ticker}' not found")

    state = _states[ticker]
    if not state.history:
        raise HTTPException(status_code=400, detail="No data loaded for this ticker")
    if not state.cd_ema.seeded:
        raise HTTPException(status_code=400, detail="CD state not seeded yet (need ~105 bars)")

    # Use pre-bar state for ema5/ema20 so futures match the day's pinned values
    ps = state._pre_bar_state
    if ps is None:
        raise HTTPException(status_code=400, detail="No pre-bar state available")
    e5 = ps.ema5.copy()
    e20 = ps.ema20.copy()
    cd_pre = ps.cd_ema.value

    ce2 = _ce2(e5, e20)
    cd_curr = _CD_DECAY * (ce2 - cd_pre) + cd_pre
    support = get_support(e5, e20, cd_curr)

    # Bullish: 2-day convergence
    from scipy.optimize import brentq
    from app.engine.indicators import TOLERANCE
    bullish = None
    try:
        f_lo = _search_bullish(0.0, e5, e20, cd_pre)
        f_hi = _search_bullish(99999.0, e5, e20, cd_pre)
        if f_lo * f_hi < 0:
            bullish = brentq(_search_bullish, 0.0, 99999.0, args=(e5, e20, cd_pre), xtol=TOLERANCE)
    except Exception:
        pass

    def check_trial(w):
        ce2_d1 = _ce2(e5, e20)
        cd_d1 = _CD_DECAY * (ce2_d1 - cd_pre) + cd_pre
        e5_d1 = e5.copy(); e5_d1.update(w)
        e20_d1 = e20.copy(); e20_d1.update(w)
        ce2_d2 = _ce2(e5_d1, e20_d1)
        cd_d2 = _CD_DECAY * (ce2_d2 - cd_d1) + cd_d1
        s2 = get_support(e5_d1, e20_d1, cd_d2)
        return s2

    buf = io.StringIO()
    w = csv.writer(buf)
    w.writerow(["Debug trace for", ticker])
    w.writerow(["EMA5", f"{e5.value:.4f}", "EMA20", f"{e20.value:.4f}", "CD_pre", f"{cd_pre:.4f}"])
    w.writerow(["CE2", f"{ce2:.4f}", "CD_curr", f"{cd_curr:.4f}"])
    w.writerow(["Support", f"{support:.4f}" if support else "None"])
    w.writerow(["Bullish", f"{bullish:.4f}" if bullish else "None"])
    w.writerow([])
    if bullish:
        s2 = check_trial(bullish)
        w.writerow(["Verify Bullish W:", f"{bullish:.4f}", "Support(Day2):", f"{s2:.4f}" if s2 else "None"])
    return buf.getvalue()
