# -*- coding: utf-8 -*-


from __future__ import print_function

import time as time
from bs4 import BeautifulSoup
import requests
import pandas as pd
from datetime import datetime

# Derivative
# https://www1.nseindia.com/live_market/dynaContent/live_watch/get_quote/getCIDHistoricalData.jsp?underlying=GBPINR&instrument=FUTCUR&expiry=27MAY2020&type=-&strike=0&fromDate=undefined&toDate=undefined&datePeriod=1day

# underlying == symbol
# instrument == FUTCUR/FURIDX/FURSTK
# expiry == ddMMMYYYY
# datePeriod == 1day/7day/2weeks/1month/3months
class Ticker():
    def __init__(self, ticker):
        self.ticker = ticker.upper()
        self.headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:76.0) Gecko/20100101 Firefox/76.0',
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
            'Accept-Language': 'en-US,en;q=0.5',
            'Accept-Encoding': 'gzip, deflate, br'
        }
        self._base_url = 'https://www1.nseindia.com/live_market/dynaContent/live_watch/get_quote'
        
    def future_cur(self, period="3months", expiry=None, type = "FUTCUR", **kwargs):
        """
        :Parameters:
            period : str
                Valid periods: 1day, 7day, 2weeks, 1month, 3months
            expiry: str
                Download start date string (DDMMMYYYY) or _datetime.
                Eg : 24May2010
            type: str
                Valid types : FUTCUR, OPTCUR ,FURIDX, OPTIDX, FUTSTK, OPTSTK
                Default is FUTCUR
            **kwargs: dict
                debug: bool
                    Optional. If passed as False, will suppress
                    error message printing to console.
        """
        if period is None or period.lower() == "max":
            if period is None:
                period = "3months"        
        params = { 'datePeriod' : period }
        
        if expiry == None:
            raise RuntimeError("Please set a valid expiry date.")
        elif isinstance(expiry, datetime):
            expiry = datetime.strftime(expiry, '%d%b%Y').upper()
        
        params['expiry'] = expiry        
        params['instrument'] = type
        params['underlying'] = self.ticker
        
        url = "{}/getCIDHistoricalData.jsp".format(self._base_url)
        data = requests.get(url, params=params ,headers=self.headers, timeout=5)
        
        if "Your request could not be processed due to technical difficulties" in data.text:
            raise RuntimeError("*** NSE IS CURRENTLY DOWN! ***\n"
                                "Thank you for your patience.")
        if "No Data" in data.text:
            raise RuntimeError(f"*** Unable to fetch data from NSE! ***\n"
                                "Please check expiry date for {}.", self.ticker)
            
        soup = BeautifulSoup(data.text, features="lxml")        
        table = soup.find_all('table')
        df = pd.read_html(str(table))[0]
        
        return df
        
        
        
        