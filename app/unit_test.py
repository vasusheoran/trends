# -*- coding: utf-8 -*-
"""
Created on Tue Apr  7 14:28:48 2020

@author: vedvp
"""


import pandas as pd
import sys, os
from dateutil.relativedelta import relativedelta

# insert at 1, 0 is the script path (or '' in REPL)
sys_path = os.getcwd() + os.sep + 'modules'
if sys_path not in sys.path:
    sys.path.append(sys_path)
else:
    print('Modules path exists ..')
sys_path = os.getcwd() + os.sep + 'wrapper'
if sys_path not in sys.path:
    sys.path.append(sys_path)
else:
    print('Wrapper path exists ..')
    
from calculate import Calculate
from datetime import datetime, timedelta
from yfinance_wrapper import get_historical_data_by_time

def set_up(index):
    
    test_data = pd.read_excel("E:/test-daily.xlsx")
    start_date = test_data.at[0,'Date'] - relativedelta(years=3)
    date = test_data.at[len(test_data) - 1, 'Date'] +  timedelta(days=1)
    end_date = date
    
    # Fetch Yahoo data
    data = get_historical_data_by_time(index , "1d", start=start_date, end=end_date)
    
    return data, test_data

def test_data(test_row, max_length):
    # global test, data
    # Extract data from global finance data
    # test_row = test.iloc[index].to_dict()
    test_date = test_row['Date']
    test_row_index = data.index[data['Date'] == test_date].to_list()
    
    if bool(test_row_index):
        
        if test_row_index == len(data) + 1:
            test_sample_data = data
        else:
            test_sample_data = data[:test_row_index[0]]
        
        # Create Index
        calc = Calculate()
        calc.process_file_or_df(test_sample_data,"nse")
        freeze_value = calc.freeze_value({'CP' : test_row['Freeze_CP'], 
                           'HP' : test_row['Freeze_HP'], 
                           'LP' : test_row['Freeze_LP'], 
                           'Date' : test_row['Date']})
        calc.update(test_row)
        
        values = calc.fetch_values()
        
        freeze_diff = test_row['Freeze_BI'] - freeze_value['bi']
        
        bi_diff = test_row['BI'] - values['dashboard']['cards'][0]['value']
        
        # print(test_date , values['bi'])
        
        return  {'bi' : bi_diff, 'freeze' : freeze_diff, 'time' : test_date}
    else:
        print("Not possible")
        return (test_date, "")
    
def print_formatter(decimal):
    print("BI differs by {:.4f}".format(decimal))
    

data, test = set_up("^nsei")
test = test.dropna().reset_index(drop=True)
test_row = test.iloc[10]
test_row_len = len(test)
max_length = len(data)

li = []
for index, row in test.iterrows():
    # test_row = row.to_dict()
    temp = test_data(row, max_length)
    li.append(temp)
    print(f"{temp['freeze']} --  {temp['bi']}")
    
print(li)
    # break
# index = "^nsei"
# print(li)
    
