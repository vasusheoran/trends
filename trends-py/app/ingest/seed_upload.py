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

    # Columns B=2, W=23, X=24, Y=25, Z=26 (1-based)
    # iter_rows streams the file in one pass — much faster than per-cell access
    for row in ws.iter_rows(min_col=2, max_col=26, values_only=True):
        # row[0]=B, row[21]=W, row[22]=X, row[23]=Y, row[24]=Z
        date  = row[0]
        close = row[21]
        open_ = row[22]
        high  = row[23]
        low   = row[24]

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
