# -*- coding: utf-8 -*-
"""
Created on Mon Mar 30 21:58:19 2020

@author: vsheoran
"""

# import our libraries

import pythoncom
import win32com.client
import logging
from logging import handlers
from async_task import AsyncUpdateTask


def get_logger(filename):
    base_path = "C:/Users/vedvp/sasonline/app/data/"
    # log_format = "[%(levelname)s]   | [%(name)s]    | %(asctime)s | %(filename)s    | Line : %(lineno)d | %(message)s"
    log_format = "[%(levelname)s] | %(asctime)s | %(thread)s | Line : %(lineno)d | [%(name)s] | %(message)s"
    
    logger = logging.getLogger(str(filename))
    logger.handlers = []
    
    # To override the default severity of logging
    logger.setLevel('DEBUG')
    
    # Use FileHandler() to log to a file
    # file_handler = logging.FileHandler(base_path + "app.log")
    file_handler = handlers.RotatingFileHandler(base_path + "excel.log", 
                                                        mode='a', 
                                                        maxBytes=1*1024*1024,
                                                        backupCount=5, 
                                                        encoding=None, 
                                                        delay=0)
    formatter = logging.Formatter(log_format)
    file_handler.setFormatter(formatter)
    
    # Don't forget to add the file handler
    logger.addHandler(file_handler)
    
    return logger

logger = get_logger("workbook.py")

class PythonObjectLibrary:
    
    # This will create a GUID to register it with Windows, it is unique.
    _reg_clsid_ = pythoncom.CreateGuid()

    # Register the object as an EXE file, the alternative is an DLL file (INPROC_SERVER)
    _reg_clsctx_ = pythoncom.CLSCTX_LOCAL_SERVER

    # the program ID, this is the name of the object library that users will use to create the object.
    _reg_progid_ = "Python.ObjectLibrary"

    # this is a description of our object library.
    _reg_desc_ = "This is our Python object library."

    # a list of strings that indicate the public methods for the object. If they aren't listed they are conisdered private.
    _public_methods_ = ['SAS', 'SASO']
    
    
    def __init__(self):
        self.count = 0

    # multiply two cell values.
    def pythonMultiply(self, a, b):
        return a * b

    # multiply two cell values.
    def pythonAdd(self, a, b):
        return a + b

    # multiply two cell values.
    def SAS(self, op, close, high, low, date, index, reqUrl ):      
        
        # logger.info("Entering SAS")
        
        data = dict()
        
        data.update({'CP' : close,
                     'HP' : high,
                     'LP' : low,
                     'Date' : date,
                     'index' : index
            })
        
        URL = str()
        
        if bool(reqUrl):
            URL = reqUrl
        else:        
            URL = "http://localhost:5000/listing/update"
            

        self.update_values_async(URL, data)        
        return URL

    # multiply two cell values.
    def SASO(self, q):
        self.count = self.count+1
        self.update_values_async(q, 1)
        return self.count
    
    def update_values_async(self, URL, data):
        async_task = AsyncUpdateTask(q=URL, task_type = data)
        async_task.start()
        return 
if __name__ == '__main__':
    import win32com.server.register
    win32com.server.register.UseCommandLine(PythonObjectLibrary)