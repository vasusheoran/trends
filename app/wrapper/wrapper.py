# -*- coding: utf-8 -*-
"""
Created on Mon Mar 23 19:09:43 2020

@author: vsheoran
"""


#from lib.calculate import Calculate
import pandas as pd
from .index import index
from .utilities import get_logger, Utilities
from .finance import Ticker
from .database import DB
from flask_socketio import emit

logger = get_logger("wrapper.py")

def set_up_socketio(io):
    global socketio
    socketio = io
    

#  TO be done in wrapper.index.py
def reset_current_index():
   return RuntimeError("Pending Implementation")
        
    
def push_notifications(eventName, update):
    logger.info("Inside push_notifications")
    socketio.emit(eventName, update)
    logger.info("Exiting push_notifications")
    
def paginate(page, size):
    offset = 3
    df = index.get_df()
        
    if size > 0 and page > 0:
        if page > 1:
            start_index = ((page - 1) * size) + offset
        else:
            start_index = 0 + offset
        end_index = ((page - 1) * size) + size + offset
        resp = df.iloc[start_index : end_index]
        resp =resp[['CP', 'HP', 'LP', 'Date']]
        resp['Date'] = pd.to_datetime(resp['Date'])
        return resp.to_dict('records')
    return []

utils = Utilities()
db = DB()

def fetch_expires(data):
    logger.info("Inside fetch_expires")
    try:
        filename = utils.combine_dict_values(data)
        logger.debug(f"File name : {filename}")
        
        meta_data = db.get_metadata(filename)        
        if not utils.keys_exists(meta_data, "expiries"):
            logger.info("Fetching expiries from API")
            if not meta_data:
                meta_data = data
                
            ticker = Ticker(data['symbol'])
            expires = ticker.quotes(instrument=data['instrument'])
            
            logger.debug(f"Storing Metadata for expiries")
            
            meta_data.update({'expiries' : expires['expiries']})
            db.set_metadata(meta_data, filename)
            
            logger.info(f"Exiting fetch_strike_prices")
            return expires['expiries']
        else:
            logger.info("Fetching expiries from database")
            return meta_data['expiries']
        
    except Exception as ex:
        logger.info(ex)
        
    logger.error("Unable to fetch data for expiries...")
    return []

def fetch_strike_prices(data):
    logger.info("Inside fetch_strike_prices")
    try:
        filename = utils.combine_dict_values(data)        
        logger.debug(f"File name : {filename}")
        
        meta_data = db.get_metadata(filename)
        
        if not utils.keys_exists(meta_data, "strikePrices"):
            logger.info("Fetching strikePrices from API")
            
            if not meta_data:
                meta_data = data
            
            ticker = Ticker(data['symbol'])
            sp = ticker.quotes(instrument=data['instrument'], 
                                 expiry=data['expiry'], 
                                 optionType=data['optionType']) 
            
            logger.debug(f"Storing Metadata for strike prices")
            
            meta_data.update({'strikePrices' : sp['strikePrices']})
            db.set_metadata(meta_data, filename)
            
            logger.info(f"Exiting fetch_strike_prices")
            return data['strikePrices']
        else:
            logger.info("Fetching strikePrices from database")
            return meta_data['strikePrices']
        
    except Exception as ex:
        logger.info(ex)
        
    logger.error("Unable to fetch data for strike prices...")
    return []