# -*- coding: utf-8 -*-
"""
Created on Fri Mar 27 02:50:36 2020

@author: vsheoran
"""

import os
from flask import Flask, jsonify, request
from flask_cors import CORS
import wrapper as mods
from flask_socketio import SocketIO
from routes.symbols import symbols
from routes.index import index_route

async_task_1 = None

logger = mods.get_logger("main.py") 

app = Flask(__name__) 

app.config['UPLOAD_FOLDER'] = os.getcwd() + os.sep + 'files'
app.config['SECRET_KEY'] = 'secret!'

app.register_blueprint(symbols)
app.register_blueprint(index_route)

CORS(app, resources={r"/*": {"origins": "*"}})


socketio = SocketIO(app, cors_allowed_origins="*")
  
@app.route('/')
def ping():
    return {'status' : "Success"} 

@socketio.on('message')
def handle_message(message):
    logger.info('Connect to client: ' + message) 
    
@socketio.on('updateui')
def connect(message):
    logger.info('Connect to client: ' + message)    

@app.route('/test_socket/<msg>')
def test(msg):
    mods.push_notifications('updateui', {'data' : msg})
    return {'status' : 'Success'}    

@app.errorhandler(400)
def custom400(error):
    response = jsonify({'message': error.description})
    response.status_code = 400
    return response

def set_up(isEnabled):
    global async_task_1, daily_task, socketio
    if isEnabled:
        mods.set_up_socketio(socketio)
        logger.info("Starting thread to dump data.")
        async_task_1 = mods.FlushToDatabase()
        async_task_1.start()
        daily_task = mods.DailyCleanup()
        daily_task.start()
    else:    
        logger.info("Stopping thread.")
        async_task_1.stop_thread()
        daily_task.stop_thread()
    
if __name__ == '__main__':
    # Start update server
    set_up(True)
    socketio.run(app, host='0.0.0.0') 
    set_up(False)