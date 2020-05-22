import { Component, OnInit, Input, EventEmitter, OnChanges, SimpleChange, Output, OnDestroy } from '@angular/core';
import * as Highcharts from 'highcharts/highstock';

import HC_exporting from 'highcharts/modules/exporting';
import { IListing } from '../../models/listing';
import { ConfigService } from '../../services/config.service';
import { IUpdateResponse, UpdateResponse } from '../../models/listing-response';
import { MatSnackBar } from '@angular/material/snack-bar';
import { SharedService } from '../../services/shared.service';
import { CountdownTimerService } from '../../services/countdown-timer.service';
import { CountdownSnackbarComponent } from '../countdown-snackbar/countdown-snackbar.component';
import { WebSocketsService } from '../../services/web-sockets.service';
import * as moment from 'moment';
import { StockService } from 'src/app/share/widgets/stock/stock.service';

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
export class StockComponent implements OnInit, OnChanges, OnDestroy {

    @Output() updatedValueForCards:EventEmitter<any> = new EventEmitter();
    
    @Input() listing: EventEmitter<IListing> = new EventEmitter();
    @Input() buy: EventEmitter<boolean> = new EventEmitter();
    @Input() support: EventEmitter<boolean> = new EventEmitter();
    @Input() sell: EventEmitter<boolean> = new EventEmitter();
    @Input() open: EventEmitter<boolean> = new EventEmitter();

    _chart: Highcharts.Chart;
    chartOptions: {};
    Highcharts: typeof Highcharts = Highcharts;
    interval: any;
    openInterval:any;
    currentData: IUpdateResponse;
    plotLineWidth:number = 2;
    isReload:boolean = false;
    OPEN_INTERVAL_PERIOD = 900;
    OPEN_INTERVAL_DELAY:number;

    constructor(private _config:ConfigService,
        private _snack : MatSnackBar,
        private _shared : SharedService,
        private _countdown : CountdownTimerService,
        private _socket : WebSocketsService,
        private _stockHelper : StockService ) {
            this._shared.resetListing(resp => {
                if (resp){
                    this._chart.destroy();
                    this._chart = null;
                }
            });
        }

    ngOnInit(): void { 
        // this.currentData['OP'] = 10;

        this._config.fetchIndexIfSet().subscribe((resp) => {
            if(resp['status'] == 'Success'){
                this.listing = resp['listing'];

                let realTimeData:Array<IUpdateResponse> = resp['data'];
                let historicalData:Array<HistoricalResponse> = resp['historical_data'];

                this.initChart(this.parseData(realTimeData, historicalData),  resp['listing']);                        
                this.isReload = true;
            }else{
                this.isReload = true;
                this._snack.open('Please set a Listing to view chart.');
            }
        });

            
        this._socket.subsribeForUpdates().subscribe((res) =>{
            var resp:UpdateResponse = new UpdateResponse(res);
            if (typeof resp != "function"){
                console.log(resp.count);

                if(this._chart && this._chart.series && this._chart.series[0]){

                    var x;
                    if(resp.Date == null || resp.Date == undefined){
                        x = new Date().valueOf();
                    }else{
                        // x = new Date(resp.Date).valueOf();
                        x = moment(resp.Date, "M:D:YYYY H:mm:ss").valueOf();
                    }
        
                    this._chart.series[0].addPoint([x, resp.CP], true, true);
        
                    this.updatePlotLine(resp);
                    
                    this.currentData = resp;
                    this.updateCardsData();
                }else{
                    this._snack.open("Please set the listing. Updates from server have started.")
                }
            
            }
        });
    }

    updateASyncData(resp){

    }

    setOpenTimer(){
        var d = new Date();
        var seconds = (d.getMinutes() * 60) + d.getSeconds()
        let duration = Math.abs(seconds - ((900 - (seconds % 900)) + seconds));
        setTimeout(() =>{
            this.openInterval = setInterval(() => {
                if(this.currentData && this.currentData.CP){
                    this.currentData.OP = this.currentData.CP;
                }
            },  10000);
        }, 10000);
        
        this._snack.open('Updating sackbar in ' + duration + 'seconds.','Okay', {
            duration:2000
        });
        // 
        // this._snack.openFromComponent(CountdownSnackbarComponent, {duration});
        // this._countdown.start(duration);
    }

    ngOnDestroy(): void {
        this._chart = null;
        clearInterval(this.interval);
        clearInterval(this.openInterval);
    }

    ngOnChanges(changes: import("@angular/core").SimpleChanges): void {

        if (this.isReload && changes['listing'] && changes.listing.currentValue && changes.listing.currentValue['YahooSymbol'] !=null ) {

            if (this._chart) {
                this._chart.destroy();
            }
            this._config.setListing(changes.listing.currentValue).subscribe(resp => {
                this._snack.open('Rendering chart. Please wait ...');
                let realTimeData:Array<IUpdateResponse> = resp['data'];
                let historicalData:Array<HistoricalResponse> = resp['historical_data'];
                this.initChart(this.parseData(realTimeData, historicalData),  changes.listing.currentValue);
            },(err) =>{
                this._snack.open('Error. Unbable to fetch data.');
            });
        }

        if(changes['buy'] || changes['sell'] || changes['support'] || changes['open']){
            let buy = this.buy;
            let sell = this.sell;
            let support = this.support;
            let open = this.open;


            if(this._chart){
                var plotLines = [];
                if(buy && this.currentData.bi){
                    plotLines.push({color: '#ff0000b8', value: this.currentData.bi, width: this.plotLineWidth }); //Green
                }if(sell && this.currentData.bk){
                    plotLines.push({color: '#74992e', value: this.currentData.bk, width: this.plotLineWidth }); // Red
                }if(support && this.currentData.bj){
                    plotLines.push({color: '#0000ff7a', value: this.currentData.bj, width: this.plotLineWidth }); //Voilet
                }if(open && this.currentData.OP){
                    plotLines.push({color: '#554e2bbf', value: this.currentData.OP, width: this.plotLineWidth }) // Custom
                }
                this._chart.update({yAxis: { plotLines: plotLines }}, true);
            }
        }

    }

