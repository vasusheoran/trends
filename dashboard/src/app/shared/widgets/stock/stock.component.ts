import { Component, OnInit, Input, EventEmitter, OnChanges, SimpleChange, Output, OnDestroy } from '@angular/core';
import * as Highcharts from 'highcharts/highstock';

import HC_exporting from 'highcharts/modules/exporting';
import { IListing } from '../../models/listing';
import { ConfigService } from '../../services/config.service';
import { IUpdateResponse, UpdateResponse } from '../../models/listing-response';
import { MatSnackBar } from '@angular/material/snack-bar';
import { SharedService, ListingResponse } from '../../services/shared.service';
import { CountdownTimerService } from '../../services/countdown-timer.service';
import { WebSocketsService } from '../../services/web-sockets.service';
import { StockService } from 'src/app/share/widgets/stock/stock.service';
import { IRealTimeDataResponse } from '../../models/reat-time-response';
import * as moment from 'moment';

export interface HistoricalResponse{
    CP:number;
    HP:number;
    LP:number;
    date:Date;
}

@Component({
    selector: 'app-widget-stock',
    templateUrl: './stock.component.html',
    styleUrls: ['./stock.component.css']
})
export class StockComponent implements OnInit, OnDestroy {

    currentValues:ListingResponse;
    isUpdated:boolean;
    listing:any;
    
    constructor(private _config:ConfigService,
        private _snack : MatSnackBar,
        private _shared : SharedService,
        private _socket : WebSocketsService,
        private _stockHelper : StockService ) { 
            this.isUpdated = false;
        }

    ngOnInit(): void {         
        this._shared.resetListing(resp => {
            if (resp){
                this._stockHelper.destroyChart();
            }
        });
        this._socket.listen('updateui').subscribe((resp) =>{
            // this.options.fn._shared.ne
            this._stockHelper.addPoint(resp['stocks'], resp['dashboard']['cards']);
            this._shared.nextUpdateResponse(resp['dashboard']);
        });

        // Subscribe to refresh
        this._config.fetchIndexIfSet().subscribe(resp => {
            this.listing = resp['chart']['listing'];
            this._stockHelper.setRealTimeData(resp['chart'], resp['data']['dashboard']['cards']); 
            this._shared.nextUpdateResponse(resp['data']['dashboard']);
            this.isUpdated = true;
            this._shared.nextListing(this.listing);
        }, (err) => {            
            this._snack.open('Please set a Listing to view chart.');
            // this._stockHelper.enableLoading('Please set a Listing to view chart.');
        });

        // Subscribe to 
        this._shared.sharedListing.subscribe(resp => {
            if(typeof resp != 'function' && !this.isUpdated){
                this.listing = resp;
                this._config.setListing(resp).subscribe(resp => {
                this._stockHelper.setRealTimeData(resp['chart'], resp['data']['dashboard']['cards']); 
                this._shared.nextUpdateResponse(resp['data']['dashboard']);
                },(err) =>{
                    this._snack.open('Unbable to fetch data. Please set a Listing to view chart.');
                });
            } 
         },(err) =>{
             this._snack.open('Error. Unbable to fetch data.');
         });
    }
    
    ngOnDestroy(): void {
        this._stockHelper.destroyChart();
    }

}
