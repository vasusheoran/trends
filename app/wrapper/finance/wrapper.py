# -*- coding: utf-8 -*-

from . import Ticker
from .yahoo import historical_data
from logging import handlers
import logging, os

def derivative(ticker=None, period="3months", expiry=None, instrument = "FUTCUR", include = ["Date", "Close", "High", "Low"]):
    try:
        logger.info("Inside derivative - finance:wrapper")
        tk = Ticker(ticker)
        
        df = tk.future_cur(period=period, expiry=expiry, instrument=instrument, include=include)
        return reset(df)
    except Exception as ex:
        logger.info(ex)        

def reset(df):
    df = df.rename(columns={"High": "HP", "Close": "CP", "Low":"LP"})
    return df

def history(ticker, period = "2y", interval = "1d"):
    df = historical_data(listing = ticker, period = period, interval = interval)
    return df


base_path = "data" + os.sep

def get_logger(filename):
    global base_path
    # log_format = "[%(levelname)s]   | [%(name)s]    | %(asctime)s | %(filename)s    | Line : %(lineno)d | %(message)s"
    log_format = "[%(levelname)s] | %(asctime)s | %(thread)s | Line : %(lineno)d | [%(name)s] | %(message)s"
    
    logger = logging.getLogger(str(filename))
    logger.handlers = []
    
    # To override the default severity of logging
    logger.setLevel('DEBUG')
    
    # Use FileHandler() to log to a file
    # file_handler = logging.FileHandler(base_path + "app.log")
    file_handler = handlers.RotatingFileHandler(base_path + "app.log", 
                                                        mode='a', 
                                                        maxBytes=5*1024*1024,
                                                        backupCount=5, 
                                                        encoding=None, 
                                                        delay=0)
    formatter = logging.Formatter(log_format)
    file_handler.setFormatter(formatter)
    
    # Don't forget to add the file handler
    logger.addHandler(file_handler)
    
    return logger

logger = get_logger('fin-wrapper.py')