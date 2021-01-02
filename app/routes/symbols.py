# -*- coding: utf-8 -*-
"""
Created on Sat Jul 18 23:10:53 2020

@author: vedvp
"""

from flask import Blueprint, jsonify, request
from wrapper import Symbol

symbols = Blueprint('symbol', __name__, url_prefix='/symbol')

sb = Symbol()
    
@symbols.route('', methods = ['GET'])
def get():
    li = sb.get()
    return jsonify(li)
    
@symbols.route('', methods = ['POST'])
def post():    
    data = request.get_json() 
    sb.post(data)
    return { 'status' : 'success'}
    
@symbols.route('/<sid>', methods = ['PUT'])
def put(sid):
    data = request.get_json() 
    sb.put(data, sid)
    return { 'status' : 'success'}
    
@symbols.route('/<sid>', methods = ['DELETE'])
def delete(sid):
    sb.delete(sid)
    return { 'status' : 'success'}