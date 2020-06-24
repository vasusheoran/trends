# -*- coding: utf-8 -*-
"""
Created on Mon Mar 30 13:44:04 2020

@author: vsheoran
"""

import yfinance as yf
import pandas as pd
from datetime import datetime

def historical_data(listing, period, interval):
    df = yf.download(listing, period=period, interval=interval)    
    return reset_columns(pd.DataFrame(df))

def historical_data_by_time(listing, interval, start, end):
    df = yf.download(listing, start=start, end=end, interval=interval)
    return reset_columns(pd.DataFrame(df))   


def reset_columns(data):
    data = data.rename(columns={"High": "HP", "Close": "CP", "Low":"LP"})
    data = data[['CP', 'HP', 'LP']].reset_index()
    
    # Remove current day if exists
    
    todays_date = str(datetime.now().date())
    valid_dates = data["Date"] < todays_date
    
    filtered_dates = data.loc[valid_dates]
    return filtered_dates

#data = get_historical_data("nse", "1y", "1mo", "in")
