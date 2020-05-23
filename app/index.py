# -*- coding: utf-8 -*-

from datetime import  datetime
from flask import Blueprint, request
from wrapper import set_current_listing, reset_current_index, fetch_updated_or_frozen, update_values
from wrapper import AsyncUpdateRealTimeTask

listing = Blueprint('listing', __name__, url_prefix='/listing')

#  Pass SAS and Yahoo index both
@listing.route('/set', methods = ['POST'])
def get_index(): 
    data = request.get_json() 
    return set_current_listing(data)

#  Pass SAS and Yahoo index both
@listing.route('/reset', methods = ['POST'])
def reset_index(): 
    return reset_current_index()

@listing.route('/freeze' , methods = ['POST'])
def freeze():
    data = request.get_json()
    if data['Date'] == None:
        data['Date'] = datetime.today().strftime("%m:%d:%Y %H:%M:%S")
    update_values(data, True)
    return fetch_updated_or_frozen(False)

# TODO: Update Timestamps
@listing.route('/update' , methods = ['POST'])
def update_values_by_time():
    data = request.get_json()
    async_task = AsyncUpdateRealTimeTask(task_details=data)
    async_task.start()
    return "Success"