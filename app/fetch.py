# -*- coding: utf-8 -*-
"""
Created on Sat May 23 00:05:44 2020

@author: vedvp
"""

from flask import Blueprint, jsonify, request
from wrapper import fetch_updated_or_frozen, fetch_symbol_list, paginate, fetch_index_if_set, fetch_data_by_start_end

fetch = Blueprint('fetch', __name__, url_prefix='/fetch')

@fetch.route('/value')
def fetch_values():
    return jsonify(fetch_updated_or_frozen(True))

@fetch.route('/listings', methods = ['GET'])
def fetch_stock_listings():
    return jsonify(fetch_symbol_list())

@fetch.route('/<int:page>/<int:size>')
def fetchHistoricalData(page, size):
    df = paginate(page, size)
    return jsonify(df)

#  Pass SAS and Yahoo index both
@fetch.route('/index', methods = ['GET'])
def get_index_if_present(): 
    return jsonify(fetch_index_if_set())

@fetch.route('/data', methods=['GET', 'POST'])
def fetch_data():
    start = request.args.get('start')
    end = request.args.get('end')
    
    if start and end:
        start = int(start)
        end = int(end)
        
    return fetch_data_by_start_end(start, end)

@fetch.route('/freeze' , methods = ['GET'])
def fetch_freeze():
    return jsonify(fetch_updated_or_frozen(False))
