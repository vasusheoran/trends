# -*- coding: utf-8 -*-


from .utilities import get_logger
from .database import DB

logger = get_logger("symbols.py")

class Symbol:
    def __init__(self):
        self.db = DB()
    
    def get(self):
        logger.info("Handlling get")
        return self.db.get_listings()
    
    def put(self, symbol, sid):
        logger.info("Handlling put")
        symbols = self.db.get_listings_df()
        # Delete existing mapping
        symbols = symbols[symbols['Symbol'] != sid]
        # Append new mapping
        symbols = symbols.append(symbol, ignore_index=True)        
        self.db.set_listings(symbols)
    
    def post(self, symbol):
        logger.info("Handlling post")
        symbols = self.db.get_listings_df()
        symbols = symbols.append(symbol, ignore_index=True)
        self.db.set_listings(symbols)
        
    def delete(self, sid):
        logger.info("Handlling delete")
        symbols = self.db.get_listings_df()
        symbols = symbols[symbols['Symbol'] != sid]
        self.db.set_listings(symbols)