    parseData(realTimeData: IUpdateResponse[], historicalData:HistoricalResponse[]): any[] {
        realTimeData = realTimeData.filter(
            (thing, i, arr) => arr.findIndex(t => t.Date === thing.Date) === i
        );
        historicalData = historicalData.filter(
            (thing, i, arr) => arr.findIndex(t => t.date === thing.date) === i
        );

        let data = [];
        realTimeData.forEach(element => {
            if(element.CP && element.Date){
                data.push([
                    new Date(element.Date).valueOf(),
                    element.CP
                ]);
            }
        });

        
        historicalData.forEach(element => {
            if(element.CP && element.date){
                data.push([
                    new Date(element.date).valueOf(),
                    element.CP
                ]);
            }
        });

        data.sort((a, b) => {
            return a[0] - b[0];
    
        });
        // this.updateCardsData(data[data.length-1]);

        return data;
    }

    updateCardsData(){
        this.updatedValueForCards.emit({'BI' : this.currentData['bi'], 'BJ': this.currentData['bj'], 
        'BK' : this.currentData['bk'], 'OP' : this.currentData['CP'], 'CP' : this.currentData['CP']})
    }

    // updateData(){

    //     this._config.fetchValues().subscribe((resp:IUpdateResponse) => {
    //         if(resp == null){
    //             clearInterval(this.interval);
    //             console.log("Server Busy. PLease try again after some time.")
    //         }
    //         this.updateCardsData();

    //         var x;
    //         if(resp.Date == null || resp.Date == undefined){
    //             x = new Date().valueOf();
    //         }else{
    //             x = new Date(resp.Date).valueOf()
    //         }

    //         this._chart.series[0].addPoint([x, resp.CP], true, false);

    //         // this.updatePlotLine(resp);
            
    //         this.currentData = resp;
    //     });

    // }

    // updatePlotLine(resp: IUpdateResponse) {
    //     if(this._chart){
    //         var plotLines = [];
    //         if(this.buy && this.currentData.bi != resp.bi){
    //             plotLines.push({color: '#74992e', value: this.currentData.bi, width: this.plotLineWidth }); //Green
    //         }if(this.sell && this.currentData.bk != resp.bk){
    //             plotLines.push({color: '#ff0000b8', value: this.currentData.bi, width: this.plotLineWidth }); // Red
    //         }if(this.support && this.currentData.bj != resp.bj){
    //             plotLines.push({color: '#0000ff7a', value: this.currentData.bi, width: this.plotLineWidth }); //Voilet
    //         }if(this.open && this.currentData.OP != resp.OP){
    //             plotLines.push({color: '#554e2bbf', value: this.currentData.bi, width: this.plotLineWidth }) // Custom
    //         }

    //         if(plotLines.length >0){
    //             this._chart.update({yAxis: { plotLines: plotLines }}, true);
    //             console.log('Modified Plotlines : ' + plotLines);
    //         }
    //     }
    // }

    


    // initChart(data=[], name:IListing) {
    //     this.chartOptions = {
    //         time: {
    //             useUTC: false
    //         },
    //         tooltip:{
    //             formatter: function() {
    //                 var series = '<b>' +  this.points[0].series.name+ ' : </b> ' + '<span>' + Highcharts.numberFormat(this.y, 2, ",");
    //                 var date = '</span><br/><br/><span><b>Time :</b></span><span>' + moment().utc(this.x).format('LLL') + '</span>' ;
    //                 return series + date ;
    //             }
    //         },

    //         rangeSelector: {
    //             buttons: [{
    //                 type: 'all',
    //                 text: 'All'
    //             },{
    //                 count: 5,
    //                 type: 'minute',
    //                 text: '5M'
    //             },{
    //                 count: 30,
    //                 type: 'minute',
    //                 text: '30M'
    //             },{
    //                 count: 1,
    //                 type: 'hour',
    //                 text: '1H'
    //             }],
    //             inputEnabled: false,
    //             selected: 0
    //         },

    //         title: {
    //             text: name.CompanyName
    //         },

    //         exporting: {
    //             enabled: false
    //         },

    //         yAxis: {
    //             formatter: function() {
    //                 return '<b>' +  this.points[0].series.name+ ' : </b> ' +  Highcharts.numberFormat(this.y, 0);
    //             }
    //         },

    //         series: [{
    //             name: 'Close Price',
    //             data: Object.assign([], data)

    //         }]
    //     }
    //     this._shared.nextListing(this.listing);
    //     this._chart = Highcharts.stockChart('canvas', this.chartOptions);
    // }

}
