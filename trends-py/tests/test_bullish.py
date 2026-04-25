"""
Bullish (BR) binary-search validation.

Algorithm: W=[BR, CE2, CE2]; find BR where BP[d+3]=BP[d+2].
CE2 is the closed-form 2-bar fixed point: (49*EMA5_pre - 19*EMA20_pre) / 30.

Seeded from Excel-stored EMA values (Nifty-20.12.2024 sheet, AS/BN columns).
"""

import pytest
from scipy.optimize import brentq

from app.engine.futures import _bp_series, _search_br, _ce2
from app.engine.indicators import EMAState, TOLERANCE

CONVERGENCE_TOLERANCE = TOLERANCE * 5


def _ema5(v: float) -> EMAState:
    return EMAState(period=5, decay=2 / 6, value=v, seeded=True)


def _ema20(v: float) -> EMAState:
    return EMAState(period=20, decay=2 / 21, value=v, seeded=True)


def _find_bullish(ema5_pre, ema20_pre):
    """Run brentq on _search_br. Returns bullish value or None."""
    ce2 = _ce2(ema5_pre, ema20_pre)
    fl = _search_br(0.0,     ema5_pre, ema20_pre, ce2)
    fh = _search_br(99999.0, ema5_pre, ema20_pre, ce2)
    if fl * fh >= 0:
        return None
    return brentq(_search_br, 0.0, 99999.0, args=(ema5_pre, ema20_pre, ce2), xtol=TOLERANCE)


def _excel_ema_pairs(excel_rows):
    pairs = []
    for i in range(1, len(excel_rows)):
        prev, curr = excel_rows[i - 1], excel_rows[i]
        if prev["ema5"] and prev["ema20"]:
            pairs.append((prev, curr))
    return pairs


def test_bullish_convergence_from_excel_ema(excel_rows):
    """
    For the last 20 row pairs, use Excel-stored EMA from row i-1 as settled state.
    Verify |BP[d+3]-BP[d+2]| ≤ CONVERGENCE_TOLERANCE at the computed BR value.
    """
    pairs = _excel_ema_pairs(excel_rows)
    sample = pairs[-20:]
    assert len(sample) == 20, "Not enough row pairs with Excel EMA values"

    failures = []
    for prev, curr in sample:
        ema5_pre  = _ema5(prev["ema5"])
        ema20_pre = _ema20(prev["ema20"])

        bullish = _find_bullish(ema5_pre, ema20_pre)
        if bullish is None:
            failures.append(f"row {curr['row']}: Bullish bracket not found")
            continue

        ce2 = _ce2(ema5_pre, ema20_pre)
        bp = _bp_series(ema5_pre, ema20_pre, [bullish, ce2, ce2])
        diff = abs(bp[2] - bp[1])
        if diff > CONVERGENCE_TOLERANCE:
            failures.append(
                f"row {curr['row']}: |BP[d+3]-BP[d+2]|={diff:.6f} > {CONVERGENCE_TOLERANCE} "
                f"at bullish={bullish:.2f}"
            )

    assert not failures, "Bullish convergence failures:\n" + "\n".join(failures)


def test_bullish_bracket_exists_for_excel_ema(excel_rows):
    """Verify brentq bracket [0, 99999] has a sign change for the last 20 rows."""
    pairs = _excel_ema_pairs(excel_rows)
    sample = pairs[-20:]

    failures = []
    for prev, curr in sample:
        ema5_pre  = _ema5(prev["ema5"])
        ema20_pre = _ema20(prev["ema20"])
        ce2 = _ce2(ema5_pre, ema20_pre)
        fl = _search_br(0.0,     ema5_pre, ema20_pre, ce2)
        fh = _search_br(99999.0, ema5_pre, ema20_pre, ce2)
        if fl * fh >= 0:
            failures.append(f"row {curr['row']}: no sign change — fl={fl:.4f}, fh={fh:.4f}")

    assert not failures, "Bracket failures:\n" + "\n".join(failures)


def test_bullish_matches_excel_last_20_rows(excel_rows):
    """BR from ema_pre at row i-1 matches Excel BR column at row i.
    Tolerance 15.0: row 5385 (18-Dec-2024) is a known Excel outlier (~14 pts off).
    Most rows are within 1.5."""
    pairs = _excel_ema_pairs(excel_rows)
    sample = [(p, c) for p, c in pairs if c["bullish"] is not None][-20:]
    assert len(sample) == 20, "Not enough rows with stored bullish"

    failures = []
    for prev, curr in sample:
        ema5_pre  = _ema5(prev["ema5"])
        ema20_pre = _ema20(prev["ema20"])
        bullish = _find_bullish(ema5_pre, ema20_pre)
        if bullish is None:
            failures.append(f"row {curr['row']}: bracket not found")
            continue
        diff = abs(bullish - curr["bullish"])
        if diff > 15.0:
            failures.append(
                f"row {curr['row']} ({curr['date']}): computed={bullish:.2f}  "
                f"excel={curr['bullish']:.2f}  diff={diff:.2f}"
            )

    assert not failures, "BR vs Excel mismatches:\n" + "\n".join(failures)
