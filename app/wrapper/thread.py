# -*- coding: utf-8 -*-
"""
Created on Tue Mar 31 00:41:53 2020

@author: vsheoran
"""
import threading,time
from datetime import datetime, timedelta
from .wrapper import push_notifications
from .utilities import get_logger, real_time_field_names, min_freeze_time, max_freeze_time
from .database import DB
from .index import index

logger = get_logger("thread.py")
real_time_data = dict()

class AsyncUpdateRealTimeTask(threading.Thread):
    real_time_field_names = list()
    def __init__(self, task_details):
        super().__init__()
        self.task_details = task_details
        
        global real_time_field_names
        self.real_time_field_names = real_time_field_names
        
    def run(self):
        logger.info("Inside AsyncUpdateRealTimeTask")
        try:
            task = self.task_details
            if isinstance(task['CP'], str):
                task['CP'] = float(task['CP'])
            if isinstance(task['HP'], str):
                task['HP'] = float(task['HP'])
            if isinstance(task['LP'], str):
                task['LP'] = float(task['LP'])
            # for task in self.task_details:
            should_freeze = False
            if 'reset_freeze_value' in task:
                should_freeze = task['reset_freeze_value']
            else:
                current_time = datetime.strptime(task['Date'], '%m:%d:%Y %H:%M:%S')                
                should_freeze = self.isFreezeEnabled(task, current_time, should_freeze)            
                
            task['Date'] = int(current_time.timestamp() * 1000)

            # Calcuate response if listing matches
            if task['index'] == index.name():
                values, ok = index.put(task)
                push_notifications('updateui', values)
                
            self.update_queue(task)
            # self.update_data(data)               

        except KeyError as err:
            logger.info(f"Key Error : {err}")
            # print(err)
        except Exception as ex:
            logger.info(ex)
            
    def isFreezeEnabled(self, task, current_time, should_freeze):
        global max_freeze_time, min_freeze_time
        
        if current_time <= max_freeze_time and current_time >= min_freeze_time:
            return True
        
        return should_freeze
    
    def update_queue(self, task):
        try:
            global real_time_data   
                    
            data = dict()
            for item in self.real_time_field_names:
                data.update({item : task[item]})        
                
            if data['index'] in real_time_data:
                real_time_data[data['index']].append(data)
            else:
                real_time_data.update({data['index'] : [data]})
                
        except Exception as ex:
            logger.error(ex)
            
    def update_data(self, data):
        # print(data)
        db = DB(data['index'])
        db.save([data])

class FlushToDatabase(threading.Thread):
        
    
    def __init__(self):
        super().__init__()
        self.running = True
        
        global real_time_data
    
    
    def run(self):
        delay = 60
        
        while(self.running):
            time.sleep(delay)
            self.save()
            
        logger.info("Thread to dump data stopped.")
        
    def stop_thread(self):
        self.running = False
        
    
    def save(self):
        global real_time_data
        for key in real_time_data:
            try:           
                data= real_time_data[key] 
                if len(data)>=1:
                    db = DB(key)
                    db.set_real_time_data(data)
                    real_time_data[key] = []
            except Exception as err:
                logger.error(err)
        

class DailyCleanup(threading.Thread):
        
    
    def __init__(self):
        super().__init__()
        self.running = True
        global max_freeze_time, min_freeze_time
    
    
    def run(self):
        global max_freeze_time, min_freeze_time
        x=datetime.today()
        y = x.replace(day=x.day, hour=9, minute=0, second=0, microsecond=0) + timedelta(days=1)
        delta_t=y-x
        secs=delta_t.seconds+1
        
        while(self.running):
            time.sleep(secs)
            logger.info("Cleaning Up")
            max_freeze_time =  datetime.now().replace(hour=9,minute=19,second=58)
            min_freeze_time =  max_freeze_time + timedelta(seconds=1)
            
        logger.info("Thread to dump data stopped.")
        
    def stop_thread(self):
        self.running = False
                        
                    
                
            