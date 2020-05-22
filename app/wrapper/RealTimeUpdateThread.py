# -*- coding: utf-8 -*-
"""
Created on Tue Apr 14 19:44:14 2020

@author: vedvp
"""
from flask_socketio import emit
import threading,logging
import random, time

class RandomThread(threading.Thread):
    def __init__(self):
        self.delay = 1
        super(RandomThread, self).__init__()
    def randomNumberGenerator(self):
        """
        Generate a random number every 1 second and emit to a socketio instance (broadcast)
        Ideally to be run in a separate thread?
        """
        #infinite loop of magical random numbers
        print("Making random numbers")
        while True:
            logging.info("random")
            number = random.randrange(10)
            emit('update', {'data': number})
            time.sleep(self.delay)    
        
    def run(self):
        self.randomNumberGenerator()

