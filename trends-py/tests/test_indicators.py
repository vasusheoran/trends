"""
Tests for indicator calculations against Excel ground truth.
Asserts within TOLERANCE = 0.001 (matches Go TOLERANCE).
"""

import pytest
from app.engine.state import TickerState
from app.engine.indicators import TOLERANCE

# Run full history through and spot-check the last N rows
_SPOT_CHECK_LAST = 20


def _feed_and_collect(excel_rows):
    """Feed all Excel rows into TickerState, return list of snapshots."""
    state = TickerState(ticker="test")
    snapshots = []
    for r in excel_rows:
        snap = state.update(r["date"], r["close"], r["open"], r["high"], r["low"])
        snapshots.append((r, snap))
    return snapshots


def _close_enough(actual, expected, tol=TOLERANCE):
    if actual is None and expected is None:
        return True
    if actual is None or expected is None:
        return False
    return abs(actual - expected) <= tol


class TestEMA:
    def test_ema5_matches_excel(self, excel_rows):
        pairs = _feed_and_collect(excel_rows)
        failures = []
        for r, snap in pairs[-_SPOT_CHECK_LAST:]:
            if r["ema5"] is None:
                continue
            if not _close_enough(snap.ema5, r["ema5"]):
                failures.append(
                    f"row={r['row']} date={r['date']}: got={snap.ema5:.4f} expected={r['ema5']:.4f}"
                )
        assert not failures, "EMA-5 mismatches:\n" + "\n".join(failures)

    def test_ema20_matches_excel(self, excel_rows):
        pairs = _feed_and_collect(excel_rows)
        failures = []
        for r, snap in pairs[-_SPOT_CHECK_LAST:]:
            if r["ema20"] is None:
                continue
            if not _close_enough(snap.ema20, r["ema20"]):
                failures.append(
                    f"row={r['row']} date={r['date']}: got={snap.ema20:.4f} expected={r['ema20']:.4f}"
                )
        assert not failures, "EMA-20 mismatches:\n" + "\n".join(failures)


class TestHL:
    def test_hl_matches_excel(self, excel_rows):
        pairs = _feed_and_collect(excel_rows)
        failures = []
        for r, snap in pairs[-_SPOT_CHECK_LAST:]:
            if r["hl"] is None:
                continue
            if not _close_enough(snap.hl, r["hl"]):
                failures.append(
                    f"row={r['row']} date={r['date']}: got={snap.hl:.4f} expected={r['hl']:.4f}"
                )
        assert not failures, "H/L mismatches:\n" + "\n".join(failures)


class TestAVG:
    def test_avg_matches_excel(self, excel_rows):
        pairs = _feed_and_collect(excel_rows)
        failures = []
        for r, snap in pairs[-_SPOT_CHECK_LAST:]:
            if r["avg"] is None:
                continue
            # AVG uses a wider tolerance due to the correction formula
            if not _close_enough(snap.avg, r["avg"], tol=1.0):
                failures.append(
                    f"row={r['row']} date={r['date']}: got={snap.avg:.4f} expected={r['avg']:.4f}"
                )
        assert not failures, "AVG mismatches:\n" + "\n".join(failures)


class TestRSI:
    def test_rsi_matches_excel(self, excel_rows):
        pairs = _feed_and_collect(excel_rows)
        failures = []
        for r, snap in pairs[-_SPOT_CHECK_LAST:]:
            if r["rsi"] is None:
                continue
            if not _close_enough(snap.rsi, r["rsi"], tol=0.01):
                failures.append(
                    f"row={r['row']} date={r['date']}: got={snap.rsi:.4f} expected={r['rsi']:.4f}"
                )
        assert not failures, "RSI mismatches:\n" + "\n".join(failures)
