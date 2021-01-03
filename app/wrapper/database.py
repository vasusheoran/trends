# -*- coding: utf-8 -*-
"""
Created on Sun Mar 29 18:32:11 2020

@author: vsheoran
"""

import pandas as pd
from datetime import datetime
from csv import DictWriter
from .utilities import get_logger, real_time_field_names
import glob, os , pickle

logger = get_logger("db.py")
    
class DB:
    base_path = os. getcwd()
    path_to_symbols = base_path + os.sep + "data" + os.sep + "Symbols.csv"
    base_real_time_path = base_path + os.sep + "data" + os.sep + "daily" + os.sep
    base_metadata_path = base_path + os.sep + "data" + os.sep + "metadata" 
    # "E:/Project/trends/src/data/daily/"
    base_historical_path = base_path + os.sep + "data" + os.sep + "historical" + os.sep
    # "E:/Project/trends/src/data/historical/"
    path_to_real_time_csv = str()
    path_to_historical_csv = str()
    path_to_metadata = str()
    pattern_to_historical_data = str()
    real_time_field_names = list()
    cur_index = str()
    max_length = 10990
    
    def __init__(self, listing = None, date = None):
        global real_time_field_names
        self.real_time_field_names = real_time_field_names
        if listing is not None:
            self.init_config(listing, date)
            
        if not os.path.exists(self.base_metadata_path):
            os.makedirs(self.base_metadata_path)
        logger.debug(f"Real time data path : {self.path_to_real_time_csv}")
        logger.debug(f"Historical data path : {self.path_to_historical_csv}")
        
    def init_config(self, listing, date):        
        if not bool(date):
            date = datetime.today()
            
        self.cur_index = listing
        self.__set_path_for_csv()
        
        if not os.path.isfile(self.path_to_real_time_csv):
            # Write Headers
            with open (self.path_to_real_time_csv, 'a') as csvfile:
                ob_writer = DictWriter(csvfile, delimiter=',', 
                                        lineterminator='\n',fieldnames=self.real_time_field_names)
                ob_writer.writeheader()            
                

    def __set_path_for_csv(self):
        base_path_real = self.base_real_time_path + self.cur_index + os.sep
        base_path_historical = self.base_historical_path + self.cur_index + os.sep
        
        self.path_to_real_time_csv =  base_path_real + self.__get_formatted_date(datetime.today()) + ".csv"          
        self.path_to_historical_csv = base_path_historical + self.__get_formatted_date(datetime.today()) + ".csv" 
        self.pattern_to_historical_data = base_path_historical + self.__get_formatted_date(datetime.today()) + "*.csv" 
        
        if not os.path.exists(base_path_real):
            os.makedirs(base_path_real)
            
        if not os.path.exists(base_path_historical):
            os.makedirs(base_path_historical)
        
    def __get_formatted_date(self, date):
        return date.strftime("%m-%d-%Y")
    
    def get_path_to_real_time_data(self):
        return self.path_to_real_time_csv
        
    def get_listings(self):
        symbols = pd.read_csv(self.path_to_symbols)
        return symbols.to_dict('records')
        
    def set_listings(self, symbols):
        symbols.to_csv(self.path_to_symbols, index=False)
        
    def get_listings_df(self):
        return pd.read_csv(self.path_to_symbols)
    
    def get_latest_record(self):
        try:
            if os.path.isfile(self.path_to_real_time_csv):
                resp = pd.read_csv(self.path_to_real_time_csv).dropna()
                if not resp.empty:
                    return resp.tail(1).to_dict("records")[0]            
            return {}
        except Exception as err:
            logger.error("Unable to get latest data from server.")
            logger.error(err)
            return {}
            
    def get_real_time_data(self, start = None, end = None):    
        try:
            if os.path.isfile(self.path_to_real_time_csv):
                resp = pd.read_csv(self.path_to_real_time_csv).dropna()
                # TODO : Enable Sorting
                # resp  = resp.sort_values(by=["Date"], axis = 1)
                resp  = resp[["Date","CP"]]
         
                
                if start and end:
                    resp = resp[(resp['Date'] > start)
                                & (resp['Date'] < end)] 
                
                resp = resp.sort_values(by=["Date"])
                return resp.values.tolist()
                # return resp
            else:
                return []
        except Exception as err:
            logger.error("Unable to load real time date from server.")
            logger.error(err)
            return []
          
    def set_real_time_data(self, elems):
        try:
            with open(self.path_to_real_time_csv , 'a', newline='') as write_obj:
                dict_writer = DictWriter(write_obj, fieldnames=self.real_time_field_names)
                dict_writer.writerows(elems)
        except Exception as ex:
            raise ex
    
    def get_historical_data_list(self):
        if os.path.isfile(self.path_to_historical_csv):
            csv = pd.read_csv(self.path_to_historical_csv)
            
            if len(csv) > self.max_length:
                csv = csv.tail(self.max_length)
            return csv.to_dict('records')
        else:
            return []
    
    def get_historical_data_csv(self):
        if os.path.isfile(self.path_to_historical_csv):
            csv = pd.read_csv(self.path_to_historical_csv)
            
            if len(csv) > self.max_length:
                csv = csv.tail(self.max_length)
            return csv
        else:
            return pd.DataFrame()
        
    def set_historical_data(self, df):
        li = glob.glob(self.pattern_to_historical_data)
        for item in li:
            logger.info(f"Removing file : {item}")
            os.remove(item)
        df.to_csv(self.path_to_historical_csv, index=False)  
        
    def reset(self, symbols):
        try:
            if os.path.isfile(self.path_to_historical_csv):
                os.remove(self.path_to_historical_csv)
            logger.debug(f"Deleting : {self.path_to_historical_csv}")
            
            # if os.path.isfile(self.path_to_real_time_csv):
            #     os.remove(self.path_to_real_time_csv)
            # logger.debug(f"Deleting : {self.path_to_real_time_csv}")
                
            for symbol in symbols:
                path_to_metadata = self.base_metadata_path + os.sep + symbol
                if os.path.isfile(path_to_metadata):
                    os.remove(path_to_metadata)
                logger.debug(f"Deleting : {symbol}")
                
        except Exception as ex:
            raise RuntimeError(ex)
            
    def get_metadata(self, symbol):
        # filename = f"{symbol}_{self.__get_formatted_date(datetime.today())}.pkl"
        path_to_metadata = self.base_metadata_path + os.sep + symbol
        
        if not os.path.exists(path_to_metadata):
            with open(path_to_metadata , 'wb') as dbfile:
                data = dict()
                pickle.dump(data, dbfile)
        
        with open(path_to_metadata , 'rb') as dbfile:
            return pickle.load(dbfile)
    
    def set_metadata(self, data, symbol):
        # filename = f"{data['symbol']}_{self.__get_formatted_date(datetime.today())}.pkl"
        path_to_metadata = self.base_metadata_path + os.sep + symbol
        with open(path_to_metadata , 'wb') as dbfile:
            pickle.dump(data, dbfile)

sym = None