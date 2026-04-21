"""
POST /api/seed/{ticker} — upload an Excel file to seed historical state.

Reads only: B=Date, W=Close, X=Open, Y=High, Z=Low.
All other columns are ignored — no indicator columns required.
Scans all rows, skips any row where those five cells are not all present.
Resets existing ticker state before loading.
"""

import openpyxl
from io import BytesIO

from fastapi import APIRouter, File, HTTPException, UploadFile

from app.registry import reset_state

router = APIRouter()


@router.post("/api/seed/{ticker}")
async def seed_ticker(ticker: str, file: UploadFile = File(...)):
    """
    Upload an Excel (.xlsx) file to seed historical OHLCV data for a ticker.
    The active (first) sheet is used. Columns must follow the standard layout:
      B=Date, W=Close, X=Open, Y=High, Z=Low
    Rows with any missing value in those columns are silently skipped.
    Resets existing ticker state — all prior bars are discarded.
    """
    contents = await file.read()
    try:
        wb = openpyxl.load_workbook(BytesIO(contents), data_only=True, read_only=True)
    except Exception:
        raise HTTPException(status_code=400, detail="Could not parse file — make sure it is a valid .xlsx file")

    ws = wb.active
    state = reset_state(ticker)
    count = 0

    for row_idx in range(1, (ws.max_row or 0) + 1):
        date  = ws[f"B{row_idx}"].value
        close = ws[f"W{row_idx}"].value
        open_ = ws[f"X{row_idx}"].value
        high  = ws[f"Y{row_idx}"].value
        low   = ws[f"Z{row_idx}"].value

        if not all([date, close, open_, high, low]):
            continue
        try:
            close = float(close)
            open_ = float(open_)
            high  = float(high)
            low   = float(low)
        except (TypeError, ValueError):
            continue

        if hasattr(date, "strftime"):
            date = date.strftime("%d-%b-%Y")

        state.update(str(date), close, open_, high, low)
        count += 1

    wb.close()

    return {
        "ticker": ticker,
        "bars_loaded": count,
        "status": "seeded" if count > 0 else "no_data",
    }
