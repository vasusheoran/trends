# -*- coding: utf-8 -*-

# Derivative
# https://www1.nseindia.com/live_market/dynaContent/live_watch/get_quote/getCIDHistoricalData.jsp?underlying=GBPINR&instrument=FUTCUR&expiry=27MAY2020&type=-&strike=0&fromDate=undefined&toDate=undefined&datePeriod=1day

# underlying == symbol
# instrument == FUTCUR/FURIDX/FURSTK
# expiry == ddMMMYYYY
# datePeriod == 1day/7day/2weeks/1month/3months

from bs4 import BeautifulSoup
import pandas as pd

class Base():
    base_url = 'https://www1.nseindia.com/live_market/dynaContent/live_watch/get_quote'
    headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:76.0) Gecko/20100101 Firefox/76.0',
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
            'Accept-Language': 'en-US,en;q=0.5',
            'Accept-Encoding': 'gzip, deflate, br'
        }
    
    def get_url(self, instrument):
        if "IDX" in instrument:
            url = "{}/getFOHistoricalData.jsp".format(self.base_url)
        elif "CUR" in instrument:
            url = "{}/getCIDHistoricalData.jsp".format(self.base_url)
        elif "STK" in instrument:
            url = "{}/getFOHistoricalData.jsp".format(self.base_url)
            
        return url
    
    def get_qoute_url(self, instrument):
        if "CUR" in instrument:
            url = "{}/ajaxCDGetQuoteDataTest.jsp".format(self.base_url)
            
        return url
            
            
    def parse(self, data, include=None):
        soup = BeautifulSoup(data.text, features="lxml")        
        table = soup.find_all('table')
        df = pd.read_html(str(table))[0]
        
        if include != None and isinstance(include, list):
            df = df[include]
        
        return df
        
        
        