# -*- coding: utf-8 -*-

from . import Ticker
from .yahoo import historical_data

def derivative(ticker=None, period="3months", expiry=None, instrument = "FUTCUR", include = ["Date", "Close", "High", "Low"]):
    ticker = Ticker(ticker)
    df = ticker.future_cur(period=period, expiry=expiry, instrument=instrument, include=include)
    return reset(df)


def reset(df):
    df = df.rename(columns={"High": "HP", "Close": "CP", "Low":"LP"})
    return df

def history(ticker, period = "2y", interval = "1d"):
    df = historical_data(listing = ticker, period = period, interval = interval)
    return df