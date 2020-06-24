# -*- coding: utf-8 -*-
"""
Created on Tue Mar 31 00:41:53 2020

@author: vsheoran
"""
import threading
import requests

class AsyncUpdateTask(threading.Thread):
    def __init__(self, q, task_type):
        super().__init__()
        self.q = q
        self.task_type = task_type
    def run(self):
        # print("Details : " + str(self.q))
        # print("Type : " + str(self.task_type))
        ob = self.task_type
        resp = requests.post(self.q, json = ob)
  
        
# =============================================================================
# data = dict()
# 
# data.update({
#     "OP" : 100,
# 	"CP" : 100,
# 	"HP" : 150,
# 	"LP" : 50,
# 	"time" : "12:30:1899 0:0:0",
# 	"index" : "Nifty 50"
# })
# URL = "http://localhost:5000/update"
# async_task = AsyncUpdateTask(q=URL, task_type = data)
# async_task.start()
# =============================================================================
