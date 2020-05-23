# -*- coding: utf-8 -*-

from __future__ import print_function

from .base import Base
import requests
from datetime import datetime

class Ticker(Base):
    def __init__(self, ticker):
        super()
        self.ticker = ticker.upper()
        
    def future_cur(self, period="3months", expiry=None, instrument = "FUTCUR", include=None, **kwargs):
        """
        :Parameters:
            period : str
                Valid periods: 1day, 7day, 2weeks, 1month, 3months
            expiry: str
                Download start date string (DDMMMYYYY) or _datetime.
                Eg : 24May2010
            instrument: str
                Valid instrument : FUTCUR, OPTCUR ,FURIDX, OPTIDX, FUTSTK, OPTSTK
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
        params['instrument'] = instrument
        params['underlying'] = self.ticker
        
        print(params)
        
        url = self.get_url(instrument)
        data = requests.get(url, params=params ,headers=self.headers, timeout=5)
        
        if "Your request could not be processed due to technical difficulties" in data.text:
            raise RuntimeError("*** NSE IS CURRENTLY DOWN! ***\n"
                                "Thank you for your patience.")
        if "No Data" in data.text:
            raise RuntimeError(f"*** Unable to fetch data from NSE! ***\n"
                                "Please check expiry date for {}.", self.ticker)
            
        df = self.parse(data, include)
        
        return df
            