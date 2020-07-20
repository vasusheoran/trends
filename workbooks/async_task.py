# -*- coding: utf-8 -*-
"""
Created on Tue Mar 31 00:41:53 2020

@author: vsheoran
"""
import threading
import requests
import logging
from logging import handlers

def get_logger(filename):
    base_path = "D:/wsl/data/trends/"
    log_format = "[%(levelname)s] | %(asctime)s | %(thread)s | Line : %(lineno)d | [%(name)s] | %(message)s"
    
    logger = logging.getLogger(str(filename))
    logger.handlers = []
    
    # To override the default severity of logging
    logger.setLevel('DEBUG')
    
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

logger = get_logger("async_task.py")

class AsyncUpdateTask(threading.Thread):
    def __init__(self, q, task_type):
        super().__init__()
        self.q = q
        self.task_type = task_type
    def run(self):
        try:
            ob = self.task_type
            resp = requests.put(self.q, json = ob)
        except:
            logger.error("Error updating values")
