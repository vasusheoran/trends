# -*- coding: utf-8 -*-
"""
Created on Sat Jul 18 23:10:53 2020

@author: vedvp
"""

from flask import Blueprint, jsonify, request
from wrapper import index, AsyncUpdateRealTimeTask, fetch_expires, fetch_strike_prices, paginate

index_route = Blueprint('index', __name__, url_prefix='/index')

# index = Index()
    
@index_route.route('', methods = ['POST'])
def post():
    data = request.get_json() 
    return index.post(data)
    # return jsonify(index.post(data))
    
@index_route.route('', methods = ['GET'])
def get():
    return index.get()
    
@index_route.route('', methods = ['PUT'])
def put():
    data = request.get_json()
    async_task = AsyncUpdateRealTimeTask(task_details=data)
    async_task.start()
    return "Success"
    
@index_route.route('', methods = ['DELETE'])
def delete():
    return index.delete()
    
@index_route.route('/freeze', methods = ['POST'])
def post_freeze():
    data = request.get_json()
    index.post_freeze(data)
    return get_freeze()
    
@index_route.route('/freeze', methods = ['GET'])
def get_freeze():
    return index.get_freeze()
    
@index_route.route('/name', methods = ['GET'])
def name():    
    return index.name()

@index_route.route('/expiry/<symbol>/<instrument>', methods = ['GET'])
def expiry(symbol, instrument): 
    return jsonify(fetch_expires({'symbol' : symbol, 'instrument' :instrument}))
    
@index_route.route('/strike/<symbol>/<instrument>/<expiry>/<optionType>', methods = ['GET'])
def sp(symbol, instrument, expiry, optionType): 
    return jsonify(fetch_strike_prices({'symbol' : symbol, 
                                        'instument' :instrument,
                                        'expiry' : expiry, 
                                        'optionType' :optionType}))
    
@index_route.route('/history/<int:page>/<int:size>', methods = ['GET'])
def history(page, size):
    df = paginate(page, size)
    return jsonify(df)