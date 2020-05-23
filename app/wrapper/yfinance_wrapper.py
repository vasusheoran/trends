# -*- coding: utf-8 -*-
"""
Created on Mon Mar 30 13:44:04 2020

@author: vsheoran
"""

import yfinance as yf
import pandas as pd

def get_historical_data_by_period(listing, period, interval):
    df = yf.download(listing, period=period, interval=interval)    
    return reset_columns(pd.DataFrame(df))

def get_historical_data_by_time(listing, interval, start, end):
    df = yf.download(listing, start=start, end=end, interval=interval)
    return reset_columns(pd.DataFrame(df))
    


def reset_columns(data):
    data = data.rename(columns={"High": "HP", "Close": "CP", "Low":"LP"})
    
    
    data = data[['CP', 'HP', 'LP']].reset_index()
    
    # if data.at[len(data) - 1 ,'Date'].date() == datetime.now().date():
    #     data.drop(data.tail(1).index,inplace=True)
        
    return data

#data = get_historical_data("nse", "1y", "1mo", "in")
