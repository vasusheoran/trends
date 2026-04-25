"""
Shared fixtures — loads Excel computed values for test assertions.
Uses data_only=True so we get the last-saved computed values, not formulas.
Source of truth: data/Final-bullish-ce.xlsx (Nifty-20.12.2024 sheet)
"""

import pytest
import openpyxl
from pathlib import Path

EXCEL_PATH = Path(__file__).parent.parent.parent / "data" / "Final-Bullish-CE.xlsx"
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
    wb_data    = openpyxl.load_workbook(EXCEL_PATH, data_only=True)
    wb_formula = openpyxl.load_workbook(EXCEL_PATH, data_only=False)
    ws_d = wb_data[SHEET_NAME]
    ws_f = wb_formula[SHEET_NAME]
    rows = []
    for row in range(DATA_START, DATA_END + 1):
        date = ws_d[f"B{row}"].value
        close = ws_d[f"W{row}"].value
        open_ = ws_d[f"X{row}"].value
        high = ws_d[f"Y{row}"].value
        low = ws_d[f"Z{row}"].value
        if not all([date, close, open_, high, low]):
            continue
        if hasattr(date, "strftime"):
            date = date.strftime("%d-%b-%Y")

        # Only trust HL when the Excel formula is the expected MIN(prev 3 highs) pattern.
        # Some rows have hardcoded values or broken formula references — skip those.
        hl_formula = ws_f[f"AD{row}"].value
        hl = ws_d[f"AD{row}"].value if (isinstance(hl_formula, str) and hl_formula.upper().startswith("=MIN(Y")) else None

        rows.append({
            "row": row,
            "date": str(date),
            "close": float(close),
            "open": float(open_),
            "high": float(high),
            "low": float(low),
            "hl":      hl,
            "avg":     ws_d[f"AR{row}"].value,
            "ema5":    ws_d[f"AS{row}"].value,
            "ema20":   ws_d[f"BN{row}"].value,   # BN = EMA-20 (2/21) in new file
            "rsi":     ws_d[f"BV{row}"].value,
            "support": ws_d[f"CC{row}"].value,   # CC = Support (binary search result)
            "bullish": ws_d[f"BR{row}"].value,   # BR = Bullish (binary search result)
        })
    wb_data.close()
    wb_formula.close()
    return rows
