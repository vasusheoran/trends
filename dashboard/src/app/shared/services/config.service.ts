import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';
import { IListing, Listing } from '../models/listing';
import { IHistorical } from '../models/historical';

interface StringConstructor {
  format: (formatString: string, ...replacement: any[]) => string;
}

@Injectable({
  providedIn: 'root'
})
export class ConfigService {

  private baseUrl: string = "";

  private index: string;
  private fetchIndexUrl: string;
  private symbol: string;
  private history: string;
  private freeze: string;
  private downloadLogUrl: string;
  private expiry: string;
  private strike: any;


  constructor(private _http: HttpClient) {
    this.baseUrl = environment.apiUrl + "/api/";
    this.symbol = this.baseUrl + '/symbol';
    this.history = this.baseUrl + '/history';
    this.freeze = this.baseUrl + '/index/freeze';
    this.expiry = this.baseUrl + '/index/expiry';
    this.strike = this.baseUrl + '/index/strike';
  }

  getIndexURL(sas: string) {
    return `${environment.apiUrl}/api/index/${sas}`
  }

  fetchIndex(sas: string) {
    var url = this.getIndexURL("1")
    return this._http.get(url).pipe(map(data => data));
  }

  getSymbols() {
    return this._http.get(this.symbol).pipe(map(data => data));
  }

  postSymbols(symbol: IListing) {
    return this._http.post(this.symbol, symbol).pipe(map(data => data));
  }

  putSymbols(sid: string, symbol: IListing) {
    var url = this.symbol + "/" + sid;
    return this._http.put(url, symbol).pipe(map(data => data));
  }

  deleteSymbols(sid: string) {
    var url = this.symbol + "/" + sid;
    return this._http.delete(url).pipe(map(data => data));
  }

  getHistories() {
    return this._http.get(this.history).pipe(map(data => data));
  }

  mergeHistories(date) {
    return this._http.patch(this.history, { 'date': date }).pipe(map(data => data));
  }

  postHistories(his: IHistorical) {
    return this._http.post(this.history, his).pipe(map(data => data));
  }

  putHistories(hid: Date, symbol: IHistorical) {
    var url = this.history + "/" + hid;
    return this._http.put(url, symbol).pipe(map(data => data));
  }

  deleteHistories(hid: Date) {
    var url = this.history + "/" + hid;
    return this._http.delete(url).pipe(map(data => data));
  }

  freezeBI(data) {
    return this._http.post(this.freeze, data).pipe(map(data => data));
  }

  fetchFrozenValues() {
    return this._http.get(this.freeze).pipe(map(data => data));
  }

  private fetchDataByStartAndEndUrl(start: string, end: string): string {
    let query = this.baseUrl + 'data?start={0}&end={1}';
    return query;
  }

  setListing(selectedOption: Listing) {
    var url = this.getIndexURL("1")
    // var url = this.getIndexURL(selectedOption.SAS)
    console.log(url)
    return this._http.post(url, selectedOption).pipe(map(data => data));
  }

  resetListing(options) {
    return this._http.delete(this.index).pipe(map(data => data));
  }

  downloadLogs(num) {
    return this._http.get(this.downloadLogUrl + num).pipe(map(data => data));
  }

  fetchIndexIfSet() {
    return this._http.get(this.fetchIndexUrl).pipe(map(data => data));
  }

  checkCORS(url) {
    return this._http.get(url).pipe(map(data => data));
  }

  fetchDataByStartAndEnd(start, end) {
    // let url = this.fetchDataByStartAndEndUrl(start, end);
    let url = this.baseUrl + 'data?start=' + start + '&end=' + end;
    return this._http.get(url).pipe(map(data => data));
  }

  fetchExpiry(instrument: string, symbol: string) {
    var url = this.expiry + "/" + symbol + "/" + instrument
    return this._http.get(url).pipe(map(data => data));
  }

  fetchStrikePrices(symbol: string, instrument: string, expiry: string, optionType: string) {
    var url = this.strike + "/" + symbol + "/" + instrument + "/" + expiry + "/" + optionType
    return this._http.get(url).pipe(map(data => data));
  }
}
