import { Injectable, EventEmitter, OptionalDecorator } from '@angular/core';
import * as Highcharts from 'highcharts/highstock';
import { LoggerService } from 'src/app/shared/services/logger.service';
import { IUpdateResponse, UpdateResponse } from 'src/app/shared/models/listing-response';
import { HistoricalResponse } from 'src/app/shared/widgets/stock/stock.component';
import * as moment from 'moment';
import { MatSnackBar } from '@angular/material/snack-bar';
import { IRealTimeDataResponse, RealTimeDataResponse } from 'src/app/shared/models/reat-time-response';
import { WebSocketsService } from 'src/app/shared/services/web-sockets.service';
import { SharedService } from 'src/app/shared/services/shared.service';
import { IListing } from 'src/app/shared/models/listing';
import { ConfigService } from 'src/app/shared/services/config.service';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class StockService {

  _chart: Highcharts.Chart;
  chartOptions = {};
  currentData: any;

//  buy:boolean;
//  support:boolean;
//  sell:boolean;
//  high:boolean;
 plotLinesOptions:any;
 isPlotLineEnabled:any;
 data: any[];
 seriesLength:number;
 isRefreshRequired:boolean;

  constructor(
    private _logger : LoggerService,
    private _socket: WebSocketsService,
    private _config : ConfigService) {
      this.isPlotLineEnabled =  {'Buy' : true, 'Support' : false, 'Sell' : false, 'Min_High' : false};
      this.plotLinesOptions = {'Buy' : { 'color' : '#74992e', 'lineWidth' : 2, 'value' : 0, 'width': 2, dashStyle: 'longdashdot'},
                                'Sell' : { 'color' : '#ff0000b8', 'lineWidth' : 2, 'value' : 0, 'width': 2, dashStyle: 'longdashdot'},
                                'Support' : { 'color' : '#0000ff7a', 'lineWidth' : 2, 'value' : 0, 'width': 2, dashStyle: 'longdashdot'},
                                'Min_High' : { 'color' : '#e976d8', 'lineWidth' : 2, 'value' : 0, 'width': 2, dashStyle: 'longdashdot'}
                              }
     }

  toggleClickableFields(field, isPlotLineEnabled){
    this.isPlotLineEnabled = isPlotLineEnabled;
    this.updatePlotLines();
  }
  updatePlotLines() {

    var plotLines = [];
    var plotLineWidth = 2;

    // plotLines.push({color: '#74992e', value: this.currentData[0]['value'], width: plotLineWidth, dashStyle: 'longdashdot' }); //Green
    if(this.isPlotLineEnabled["Buy"]){
      plotLines.push(this.plotLinesOptions['Buy'])
      // plotLines.push({color: this.plotLinesOptions['Buy']['color'], value: this.currentData[0].value, width: plotLineWidth, dashStyle: 'longdashdot' }); //Green
    }if(this.isPlotLineEnabled["Sell"]){
      plotLines.push(this.plotLinesOptions['Sell'])
        // plotLines.push({color: '#ff0000b8', value: this.currentData[1].value, width: plotLineWidth, dashStyle: 'longdashdot' }); // Red
    }if(this.isPlotLineEnabled["Support"]){
      plotLines.push(this.plotLinesOptions['Support'])
        // plotLines.push({color: '#0000ff7a', value: this.currentData[2].value, width: plotLineWidth, dashStyle: 'longdashdot' }); //Voilet
    }
    if(this.isPlotLineEnabled["Min_High"]){
      plotLines.push(this.plotLinesOptions['Min_High'])
    }
    this._chart.update({yAxis: { plotLines: plotLines }}, true);
  }


  afterSetExtremes(e) {

    // var chart = Highcharts.charts[0];

    // chart.showLoading('Loading data from server...');
    
    console.log("start : " + e.min + ", end : " + e.max);
    var chart = Highcharts.charts[0];
    chart.showLoading('Loading data from server...');
    let url = "http://localhost:5000/"  + 'data?start=' + e.min + '&end=' + e.max;
    chart.options['fn']._http.get(url).subscribe(resp =>{
      var data = resp['data'];
      // chart.series[0].update({
      //   data: data
      // });
      chart.series[0].setData(data, false);
      chart.hideLoading();
    }, (error) =>{
      console.error(error);
    })
    // chart.options['fn'].fetchDataByStartAndEnd(e.min, e.max).subscribe(resp =>{
    //   console.log(resp);
    //   // chart.series[0].setData(data);
    //   // chart.hideLoading();

    // });
   
   
  }

  setChartOptions():{}{        
    var options = {
        scrollbar: {
          liveRedraw: false
        },
        // type: 'spline',
      // animation: Highcharts.svg, // don't animate in old IE
        // marginRight: 10,
        time: {
            // useUTC: false,
            timezone: 'Asia/Kolkata'
            // timezoneOffset: 330
        },

        rangeSelector: {
            buttons: [{
                count: 1,
                type: 'minute',
                text: '1M'
            },{
                count: 5,
                type: 'minute',
                text: '5M'
            },{
                count: 30,
                type: 'minute',
                text: '30M'
            },{
                count: 1,
                type: 'hour',
                text: '1H'
            },{
              type: 'all',
              text: 'All'
            },],
            selected: 4,
            inputEnabled: false
        },

        title: {
            text: ''
        },

        exporting: {
            enabled: false
        },

        fn:{
          sockets:this._socket,         
          fetchDataByStartAndEnd: this._config.fetchDataByStartAndEnd,
        },

        chart: {
          events: {
              // load: function(){
              //   this.options.fn.sockets.listen('updateui').subscribe((resp) =>{
              //     // this.options.fn._shared.ne
              //     console.log(resp);

              //     if(this.series[0].length > 10000)
              //       this.series[0].addPoint(resp['update'], true, true);
              //     else
              //       this.series[0].addPoint(resp['update']);
              //     // this.options.fn.updatePlotLine(resp);
              //   });
              // }
          },
          zoomType: 'x'
        },

        // xAxis: {
        //   events: {
        //     afterSetExtremes: this.afterSetExtremes
        //   }, 
        //   minRange: 60 * 1000
        // },

        series: [{
          name: 'Close',
          data: [],
          marker: {
            enabled: true,
            radius: 1
          },
          tooltip: {
              valueDecimals: 2
          },
          states: {
              hover: {
                  markerRadiusPlus : 1
              }
          },
          lineWidth: 1,
          dataGrouping: {
            enabled: true,
            groupPixelWidth : 4
          }
        }]
    };
    return options;
  }

  getCurrentValues(){
    return this.currentData;
  }

  setRealTimeData(chart, current) { 
      
    if (chart['data'].length ==0){
      chart['data'] = [current];
    }
    var options= this.setChartOptions();
    options['title']['text'] = chart['listing']['CompanyName'];
    options['series'][0]['data'] = chart['data'];
    this._chart = Highcharts.stockChart('canvas', options);
    this.currentData = current;
    
    this.setCurrentData(current);
    this.updatePlotLines();
  }

  setCurrentData(data){
    data.forEach(element => {
      this.plotLinesOptions[element["key"]]['value'] = element['value'];
    });
  }

  destroyChart(){
    if(this._chart)
      this._chart.destroy();
  }

  addPoint(update, plotLineData){
    if(this._chart){      
      if(this._chart.series[0].options['data'] > 10000){
        this._chart.series[0].addPoint(update, true, true);
      }
      else
        this._chart.series[0].addPoint(update);
    }
    this.updatePlotLineWithResponse(plotLineData)
  }

  updatePlotLineWithResponse(resp) {
    var plotLines= [];
    var plotLineWidth = 2;
    var flag = false;

    if(this.currentData){

      if(this.isPlotLineEnabled["Buy"] && this.plotLinesOptions['Buy'].value != resp[0].value){
        this.plotLinesOptions['Buy']['value'] = resp[0].value;
        plotLines.push(this.plotLinesOptions['Buy']);
        flag = true;
      }if(this.isPlotLineEnabled["Support"] && this.plotLinesOptions['Support'].value != resp[1].value){
        this.plotLinesOptions['Support']['value'] = resp[1].value;
        plotLines.push(this.plotLinesOptions['Support']); 
        flag = true;
      }if(this.isPlotLineEnabled["Sell"] && this.plotLinesOptions['Sell'].value != resp[2].value){
        this.plotLinesOptions['Sell']['value'] = resp[2].value;
        plotLines.push(this.plotLinesOptions['Sell']); //Green
        flag = true;
      }if(this.isPlotLineEnabled["Min_High"] && this.plotLinesOptions['Min_High'].value != resp[3].value){
        this.plotLinesOptions['Min_High']['value'] = resp[3].value;
        plotLines.push(this.plotLinesOptions['Min_High']); //Green
        flag = true;
      }
      if (flag)
        this._chart.update({yAxis: { plotLines: plotLines }}, true);
    }
    this.setCurrentData(resp);
  }

}
