import { Component, OnInit, Input, EventEmitter, OnChanges, SimpleChange, Output, OnDestroy } from '@angular/core';
import * as Highcharts from 'highcharts/highstock';

import HC_exporting from 'highcharts/modules/exporting';
import { ConfigService } from '../../services/config.service';
import { MatSnackBar, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';
import { SharedService, ListingResponse } from '../../services/shared.service';
import { WebSocketsService } from '../../services/web-sockets.service';
import { StockService } from 'src/app/share/widgets/stock/stock.service';
import { Router } from '@angular/router';

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
    alertStatus:boolean = false;
    
    constructor(private _config:ConfigService,
        private _snack : MatSnackBar,
        private _shared : SharedService,
        private _socket : WebSocketsService,
        private _stockHelper : StockService,
        private _route : Router ) { 
            this.isUpdated = false;
        }

    ngOnInit(): void { 
        this._socket.listen('updateui').subscribe((resp) =>{
            this._stockHelper.addPoint(resp['stocks'], resp['dashboard']['cards']);
            this._shared.nextUpdateResponse(resp['dashboard']);
        });

        // Subscribe to refresh
        this._config.fetchIndex().subscribe(resp => {
            this.listing = resp['symbol'];
            this._stockHelper.setRealTimeData(resp['data'], resp['values']['dashboard']['cards']); 
            this._shared.nextUpdateResponse(resp['values']['dashboard']);
            this.isUpdated = true;
            this._shared.nextListing(this.listing);
        }, (err) => {     
            if(err.status == 200 || err.status == 500){
                this.openSnackBar("Please set the Symbol to continue...");
                this._route.navigateByUrl('settings');
            }
            else
                this.openSnackBar("Server unavailable...");
        });
    }
    
    ngOnDestroy(): void {
        this._stockHelper.destroyChart();
    }
    
    horizontalPosition: MatSnackBarHorizontalPosition = 'end';
    verticalPosition: MatSnackBarVerticalPosition = 'bottom';
    
    openSnackBar(msg?:string, actionName?:string) {
        if (!msg)
            msg = "Unknown Error.";
    
        if(this._snack._openedSnackBarRef){
            this._snack._openedSnackBarRef.afterDismissed().subscribe(() => {
                this._snack.open(msg, actionName, {
                    duration: 1000,
                    horizontalPosition: this.horizontalPosition,
                    verticalPosition: this.verticalPosition,
                });
            });
        }else{
            this._snack.open(msg, actionName, {
                duration: 1000,
                horizontalPosition: this.horizontalPosition,
                verticalPosition: this.verticalPosition,
            });
        }
    }

}
