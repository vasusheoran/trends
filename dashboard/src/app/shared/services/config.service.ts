import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';
import { IListing } from '../models/listing';

interface StringConstructor {
  format: (formatString: string, ...replacement: any[]) => string;
}

@Injectable({
  providedIn: 'root'
})
export class ConfigService {

  private baseUrl:string = "";

  private getIndex:string;
  private fetchIndexUrl:string;
  private getSymbol:string;
  private putSymbol:string;
  private postSymbol:string;
  private deleteSymbol:string;
  private getHistorical:string;
  private postFreeze:string;
  private getFreeze:string;
  private addNewRowUrl:string;
  private postIndex:string;
  private deleteIndex:string;
  private downloadLogUrl:string;
  private uploadSymbolsUrl:string;
  private deleteRowUrl:string;
  private getExpiry:string;
  private getStrike: any;


  constructor(private _http: HttpClient) { 
    this.baseUrl = environment.apiUrl;
    this.getIndex = this.baseUrl  + 'index';   
    this.postIndex = this.baseUrl  + 'index';  
    this.deleteIndex = this.baseUrl  + 'index';  
    this.getSymbol = this.baseUrl  + 'symbol';
    this.putSymbol = this.baseUrl  + 'symbol';
    this.postSymbol = this.baseUrl  + 'symbol';
    this.deleteSymbol = this.baseUrl  + 'symbol';
    this.getFreeze = this.baseUrl  + 'index/freeze'; 
    this.postFreeze = this.baseUrl  + 'index/freeze'; 
    this.getHistorical = this.baseUrl  + 'index/history/';
    // this.downloadLogUrl = this.baseUrl  + 'download/';  
    // this.uploadSymbolsUrl = this.baseUrl  + 'upload';
    this.getExpiry = this.baseUrl  + 'index/expiry';
    this.getStrike = this.baseUrl  + 'index/strike';
  }

  fetchIndex() {
    return this._http.get(this.getIndex).pipe(map(data => data));
  }

  getSymbols() {
    return this._http.get(this.getSymbol).pipe(map(data => data));
  }

  postSymbols(symbol:IListing) {
    return this._http.post(this.postSymbol, symbol).pipe(map(data => data));
  }

  putSymbols(sid:string, symbol:IListing) {
    var url = this.putSymbol + "/" + sid;
    return this._http.put(url, symbol).pipe(map(data => data));
  }

  deleteSymbols(sid:string) {
    var url = this.putSymbol + "/" + sid;
    return this._http.delete(url).pipe(map(data => data));
  }

  fetchHistoricalData(page, size) {
    const url = this.getHistorical + page + '/' + size;
    return this._http.get(url).pipe(map(data => data));
  }

  freezeBI(data) {
    return this._http.post(this.postFreeze, data).pipe(map(data => data));
  }

  fetchFrozenValues() {
    return this._http.get(this.getFreeze).pipe(map(data => data));
  }

  private fetchDataByStartAndEndUrl(start:string, end:string):string
  {
      let query = this.baseUrl  + 'data?start={0}&end={1}';
      return query;
  }

  addNewRow(ob) {
    return this._http.post(this.addNewRowUrl, ob).pipe(map(data => data));
  }

  setListing(selectedOption) {
    return this._http.post(this.postIndex, selectedOption).pipe(map(data => data));
  }

  resetListing(options) {
    return this._http.delete(this.deleteIndex).pipe(map(data => data));
  }

  downloadLogs(num) {
    return this._http.get(this.downloadLogUrl + num).pipe(map(data => data));
  }

  uploadSymbols(file: File) {
    const fd = new FormData;
    fd.append('file', file, file.name);
    return this._http.post(this.uploadSymbolsUrl, fd).pipe(map(data => data));
  }

  deleteRow() {
    return this._http.post(this.deleteRowUrl, null).pipe(map(data => data));
  }

  fetchIndexIfSet(){
    return this._http.get(this.fetchIndexUrl).pipe(map(data => data));
  }
  
  checkCORS(url){
    return this._http.get(url).pipe(map(data => data));
  }

  fetchDataByStartAndEnd(start, end){
    // let url = this.fetchDataByStartAndEndUrl(start, end);
    let url = this.baseUrl  + 'data?start=' + start + '&end=' + end;
    return this._http.get(url).pipe(map(data => data));
  }
  
  fetchExpiry(instrument: string, symbol:string) {
    var url = this.getExpiry + "/" + symbol + "/" + instrument
    return this._http.get(url).pipe(map(data => data));
  }

  fetchStrikePrices(symbol: string, instrument: string, expiry: string, optionType: string) {
    var url = this.getStrike + "/" + symbol + "/" + instrument + "/" + expiry + "/" + optionType
    return this._http.get(url ).pipe(map(data => data));
  }
}
