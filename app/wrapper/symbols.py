# -*- coding: utf-8 -*-


from .utilities import get_logger
from .database import DB

logger = get_logger("symbols.py")

class Symbol:
    def __init__(self):
        self.db = DB()
    
    def get(self):
        logger.info("Getting symbols list")
        return self.db.get_listings()
    
    def put(self, symbol, symbol_id):
        raise RuntimeError("Implementation pending")
    
    def post(self, symbol):
        raise RuntimeError("Implementation pending")
        
    def delete(self, symbol_id):
        raise RuntimeError("Implementation pending")