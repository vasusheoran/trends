# -*- coding: utf-8 -*-

from __future__ import print_function

from .base import Base
from .util import get_logger
import requests
from datetime import datetime

logger = get_logger('Ticker.py')

class Ticker(Base):
    def __init__(self, ticker):
        super()
        if not ticker:
            raise AttributeError("ticker is a required field")
        self.ticker = ticker.upper()
        self.logger = get_logger('Ticker | finance')
        
    def currency_derivatives(self, expiry, instrument, period="3months", option=None, strike=None, include=None, **kwargs):
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
        
        logger.info("Inside currency_derivatives")
        
        if expiry == None:
            raise RuntimeError("Please set a valid expiry date.")
        elif isinstance(expiry, str):
            expiry = expiry.upper()
        elif isinstance(expiry, datetime):
            expiry = datetime.strftime(expiry, '%d%b%Y').upper()
            expiry = datetime.strftime(expiry, '%d%b%Y').upper()
        
        params = {}
        params['underlying'] = self.ticker
        params['instrument'] = instrument
        params['expiry'] = expiry      
        
        
        if option:
            params['type'] = option
        else:
            params['type'] = "-"  
            
        if strike:
            params['strike'] = strike 
        else:
            params['strike'] = "0"
            
            
              
        params['fromDate'] = "undefined"
        params['toDate'] = "undefined"
            
        if period is None or period.lower() == "max":
            if period is None:
                period = "3months"  
                
        params['datePeriod'] = period                
        
        url = self.get_url(instrument)     
        print(params)      
        
        
        
        # params['underlying'] = "USDINR"
        # params['instrument'] = "FUTCUR"
        # params['expiry'] = "26JUN2020"     
        # params['type'] = "-"
        # params['strike'] = "0"
        # params['fromDate'] = "undefined"      
        # params['toDate'] = "undefined"
        # params['datePeriod'] = "3months"
        
        data = requests.get(url, params=params ,headers=self.headers)
        
        if "Your request could not be processed due to technical difficulties" in data.text:
            raise RuntimeError("*** NSE IS CURRENTLY DOWN! ***\n Thank you for your patience.")
        if "No Data" in data.text:
            raise RuntimeError(f"*** Unable to fetch data from NSE! ***\n Please check expiry date for {self.ticker}.")
        logger.info(data)
        df = self.parse(data, include)
        
        return df
        
    def options(self, expiry, instrument, option, strike, period="3months", include=None, **kwargs):
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
        elif isinstance(expiry, str):
            expiry = expiry.upper()
        elif isinstance(expiry, datetime):
            expiry = datetime.strftime(expiry, '%d%b%Y').upper()
            expiry = datetime.strftime(expiry, '%d%b%Y').upper()
            
        params['expiry'] = expiry   
        params['instrument'] = instrument
        params['underlying'] = self.ticker
        params['type'] = option
        params['strike'] = strike        
        
        url = self.get_url(instrument)     
# =============================================================================
#         print(params)   
#         print(url)
# =============================================================================
        
        data = requests.get(url, params=params ,headers=self.headers, timeout=15)
        
        if "Your request could not be processed due to technical difficulties" in data.text:
            raise RuntimeError("*** NSE IS CURRENTLY DOWN! ***\n Thank you for your patience.")
        if "No Data" in data.text:
            raise RuntimeError(f"*** Unable to fetch data from NSE! ***\n Please check expiry date for {self.ticker}.")
            
        df = self.parse(data, include)
        
        logger.info("Exiting currency_derivatives")
        return df
    
    def quotes(self, instrument = "OPTCUR", expiry=None, optionType = None):
        """
        :Parameters:
            instrument: str
                Valid instrument : FUTCUR, OPTCUR ,FURIDX, OPTIDX, FUTSTK, OPTSTK
                Default is FUTCUR
                
        """
        params = { 'u' : self.ticker}
        
        params['i'] = instrument  
        
        if expiry:
            params['e'] = expiry
        if optionType:
            params['o'] = optionType
            params['k'] = optionType
        
        url = self.get_qoute_url(instrument)
        data = requests.get(url, params=params ,headers=self.headers)
        
        if "Your request could not be processed due to technical difficulties" in data.text:
            raise RuntimeError("*** NSE IS CURRENTLY DOWN! ***\n Thank you for your patience.")
        if "FAILIURE" in data.text:
            raise RuntimeError(f"*** Unable to fetch data from NSE! ***\n Please check expiry date for {self.ticker}." )
        
        try:
            json = data.json()
            return json
        except Exception as ex:
            print("Error in Ticker ..." + ex)
            return {}
            