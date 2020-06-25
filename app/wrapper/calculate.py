# -*- coding: utf-8 -*-
"""
Created on Wed Mar 25 20:36:00 2020

@author: vsheoran
"""

#from modules.utilities import Utilities
import pandas as pd
from .utilities import Utilities, get_logger
from .calchelper import find_BI
from datetime import datetime
# import alpha_vantage_wrapper

logger = get_logger("calcuate.py")

class Calculate:
    df = pd.DataFrame()
    curIndex = None
    curValues = dict()
    db = dict()
    frozen = dict()
    def __init__(self, df = None):
        self.bk = float()
        self.util = Utilities()
        self.curRow = 2
        self.cols= ['CP', 'CP_CI_HP']
        
        self.db = dict()
        self.frozen = dict()
        self.back_ground = dict()
        
        if isinstance(df, pd.DataFrame) and df.empty == False:
            self.process_data(df)
            self.set_up() 
            
    def process_file_or_df(self,val, index):    
        dataset = pd.DataFrame()
        if isinstance(val , str):
            dataset = pd.read_excel(val)
        elif isinstance(val, pd.DataFrame):
                dataset = val
        dataset = dataset.dropna(how='any',axis=0) 
        self.__process_data(dataset)
        self.__set_up()
        self.curValues.update({'index' : index, 'Date' : (datetime.today().timestamp() )} )

    def process_new_values(self, df):
        if isinstance(df, pd.DataFrame):
            logger.info("Processing ..")
            self.__process_data(df)
            self.__set_up_new_only()
            

    def __find_ema(self, spans, cols ,df):
        for span in spans:
            df = self.util.ema_rolling(df, cols, span)
        return df  
    
    def __set_up(self):
        self.df = self.__find_ema([5,20, 50], self.cols, self.df)
        self.df = self.util.findMin(self.df, ['HP'],3)
        self.df = self.util.av_rolling(self.df,['CP_CI_HP'],50)
        self.df = self.util.av_rolling(self.df,['CP_CI_HP'],10)  
        
    def __set_up_new_only(self, num = 2):
        spans = [5, 20, 50]
        avColName = 'CP_CI_HP'
        for i in range(num, -1 , -1):
            # Set up new EMI Values
            for col in self.cols:
                for span in spans:
                    col_name = 'ema' + str(col) + str(span)
                    prev_ema = self.df.at[i + 1, col_name]
                    new_ema = ((2/(span+1)) * (self.df.at[i, col] - prev_ema)) + prev_ema
                    self.df.at[i, col_name] = new_ema
                    
            # Set up new minHP3 Values
            val = min(self.df.at[i, 'HP'], self.df.at[i + 1, 'HP'], self.df.at[2, 'HP'])
            self.df.at[i, 'minHP3'] = val
            for span in [10, 50]:
                col_name = 'av' + str(avColName) + str(span)
                sma = self.df[:span + 3]
                self.df.loc[:3,col_name] = sma[avColName].rolling(window=span).mean()[-4::].reset_index(drop = True)
              
        
    def get_dataframe(self):
        return self.df
        
    def set_dataframe(self, df, index):
        self.curValues.update({'index' : index})
        self.df = df
        
    def __process_data(self, df):
        df = df.reindex(index=df.index[::-1]).reset_index(drop=True)
        
        df['CP_CI_HP'] = df['CP']
        
        last_day = df.iloc[0].to_dict()
        currentDayRow = pd.DataFrame(
                {'CP': last_day['CP'], 'HP':last_day['HP'], 
                 'LP': last_day['LP'], 'CP_CI_HP' : last_day['HP']}, index =[0]) 
        
        df = pd.concat([currentDayRow, df]).reset_index(drop = True)
        
        
#        Using dependents
            
        HP = df.iloc[:5]
        HP = self.util.findMin(HP, ['HP'], 3)
        
        last_day = df.iloc[0].to_dict()
        last_day['minHP3'] = HP.at[0,'minHP3']
