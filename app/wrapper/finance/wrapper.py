# -*- coding: utf-8 -*-

from . import Ticker
from .util import get_logger
import pandas as pd
from .yahoo import historical_data

logger = get_logger('fin-wrapper.py')

def derivative(ticker, expiry, instrument, period="3months", option="-", strike = "0",  include = ["Date", "Close", "High", "Low"]):
    # try:
    logger.info("Inside derivative")
    tk = Ticker(ticker)
    
    df = tk.currency_derivatives(period=period, expiry=expiry, instrument=instrument, option=option, strike = strike, include=include)
    logger.info("Exiting derivative")
    return reset(df)

def reset(df):
    df = df.rename(columns={"High": "HP", "Close": "CP", "Low":"LP"})
    return df

def history(ticker, period = "2y", interval = "1d"):
    logger.info("Inside history")
    df = historical_data(listing = ticker, period = period, interval = interval)
    # df=pd.read_csv("/Users/vasusheoran/git/trends/app/data/daily/Nifty 50/2021-09-01.csv")
    logger.info("Exiting history")
    return df
