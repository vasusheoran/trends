# -*- coding: utf-8 -*-
"""
Created on Mon Mar 23 19:09:43 2020

@author: vsheoran
"""


#from lib.calculate import Calculate
import inspect,os
import pandas as pd
from calculate import Calculate
from yfinance_wrapper import get_historical_data_by_period
from async_es_task import AsyncUpdateHistoricalTask
from DB import DB
import numpy as np
from flask import Flask, abort, Response
from flask_socketio import SocketIO, emit
from flask_cors import CORS
from utilities import get_logger
from datetime import datetime

calc = None
df = None
current = dict()
db = DB()
count = 0

app = Flask(__name__) 
app.config['UPLOAD_FOLDER'] = os.getcwd() + os.sep + 'files'
app.config['SECRET_KEY'] = 'secret!'
CORS(app, resources={r"/*": {"origins": "*"}})
socketio = SocketIO(app, cors_allowed_origins="*")

logger = get_logger("wrapper.py")

    
def fetch_symbol_list():
    global db
    response = db.get_listings()
    return response

def reset_current_index():
    global calc, df, current
    calc = None
    df = None
    current = dict()
    
    return {"data": "Reset Successfull"}
    
def set_current_listing(jsonData):
    global calc, df, global_data, current_listing, db
    current.update({'listing' : jsonData})
    
    data = pd.DataFrame()
    
    yahoo_index = jsonData['YahooSymbol']
    db_index = jsonData['SASSymbol']  
    period = "2y"
    interval = "1d"
    db = DB(listing=db_index)
    historical_data = db.get_historical_data()
    
    try:
        if not historical_data:
            logger.info("Fetching data from yahoo finance .. ")
            data = pd.DataFrame()
            data = get_historical_data_by_period(yahoo_index , period, interval)
            db.set_historical_data(data)
        else:
            data = pd.DataFrame.from_records(historical_data)
            
        calc = Calculate()
        calc.process_file_or_df(data, db_index)
        
        return fetch_index_if_set()
        
    except KeyError:
        abort(500, description="No data found for this date range, symbol may be delisted.")

def fetch_index_if_set():
    global current, db, calc
# =============================================================================
#     real_time_data = db.get_real_time_data()
#     df = calc.get_dataframe()
#     response = {'chart' : {'update' : [(datetime.today().timestamp()) , df.at[2, 'CP'] ], 
#                            'listing' : current['listing'], 
#                            'data' : real_time_data}}
#     if len(real_time_data) > 0:
#         last_row = db.get_latest_record()
#         values = calc.update(last_row)
#         # response.update({'data' : values})
#     return response
# =============================================================================
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
            historical_data.update({'index' : current['listing']['SASSymbol']})
            values = calc.update(historical_data)
            response.update({'data' : values})
            
        return response
    except KeyError as er:
        logger.error(er)
        abort(500, description="No data found for this date range. Key Error.")
    except Exception as er:
        logger.error(er)
        abort(500, description="Please choose a index from search option.")
    
def add_new_rows(ob):
    global calc, current
    resp = calc.add_new_days(ob)
    vals = dict()
    vals.update({'data' : calc.get_dataframe(), 'index' : current['listing']['YahooSymbol']})
    async_task = AsyncUpdateHistoricalTask(task_details=vals)
    async_task.start()
    return resp

def delete_new_rows():
    global calc, current
    resp = calc.delete_new_days()
    vals = dict()
    vals.update({'data' : calc.get_dataframe(), 'index' : current['listing']['YahooSymbol']})
    async_task = AsyncUpdateHistoricalTask(task_details=vals)
    async_task.start()
    return resp
    
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
    global calc, current
    
    try:
        if isFreezeEnabled and ob['index'] == current['listing']['SASSymbol']:
            calc.freeze_value(ob)   
        else:    
            values = calc.update(ob)
            
            push_notifications('updateui', values)
    except Exception as err:
        logger.error(err)
    #     return resp
    # return {}
    
def push_notifications(eventName, update):
    socketio.emit(eventName, update)
    
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