# -*- coding: utf-8 -*-
"""
Created on Sat Jul 18 23:10:53 2020

@author: vedvp
"""

from flask import Blueprint, jsonify
from wrapper import Symbol

symbols = Blueprint('symbol', __name__, url_prefix='/symbol')

sb = Symbol()
    
@symbols.route('/')
def fetch_values():
    li = sb.get()
    return jsonify(li)