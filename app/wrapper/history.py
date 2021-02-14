# -*- coding: utf-8 -*-
import datetime
import pandas as pd

from .utilities import get_logger
from .database import DB
from .index import index

logger = get_logger("history.py")

class History:
    def __init__(self):
        self.db = DB()
        self.symbol = None
    
    def get(self):
        logger.info("Handlling get")
        self.symbol = index.symbol()
        self.db = DB(listing=self.symbol['SAS'])

        data = self.db.get_historical_data_list()
        data = data[::-1]
        return data
    
    def put(self, history, sid):
        logger.info("Handlling put")
        his = self.db.get_historical_data_csv()
        # Delete existing mapping
        his = his[his['Date'] != sid]
        # Append new mapping
        his = his.append(history, ignore_index=True)        
        self.db.set_historical_data(his)
        index.refresh()
    
    def post(self, history):
        logger.info("Handlling refresh")
        his = self.db.get_historical_data_csv()
        his = his.append(history, ignore_index=True)
        self.db.set_historical_data(his)
        index.refresh()

    def delete(self, sid):
        logger.info("Handlling delete")
        his = self.db.get_historical_data_csv()
        his = his[his['Date'] != sid]
        self.db.set_historical_data(his)
        index.refresh()

    def patch(self, date):
        logger.info("Handlling patch")

        self.symbol = index.symbol()
        self.db = DB(listing=self.symbol['SAS'])

        his_today = self.db.get()
        his_old = self.db.get(date)
        result = pd.concat([his_old, his_today]).drop_duplicates().reset_index(drop=True)
        self.db.set_historical_data(result)
        index.refresh()
