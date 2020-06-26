# -*- coding: utf-8 -*-
"""
Created on Mon Mar 23 19:09:43 2020

@author: vsheoran
"""


#from lib.calculate import Calculate
import pandas as pd
from .utilities import get_logger, Utilities
from .calculate import Calculate
from .finance import *
from .database import DB
from flask import abort, Response, jsonify
from datetime import datetime
from flask_socketio import emit

calc = None
df = None
current = dict()
db = DB()
utils = Utilities()
count = 0

socketio = None

logger = get_logger("wrapper.py")

def set_up_socketio(io):
    global socketio
    socketio = io
    
def fetch_symbol_list():
    global db
    response = db.get_listings()
    return response

def reset_current_index():
    try:
            
        global calc, df, current, db
        logger.info("Inside reset_current_index")
        
        calc = None
        df = None
        
        li =[]
        
        util = Utilities()
        if util.keys_exists(current, 'listing', 'options'):
            options = current['listing']['options']
            
            strike_filename = util.combine_dict_values(options)
            
            expiry_filename = util.combine_dict_values({'symbol' : options['symbol'],
                                                        'instrument' : options['instrument']})
            li.append(expiry_filename)
            li.append(strike_filename)
            current = dict()
        db.reset(li) 
        
        db = DB()
        logger.info("Exiting reset_current_index")
        return {"data": "Reset Successfull"}
    except Exception as err:
        logger.error(err)
        response = jsonify({'message': repr(err)})
        response.status_code = 400
        return response
        
    # return db
    
def set_current_listing(jsonData):
    logger.info("Inside set_current_listing")
    global calc, df, db, current
    
    calc = None
    df = None
    current = dict()
    
    current.update({'listing' : jsonData})
    
    data = pd.DataFrame()
    
    yahoo_index = jsonData['Symbol']
    period = "2y"
    interval = "1d"
        
    db_index = jsonData['SAS']  
        
    logger.info(f"DB Index File Name : {db_index}")
    db = DB(listing=db_index)
    historical_data = db.get_historical_data()
        
    try:
        if not historical_data:
            logger.info("Fetching data from nse .. ")
            data = pd.DataFrame()
            if 'Derivative' in jsonData['Series']:
                if utils.keys_exists(jsonData, "options", "strikePrice"):
                    strike=jsonData['options']['strikePrice']
                    
                else:
                    strike=None
                if utils.keys_exists(jsonData, "options", "option"):
                    option=jsonData['options']['option']
                    
                else:
                    option=None
                    
                data = derivative(ticker=jsonData['Symbol'], 
                                  period=jsonData['options']['period'], 
                                  expiry = jsonData['options']['expiry'], 
                                  instrument = jsonData['options']['instrument'],
                                  option=option,
                                  strike=strike)
                
            else:
                data = history(yahoo_index , period, interval)
                
            db.set_historical_data(data)
        else:
            logger.info("Fetching data from yahoo finance .. ")
            data = pd.DataFrame.from_records(historical_data)
            
        logger.info(f"Processing data...")
        calc = Calculate()
        calc.process_file_or_df(data, db_index)
        
        logger.info(f"Exiting set_current_listing")
        
        if historical_data and 'Derivative' in jsonData['Series']:
            return {'status' : "success", "msg" : "Symbol already exists. Please reset the symbol to change options."}
        else:
            return {'status' : "success", "msg" : "Symbol set successfully."}
        
    except KeyError as err:
        logger.error(err)
        # abort(500, description=f"Key Error. {err}")
            
        response = jsonify({'message': repr(err)})
        response.status_code = 400
        return response
    except Exception as err:
        logger.error(err)
        # abort(500, description=err)
            
        response = jsonify({'message': repr(err)})
        response.status_code = 400
        return response

def fetch_index_if_set():
    global current, db, calc
    try:
        real_time_data = db.get_real_time_data()
        df = calc.get_dataframe()
        response = {'chart' : {'update' : [(datetime.today().timestamp()) , df.at[2, 'CP'] ], 
                               'listing' : current['listing'], 
                               'data' : real_time_data}}
        # Card values for the latest point
        if len(real_time_data) > 0:
            last_row = db.get_latest_record()
            values = calc.update(last_row)
            response.update({'data' : values})
        else:
            historical_data = paginate(1, 1)[0]
            historical_data.update({'index' : current['listing']['SAS']})
            values = calc.update(historical_data)
            response.update({'data' : values})
            
        return response
    except KeyError as er:
        logger.error(er)
        abort(500, description="No data found for this date range. Key Error.")
    except Exception as er:
        logger.error(er)
        abort(500, description="Please choose a index from search option.")
    
def fetch_updated_or_frozen(isUpdateEnabled = True):
    global calc
    try: 
        response = None
        if isUpdateEnabled:
            response = {'data' : calc.fetch_values()}
        else:
            response = {'data' : calc.fetch_frozen_values()}
        return response
    except Exception as ex:
        logger.error(ex)
        abort(500, description="The server encountered an internal error and was unable to complete your request. Try setting the listing again.")

def update_values(ob, isFreezeEnabled = False):
    logger.info("Inside set_current_listing")
    global calc, current
    
    try:
        if isFreezeEnabled and ob['index'] == current['listing']['SAS']:
            calc.freeze_value(ob)   
        else:    
            values = calc.update(ob)
            
            push_notifications('updateui', values)
    except Exception as err:
        logger.error(err)
    logger.info("Exiting set_current_listing")
    
def push_notifications(eventName, update):
    logger.info("Inside push_notifications")
    socketio.emit(eventName, update)
    logger.info("Exiting push_notifications")
    
def paginate(page, size):
    offset = 3
    global calc
    if calc == None:
        return []
    df = calc.get_dataframe()
        
    if size > 0 and page > 0:
        if page > 1:
            start_index = ((page - 1) * size) + 1 + offset
        else:
            start_index = 0 + offset
        end_index = ((page - 1) * size) + size + offset
        resp = df.iloc[start_index : end_index]
        resp =resp[['CP', 'HP', 'LP', 'Date']]
        resp['Date'] = pd.to_datetime(resp['Date'])
        return resp.to_dict('records')
    return []

def fetch_data_by_start_end(start, end):
    global db
    data = db.get_real_time_data(start, end)    
    return {'start' : start, 'end' : end, 'data' : data}

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
            
def fetch_current():
    return current