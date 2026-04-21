"""
POST /api/seed/{ticker} — upload an Excel file to seed historical state.

Column detection: scans every row for a header row containing Date, Close, Open, High, Low
(case-insensitive). Data rows start immediately after the detected header.

Fallback: if no header row is found, assumes the fixed layout used in the dev Excel:
  B=Date, W=Close, X=Open, Y=High, Z=Low

Rows where any of those five values are missing or non-numeric are silently skipped.
Resets existing ticker state before loading.
"""

import openpyxl
from io import BytesIO
from typing import Optional

from fastapi import APIRouter, File, HTTPException, UploadFile

from app.registry import reset_state

router = APIRouter()

_REQUIRED = {"date", "close", "open", "high", "low"}
# Fallback fixed column indices (0-based within iter_rows min_col=2)
_FALLBACK = {"date": 0, "close": 21, "open": 22, "high": 23, "low": 24}


def _detect_header(ws) -> Optional[dict]:
    """
    Scan up to the first 20 rows for a header row containing all required column names.
    Returns a dict mapping field name → 0-based column index within the full row (col 1 = index 0).
    Returns None if no header found.
    """
    for row in ws.iter_rows(max_row=20, values_only=True):
        normalized = [str(c).strip().lower() if c is not None else "" for c in row]
        if _REQUIRED.issubset(set(normalized)):
            return {field: normalized.index(field) for field in _REQUIRED}
    return None


@router.post("/api/seed/{ticker}")
async def seed_ticker(ticker: str, file: UploadFile = File(...)):
    """
    Upload an Excel (.xlsx) file to seed historical OHLCV data for a ticker.

    The active (first) sheet is used. Column detection works in two ways:
    1. Header scan: looks for a row in the first 20 rows containing headers named
       Date, Close, Open, High, Low (case-insensitive). Data starts after that row.
    2. Fallback: if no header found, assumes fixed layout B=Date, W=Close, X=Open, Y=High, Z=Low.

    Rows with missing or non-numeric OHLCV values are silently skipped.
    Resets existing ticker state — all prior bars are discarded.
    """
    contents = await file.read()
    try:
        wb = openpyxl.load_workbook(BytesIO(contents), data_only=True, read_only=True)
    except Exception:
        raise HTTPException(status_code=400, detail="Could not parse file — make sure it is a valid .xlsx file")

    ws = wb.active
    col_map = _detect_header(ws)
    using_fallback = col_map is None

    if using_fallback:
        col_map = _FALLBACK

    state = reset_state(ticker)
    count = 0
    header_found = False

    for row in ws.iter_rows(values_only=True):
        # In fallback mode, skip until we're past any header-like rows by checking numeric values
        date  = row[col_map["date"]]  if col_map["date"]  < len(row) else None
        close = row[col_map["close"]] if col_map["close"] < len(row) else None
        open_ = row[col_map["open"]]  if col_map["open"]  < len(row) else None
        high  = row[col_map["high"]]  if col_map["high"]  < len(row) else None
        low   = row[col_map["low"]]   if col_map["low"]   < len(row) else None

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

    # Switch to live mode — subsequent PUTs restore from this checkpoint,
    # making them idempotent for the same date/OHLCV.
    state.commit()

    return {
        "ticker": ticker,
        "bars_loaded": count,
        "status": "seeded" if count > 0 else "no_data",
        "column_detection": "fallback (B/W/X/Y/Z)" if using_fallback else "header row detected",
    }
