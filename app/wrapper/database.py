# -*- coding: utf-8 -*-
"""
Created on Sun Mar 29 18:32:11 2020

@author: vsheoran
"""

import pandas as pd
from datetime import datetime
from csv import DictWriter
from .utilities import get_logger, real_time_field_names
import glob, os 

logger = get_logger("db.py")
    
class DB:
    base_path = os. getcwd()
    path_to_symbols = base_path + os.sep + "data" + os.sep + "Symbol.xlsx"
    base_real_time_path = base_path + os.sep + "data" + os.sep + "daily" + os.sep
    # "E:/Project/trends/src/data/daily/"
    base_historical_path = base_path + os.sep + "data" + os.sep + "historical" + os.sep
    # "E:/Project/trends/src/data/historical/"
    path_to_real_time_csv = str()
    path_to_historical_csv = str()
    # path_to_real_time_pkl = str()
    pattern_to_historical_data = str()
    real_time_field_names = list()
    cur_index = str()
    max_length = 10990
    
    def __init__(self, listing = None, date = None):
        global real_time_field_names
        self.real_time_field_names = real_time_field_names
        if listing is not None:
            self.init_config(listing, date)
        
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
                
        # Write to pickle
        # with open (self.path_to_real_time_pkl, 'wb') as dump:
        #     db = list()
        #     pickle.dump(db, dump)
            
                

    def __set_path_for_csv(self):
        base_path_real = self.base_real_time_path + self.cur_index + os.sep
        base_path_historical = self.base_historical_path + self.cur_index + os.sep
        
        self.path_to_real_time_csv =  base_path_real + self.__get_formatted_date(datetime.today()) + ".csv"    
        # self.path_to_real_time_pkl = base_path_real + self.__get_formatted_date(datetime.today()) + ".pkl" 
        
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
        symbols = pd.read_excel(self.path_to_symbols)
        return symbols.to_dict('records')
    
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
                
                return resp.values.tolist()
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
    
    def get_historical_data(self):  
        if os.path.isfile(self.path_to_historical_csv):
            csv = pd.read_csv(self.path_to_historical_csv)
            
            if len(csv) > self.max_length:
                csv = csv.tail(self.max_length)
            return csv.to_dict('records')
        else:
            return []
        
    def set_historical_data(self, df):
        li = glob.glob(self.pattern_to_historical_data)
        for item in li:
            logger.info(f"Removing file : {item}")
            os.remove(item)
        df.to_csv(self.path_to_historical_csv, index=False)  
        
    def reset(self):
        if os.path.isfile(self.path_to_historical_csv):
            os.remove(self.path_to_historical_csv)
            
        if os.path.isfile(self.path_to_real_time_csv):
            os.remove(self.path_to_real_time_csv)
    
    # def get_real_time_pkl(self): 
    #     db = dict()
    #     if os.path.isfile(self.path_to_real_time_pkl):
    #         with open(self.path_to_real_time_pkl , 'rb') as dbfile:
    #             db = pickle.load(dbfile)
    #     return db
    
    # def set_real_time_pkl(self, data):
    #     db = list()
    #     with open(self.path_to_real_time_pkl , 'rb') as dbfile:
    #         db = pickle.load(dbfile)
        
    #     db.extend(data)
        
    #     with open (self.path_to_real_time_pkl, 'wb') as dump:
    #         pickle.dump(db, dump)
            
    #     return db
