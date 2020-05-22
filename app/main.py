# -*- coding: utf-8 -*-
"""
Created on Fri Mar 27 02:50:36 2020

@author: vsheoran
"""
import sys, os

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


from wrapper import fetch_updated_or_frozen, fetch_symbol_list, paginate, fetch_index_if_set, update_values,push_notifications
from wrapper import set_current_listing, add_new_rows, push_notifications ,socketio, app,reset_current_index, fetch_data_by_start_end
from utilities import get_logger
from flask import request,jsonify
from datetime import datetime
from async_update_task import AsyncUpdateRealTimeTask, FlushToDatabase, DailyCleanup
from async_es_task import AsyncUpdateSymbolsTask
from werkzeug.utils import secure_filename
from flask_swagger import swagger

logger = get_logger("main.py")

@app.route("/spec")
def spec():
    swag = swagger(app)
    swag['info']['version'] = "1.0"
    swag['info']['title'] = "My API"
    return jsonify(swag)

@app.route('/fetch/value')
def fetch_values():
    return jsonify(fetch_updated_or_frozen(True))

@app.route('/fetch/listings', methods = ['GET'])
def fetch_stock_listings():
    return jsonify(fetch_symbol_list())

@app.route('/fetch/<int:page>/<int:size>')
def fetchHistoricalData(page, size):
    df = paginate(page, size)
    return jsonify(df)

#  Pass SAS and Yahoo index both
@app.route('/fetch/index', methods = ['GET'])
def get_index_if_present(): 
    return jsonify(fetch_index_if_set())

#  Pass SAS and Yahoo index both
@app.route('/set', methods = ['POST'])
def get_index(): 
    data = request.get_json() 
    return set_current_listing(data)

#  Pass SAS and Yahoo index both
@app.route('/reset', methods = ['POST'])
def reset_index(): 
    data = request.get_json() 
    return reset_current_index()

@app.route('/freeze' , methods = ['POST'])
def freeze():
    data = request.get_json()
    if data['Date'] == None:
        data['Date'] = datetime.today().strftime("%m:%d:%Y %H:%M:%S")
    update_values(data, True)
    return fetch_updated_or_frozen(False)

@app.route('/fetch/freeze' , methods = ['GET'])
def fetch_freeze():
    return jsonify(fetch_updated_or_frozen(False))

# TODO: Update Timestamps
@app.route('/update' , methods = ['POST'])
def update_values_by_time():
    data = request.get_json()
    async_task = AsyncUpdateRealTimeTask(task_details=data)
    async_task.start()
    return "Success"

# TODO: Update Timestamps
@app.route('/add' , methods = ['POST'])
def add():    
    time = datetime.now().date
    data = request.get_json()
    data.update({ 'time' :time})
    add_new_rows(data)
    
    return {'status' : 'Success'}

# TODO: Update Timestamps
@app.route('/delete' , methods = ['POST'])
def delete():
    # delete_new_rows(data)
    
    return {'status' : 'Failiure'}

@app.route('/finish' , methods = ['GET'])
def finish():
    # return jsonify(finish_day())
    return {'status' : 'Failiure'}

def allowed_file(filename):
    return '.' in filename and \
           filename.rsplit('.', 1)[1].lower() in ['csv', 'xlsx']

@app.route('/upload', methods=['GET', 'POST'])
def upload_file():
    if request.method == 'POST':
        if 'file' not in request.files:
            return {'status' : 'Failiure', 'msg' :'No file part.'}
        file = request.files['file']
        if file.filename == '':
            return {'status' : 'Failiure', 'msg' :'No file selected.'}
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
            file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))
        
            async_task = AsyncUpdateSymbolsTask(task_details=file.filename)
            async_task.start()
            
            return {'status' : 'Success', 'msg' :'File Saved.'}
        
    return {'status' : 'Failiure', 'msg' :'Unknown.'}

@app.route('/data', methods=['GET', 'POST'])
def fetch_data():
    start = request.args.get('start')
    end = request.args.get('end')
    
    if start and end:
        start = int(start)
        end = int(end)
        
    return fetch_data_by_start_end(start, end)
  
@app.route('/')
def ping():
    return {'status' : "Success"} 
    
def set_up_threading(isEnabled):
    global async_task_1, daily_task
    if isEnabled:
        logger.info("Starting thread to dump data.")
        async_task_1 = FlushToDatabase()
        async_task_1.start()
        daily_task = DailyCleanup()
        daily_task.start()
    else:    
        logger.info("Stopping thread.")
        async_task_1.stop_thread()
        daily_task.stop_thread()
              
@socketio.on('message')
def handle_message(message):
    logger.info('Connect to client: ' + message)
    

@app.route('/test_socket/<msg>')
def test(msg):
    push_notifications('updateui', {'data' : msg})
    return {'status' : 'Success'}    
    
global async_task_1
if __name__ == '__main__':
    # Start update server
    set_up_threading(True)
    socketio.run(app, host='0.0.0.0') 
    set_up_threading(False)