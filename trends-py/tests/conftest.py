"""
Shared fixtures — loads Excel computed values for test assertions.
Uses data_only=True so we get the last-saved computed values, not formulas.
Source of truth: data/Final-bullish-ce.xlsx (Nifty-20.12.2024 sheet)
"""

import pytest
import openpyxl
from pathlib import Path

EXCEL_PATH = Path(__file__).parent.parent.parent / "data" / "Final-bullish-ce.xlsx"
SHEET_NAME = "Nifty-20.12.2024"
DATA_START = 5
DATA_END = 5387


@pytest.fixture(scope="session")
def excel_rows():
    """
    Returns a list of dicts with computed Excel values for each data row.
    Column mapping (Final-bullish-ce.xlsx):
      B=Date, W=Close, X=Open, Y=High, Z=Low
      AD=H/L, AR=AVG, AS=EMA-5, BN=EMA-20 (decay 2/21), BV=RSI(14)
    """
    wb = openpyxl.load_workbook(EXCEL_PATH, data_only=True)
    ws = wb[SHEET_NAME]
    rows = []
    for row in range(DATA_START, DATA_END + 1):
        date = ws[f"B{row}"].value
        close = ws[f"W{row}"].value
        open_ = ws[f"X{row}"].value
        high = ws[f"Y{row}"].value
        low = ws[f"Z{row}"].value
        if not all([date, close, open_, high, low]):
            continue
        if hasattr(date, "strftime"):
            date = date.strftime("%d-%b-%Y")
        rows.append({
            "row": row,
            "date": str(date),
            "close": float(close),
            "open": float(open_),
            "high": float(high),
            "low": float(low),
            "hl":    ws[f"AD{row}"].value,
            "avg":   ws[f"AR{row}"].value,
            "ema5":  ws[f"AS{row}"].value,
            "ema20": ws[f"BN{row}"].value,   # BN = EMA-20 (2/21) in new file
            "rsi":   ws[f"BV{row}"].value,
        })
    wb.close()
    return rows
