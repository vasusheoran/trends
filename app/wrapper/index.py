# -*- coding: utf-8 -*-

from .utilities import get_logger, Utilities
from .calculate import Calculate
from .database import DB
from .finance import *
from flask import jsonify, abort
import pandas as pd


logger = get_logger("index.py")

class Index:    
    def __init__(self):
        # {
            # calc : __,
            # values : __,
        # }
        self.index = None
        self.index_map = dict()
        
        self.db = DB()
        self.utils = Utilities()
        
    def get(self):
        try: 
            if self.__is_index_present():
                return {'values' : self.index_map[self.index]['values'],
                        'symbol' : self.index_map[self.index]['symbol'],
                        'data' : self.db.get_real_time_data()}
            else:
                abort(500, description="The server encountered an internal error and was unable to complete your request. Try setting the listing again.")

        except Exception as ex:
            logger.error(ex)
            abort(500, description="The server encountered an internal error and was unable to complete your request. Try setting the listing again.")

    def put(self, stock_data):
        if self.__is_index_present():
            calc = self.index_map[self.index]['calc']
            calc.update(stock_data)            
            self.index_map[self.index]['values'] = calc.fetch_values()

            # Updating latest buy/sell value in freeze object for open_dailogue
            self.index_map[self.index]['freeze']['Buy'] = self.index_map[self.index]['values']['dashboard']['cards'][0]['value']
            self.index_map[self.index]['freeze']['Sell'] = self.index_map[self.index]['values']['dashboard']['cards'][2]['value']
            return self.index_map[self.index]['values'], True
        else:
            return None, False

    def post(self, symbol):
        
        logger.info("Inside post")
        
        index = symbol['SAS'] 
        self.db = DB(listing=index)
        historical_data = self.db.get_historical_data_list()
        
        try:
            if not historical_data:
                logger.info("Fetching data from nse .. ")
                # data = pd.DataFrame()
                if 'Derivative' in symbol['Series']:
                    data = self.__derivative_data(symbol)
                else:
                    data = history(symbol['Symbol'] , "2y", "1d")
                    
                self.db.set_historical_data(data)
            else:
                logger.info("Fetching data from yahoo finance .. ")
                data = pd.DataFrame.from_records(historical_data)
                
            calc = Calculate()
            calc.process_file_or_df(data, index)
             
            # Update self only if index mapping not present
            self.index = symbol['SAS']
            if self.__is_index_present() == False:  
                logger.info("New Index -- updating records")
                self.index_map[self.index] = {
                    'calc' : calc,
                    'values' : calc.fetch_values(),
                    'freeze' : calc.fetch_frozen_values(),
                    'symbol' : symbol
                }
                
            
            if historical_data and 'Derivative' in symbol['Series']:
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
        
    def delete(self):
        if self.index == None:
            return {'message' : "Please choose a symbol to delete the data."}
            
        try:
            logger.info("Inside delete")
            
            li =[]
            
            if 'options' in self.index_map[self.index]:
                options = self.index_map[self.index]['options']
                
                strike_filename = util.combine_dict_values(options)
                
                expiry_filename = util.combine_dict_values({'symbol' : options['symbol'],
                                                            'instrument' : options['instrument']})
                li.append(expiry_filename)
                li.append(strike_filename)
                current = dict()
            self.db.reset(li) 
            
            # Reseting current index
            del self.index_map[self.index]
            
            return {"data": "Reset Successfull"}
        except Exception as err:
            logger.error(err)
            response = jsonify({'message': repr(err)})
            response.status_code = 400
            return response

    def refresh(self):
        symbol = self.index_map[self.index]['symbol']
        
        logger.info("Removing hashed data")
        del self.index_map[self.index]


        logger.info("Calculate new index ..")
        self.post(symbol)
        
        
    def post_freeze(self, stock_data):           
        if self.__is_index_present() == False:
            abort(500, description="The server encountered an internal error and was unable to complete your request. Try setting the listing again.")
 
        logger.info("Inside update_freeze")
        if stock_data['Date'] == None:
            stock_data['Date'] = datetime.today().strftime("%m:%d:%Y %H:%M:%S")
            
        # Fetching calculate object 
        calc = self.index_map[self.index]['calc']
        
        # Freeze values
        calc.freeze_value(stock_data)
        
        # Update Index
        # TODO : Check if updating calc object is required
        # self.index_map[self.index]['calc'] = calc
        self.index_map[self.index]['values'] = calc.fetch_values()
        self.index_map[self.index]['freeze'] = calc.fetch_frozen_values()  
        
        
    def get_freeze(self):
        if self.__is_index_present():
            return self.index_map[self.index]['freeze']
        else:
            abort(500, description="The server encountered an internal error and was unable to complete your request. Try setting the listing again.")
 
        
    def name(self):
        return self.index

    def symbol(self):
        return self.index_map[self.index]['symbol']
    
    def get_df(self):
        if self.__is_index_present():
            return self.index_map[self.index]['calc'].get_dataframe()
        else:
            abort(500, description="The server encountered an internal error and was unable to complete your request. Try setting the listing again.")
    
    def __is_index_present(self):
        if self.index in self.index_map:
            return True
        else:
            return False
        
    def __update_index(self, calc): 
        self.index_map = {
            self.index : {
                'calc' : calc,
                'values' : calc.fetch_values()
                }
            }
        
        
    def __derivative_data(self, symbol):
        if self.utils.keys_exists(symbol, "options", "strikePrice"):
            strike=symbol['options']['strikePrice']
            
        else:
            strike=None
        if self.utils.keys_exists(symbol, "options", "option"):
            option=symbol['options']['option']
            
        else:
            option=None
            
        data = derivative(ticker=symbol['Symbol'], 
                          period=symbol['options']['period'], 
                          expiry = symbol['options']['expiry'], 
                          instrument = symbol['options']['instrument'],
                          option=option,
                          strike=strike)
        return data

index = Index()