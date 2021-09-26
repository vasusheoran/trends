# -*- coding: utf-8 -*-
"""
Created on Sat Jul 18 23:10:53 2020

@author: vedvp
"""

from flask import Blueprint, jsonify, request
from wrapper import History, paginate

history = Blueprint('history', __name__, url_prefix='/history')

his = History()
    
@history.route('', methods = ['GET'])
def get():
    li = his.get()
    return jsonify(li)
    
@history.route('', methods = ['POST'])
def post():    
    data = request.get_json() 
    his.post(data)
    return {'status': 'success'}
    
@history.route('/<hid>', methods = ['PUT'])
def put(hid):
    data = request.get_json() 
    his.put(data, hid)
    return {'status': 'success'}

@history.route('/<hid>', methods = ['DELETE'])
def delete(hid):
    his.delete(hid)
    return {'status': 'success'}

@history.route('/', methods = ['PATCH'])
def patch():
    data = request.get_json()
    his.patch(data['date'])
    return {'status': 'success'}
