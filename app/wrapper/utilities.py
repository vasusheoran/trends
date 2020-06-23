# -*- coding: utf-8 -*-
"""
Created on Mon Mar 23 19:16:30 2020

@author: vsheoran
"""

import logging, os
from datetime import datetime, timedelta
from logging import handlers

base_path = "data" + os.sep
real_time_field_names = ["CP","HP","LP","Date","index"]
max_freeze_time =  datetime.now().replace(hour=9,minute=19,second=58)
min_freeze_time =  max_freeze_time + timedelta(seconds=1)

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

class Utilities:
#    def __init__(self):
#        print("Initializing utilites")
    
    def minOrMaxInWindow(self, arr, n, k, isMin = True, includeSelf = False): 
        include = 0
        if includeSelf:
            include = 1
        
        li = [0 for i in range(k - include)]
        min = 0
        
        for i in range(k - include, n):
            min = arr[i-1 + include] 
            for j in range(2 - include, k+1 - include): 
                if isMin:
                    if arr[i - j] < min: 
                        min = arr[i - j]
                else:
                    if arr[i - j] > min: 
                        min = arr[i - j] 
            
            li.append(min)
         
        
        return li
    
    def findMin(self, df, cols, span, reverse = True):
        if reverse:
            df = df.iloc[::-1]
        for col in cols:
            col_name = 'min' + str(col) + str(span)
#            if col_name not in df:
            df[col_name] = df[col].rolling(window=span,min_periods=1).apply(min, raw=True).dropna()
            
        if reverse:
            return df.iloc[::-1]
        else:
            return df
        
    def findMax(self, df, cols, span, reverse = True):
        if reverse:
            df = df.iloc[::-1]     
        for col in cols:
            col_name = 'max' + str(col) + str(span)
#            if col_name not in df:
            df[col_name] = df[col].rolling(window=span,min_periods=1).apply(max, raw=True).dropna()
            
        if reverse:
            return df.iloc[::-1]
        else:
            return df
    
    def ema_rolling(self, df, cols, span, reverse = True):
        if reverse:
            df = df.iloc[::-1]
        for col in cols:
            col_name = 'ema' + str(col) + str(span)
            sma = df[col].rolling(span).mean()
            modPrice = df[col].copy()
            modPrice.iloc[0:span] = sma[0:span]
            df[col_name] = modPrice.ewm(span=span, adjust=False).mean()
            
        if reverse:
            return df.iloc[::-1]
        else:
            return df
        
    def av_rolling(self, df, cols, span, reverse=True):
        if reverse:
            df = df.iloc[::-1]
        for col in cols:
            col_name = 'av' + str(col) + str(span)
            sma = df[col].rolling(window=span).mean()
#            sma = sma.shift(periods=1, fill_value=0)
            df[col_name] = sma
        if reverse:
            return df.iloc[::-1].copy()
        else:
            return df.copy()
        
        
    def ema_update(self, df, row, col, span):
        vals = row        
        for i in range(vals+1):
            
            val = df.at[row, col]
            ema_col_name = 'ema' + str(col) + str(span)
#            print(ema_col_name + " -- " + str(row+1))
            prev_row_ema = df.at[row + 1, ema_col_name]
            new_ema = self.ema_calculate(span, val, prev_row_ema)
            df.at[row, ema_col_name] =  new_ema
            row = row-1
            
        return df
    
    def ema_calculate(self, span, val, prev_ema):
        return ((2/(span+1)) * (val - prev_ema)) + prev_ema

    def combine_dict_values(self, elem):
        return '_'.join(y for x, y in sorted(elem.items()))
    
    def keys_exists(self, element, *keys):
        '''
        Check if *keys (nested) exists in `element` (dict).
        '''
        if not isinstance(element, dict):
            raise AttributeError('keys_exists() expects dict as first argument.')
        if len(keys) == 0:
            raise AttributeError('keys_exists() expects at least two arguments, one given.')
    
        _element = element
        for key in keys:
            try:
                _element = _element[key]
            except KeyError:
                return False
        return True
    