# -*- coding: utf-8 -*-

from . import Ticker, util

def download(ticker=None, period="3months", expiry="27MAY2020", type = "FUTCUR"):
    ticker = Ticker(ticker)
    df = ticker.future_cur(period=period, expiry=expiry, type=type)
    return df