#        logger.info(last_day)
        
        nextDayRow = {'CP': last_day['minHP3'], 'HP':last_day['minHP3'],
                      'LP': last_day['minHP3'], 'CP_CI_HP' : last_day['HP']}
        
        new_df = pd.DataFrame([nextDayRow, nextDayRow]) 
        new_df = self.__update_dependents(last_day, new_df)
        
        
        
        
        self.df = pd.concat([new_df, df[:]], sort=True).reset_index(drop = True)
        
        self.last_day = last_day
        self.curValues.update({'CP' : self.df.at[2,'CP'], 'HP' : self.df.at[2,'HP'],'LP' : self.df.at[2,'LP']})
        
    def __update_dependents(self, last_day, df):
        
        df.at[1,'CP'] = last_day['CP']
        df.at[1,'HP'] = last_day['HP']
        df.at[1,'LP'] = last_day['CP']
        df.at[1,'CP_CI_HP'] = last_day['HP']
        
        df.at[0,'CP'] = last_day['minHP3']
        df.at[0,'HP'] = last_day['minHP3']
        df.at[0,'LP'] = last_day['minHP3']
        df.at[0,'CP_CI_HP'] = last_day['HP']
        
        return df
        
        
    def __update_cp(self, new_cp = None, new_hp = None, new_lp = None):
#        Update CP, LP, HP
        
        if new_cp < new_lp and new_cp > new_hp:
            raise RuntimeError("Fatal Error. Please check values for cp/hp/lp.")
        
        self.df.at[2,'CP'] = new_cp  
        self.df.at[2,'LP'] = new_lp   
        self.df.at[2,'HP'] = new_hp
        
        # TODO: Verification
        
        # self.df.at[2, 'CP_CI_HP'] = new_hp        
            
        last_day = self.df.iloc[2].to_dict()
            
        HP = self.df.iloc[:5]
        HP = self.util.findMin(HP, ['HP'], 3)
        
        last_day['minHP3'] = HP.at[2,'minHP3']
        
        self.df = self.__update_dependents(last_day, self.df)
        self.__set_up_new_only()
    
    def update(self, val):        
        self.__update_cp(val['CP'], val['HP'], val['LP'])
        self.curValues = val.copy()  
        return self.fetch_values()  
    
    def freeze_value(self, val):
        temp = {'CP' : self.df.at[2,'CP'], 'HP' : self.df.at[2,'HP'],'LP' : self.df.at[2,'LP']}
        self.__update_cp(val['CP'], val['HP'], val['LP'])
        res = find_BI(self.db, self.frozen, self.df, True)
        # self.__update_cp(temp['CP'], temp['HP'], temp['LP'])
        
        self.frozen.update({
            'CP': val['CP'] ,
            'HP': val['HP'],
            'LP': val['LP'],
            'Date': val['Date'],
            'bi': res
            })
        return self.frozen
              

    def fetch_rows(self, endIndex, startIndex = 0):
        temp = self.df[startIndex:endIndex+startIndex].copy()
        return temp[['CP','HP', 'LP','Date']]
        
    def fetch_values(self, isUpdated=False):
        if isUpdated:
            return self.db
        else:
            
            if self.df.at[2,'emaCP5'] == None:
                ema5 = ""
            else:
                ema5 = self.df.at[2,'emaCP5']
                
            if self.df.at[2,'emaCP20'] == None:
                ema20 = ""
            else:
                ema20 = self.df.at[2,'emaCP20']
                
            if self.df.at[2,'emaCP50'] == None:
                ema50 = ""
            else:
                ema50 = self.df.at[2,'emaCP50']
                                
            response = find_BI(self.db, self.frozen, self.df)
            self.back_ground = response
            self.db.update({
                'dashboard' :{
                    'cards' : [{'name' : 'Buy', 'key' : 'Buy', 'value': response['bi']},
                               {'name' : 'Support', 'key' : 'Support', 'value': response['bj']},
                               {'name' : 'Sell', 'key' : 'Sell', 'value': response['bk']},
                               {'name' : 'Min High', 'key' : 'Min_High', 'value': response['min.HP.2']
                                }],
                    'table' : [{'name' : 'Close', 'value': self.curValues['CP']},
                               {'name' : 'High', 'value': self.curValues['HP']},
                               {'name' : 'Low', 'value': self.curValues['LP']},
                               {'name' : 'EMA 5', 'value': ema5},
                               {'name' : 'EMA 20', 'value': ema20},
                               {'name' : 'EMA 50', 'value': ema50},
                               # {'name' : 'RSI', 'value': 0},
                               # {'name' : 'PE', 'value': 0}
                               ]},
                'Date' : {'name' : 'Date', 'value': self.curValues['Date']},
                'stocks' : [self.curValues['Date'], self.curValues['CP']]})
        return self.db
        
    def fetch_frozen_values(self):
        # Update values if not yet frozen
        if 'bi' not in self.frozen:
            bi = find_BI(self.db, self.frozen, self.df, True)
            self.frozen.update({'bi' : bi})          
        return self.frozen
    
    def fetch_back_values(self):
        return self.back_ground