import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';

interface StringConstructor {
  format: (formatString: string, ...replacement: any[]) => string;
}

@Injectable({
  providedIn: 'root'
})
export class ConfigService {

  private baseUrl:string = "";

  private fetchValuesUrl:string;
  private fetchIndexUrl:string;

  private fetchListingsUrl:string;

  private paginateHistoricalDataUrl:string;

  private freezeBIUrl:string;

  private fetchFrozenUrl:string;

  private addNewRowUrl:string;

  private setIndexUrl:string;

  private resetIndexUrl:string;

  private downloadLogUrl:string;

  private uploadSymbolsUrl:string;

  private deleteRowUrl:string;

  private fetchDataByStartAndEndUrl(start:string, end:string):string
  {
      let query = this.baseUrl  + 'data?start={0}&end={1}';
      return query;
  }


  constructor(private _http: HttpClient) { 
    this.baseUrl = environment.apiUrl;
    this.fetchValuesUrl = this.baseUrl  + 'fetch/value';  
    this.fetchIndexUrl = this.baseUrl  + 'fetch/index';
    this.fetchListingsUrl = this.baseUrl  + 'fetch/listings';
    this.fetchFrozenUrl = this.baseUrl  + 'fetch/freeze'; 
    this.paginateHistoricalDataUrl = this.baseUrl  + 'fetch/';
    this.freezeBIUrl = this.baseUrl  + 'listing/freeze';  
    this.setIndexUrl = this.baseUrl  + 'listing/set';  
    this.resetIndexUrl = this.baseUrl  + 'listing/reset';  
    this.downloadLogUrl = this.baseUrl  + 'download/';  
    this.uploadSymbolsUrl = this.baseUrl  + 'upload';
  }

  fetchValues() {
    return this._http.get(this.fetchValuesUrl).pipe(map(data => data));
  }

  fetchListings() {
    return this._http.get(this.fetchListingsUrl);
  }

  fetchHistoricalData(page, size) {
    const url = this.paginateHistoricalDataUrl + page + '/' + size;
    return this._http.get(url).pipe(map(data => data));
  }

  freezeBI(data) {
    return this._http.post(this.freezeBIUrl, data).pipe(map(data => data));
  }

  fetchFrozenValues() {
    return this._http.get(this.fetchFrozenUrl).pipe(map(data => data));
  }

  addNewRow(ob) {
    return this._http.post(this.addNewRowUrl, ob).pipe(map(data => data));
  }

  setListing(selectedOption) {
    return this._http.post(this.setIndexUrl, selectedOption).pipe(map(data => data));
  }

  resetListing(options) {
    return this._http.post(this.resetIndexUrl, options).pipe(map(data => data));
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
}
