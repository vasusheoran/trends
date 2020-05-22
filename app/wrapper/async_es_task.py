# -*- coding: utf-8 -*-
"""
Created on Tue Mar 31 00:41:53 2020

@author: vsheoran
"""
import threading
from DB import DB
import pandas as pd
from utilities import get_logger

logger = get_logger("async_update_task.py")
        
class AsyncUpdateSymbolsTask(threading.Thread):
    def __init__(self, task_details):
        super().__init__()
        self.task_details = task_details
    def run(self):
        df = pd.read_excel(f"files/{self.task_details}")
        df.columns = df.columns.str.replace(' ', '')
        
        df['SASSymbol'] = df['SASSymbol'].str.lower()
        df['SASSymbol'] = df['SASSymbol'].replace(' ', '_', regex=True)
        df['YahooSymbol'] = df['YahooSymbol'].str.lower()
        
        df.set_index("SASSymbol", inplace = True)
        
        db = DB()
        db.post_symbols_client(df)

class AsyncUpdateHistoricalTask(threading.Thread):
    previousRequest = None
    def __init__(self, task_details):
        super().__init__()
        self.task_details = task_details
    def run(self):
        # try:
        if 'index' in self.task_details.keys():
            index = self.task_details['index']
        logger.info("Index : " + index)
        self.process_index(index, self.task_details['data'])
        # except Exception as ex:
        #     logger.error('Unable to post to client')
        #     logger.info(ex)
            
        
    def process_index(self, index, data):
        db = DB()
        logger.info("Posting data to es : " + index)
        db.post_to_client(data, index, False)
