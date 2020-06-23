import { Injectable, EventEmitter, OptionalDecorator } from '@angular/core';
import * as Highcharts from 'highcharts/highstock';
import { LoggerService } from 'src/app/shared/services/logger.service';
import { WebSocketsService } from 'src/app/shared/services/web-sockets.service';
import { ConfigService } from 'src/app/shared/services/config.service';

@Injectable({
  providedIn: 'root'
})
export class StockService {

  _chart: Highcharts.Chart;
  chartOptions = {};
  currentData: any;
  
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

  setChartOptions():{}{        
    var options = {
        scrollbar: {
          liveRedraw: false
        },
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
          zoomType: 'x'
        },

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
    options['title']['text'] = chart['listing']['SAS'];
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
