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
    this.symbol = `${environment.apiUrl}/api/symbol`
  }

  getIndexURL(sas: string) {
    return `${environment.apiUrl}/api/index/${sas}`
  }

  getFreezeURL(sas: string) {
    return `${environment.apiUrl}/api/index/${sas}/freeze`
  }


  getSymbolURL(sas: string) {
    return `${this.symbol}/${sas}`
  }

  getHistoryURL(sas: string) {
    return `${environment.apiUrl}/api/history/${sas}`
  }

  setListing(listing: Listing) {
    return this._http.post(this.getIndexURL(listing.SAS), null).pipe(map(data => data));
  }

  fetchIndex(sas: string) {
    return this._http.get(this.getIndexURL(sas)).pipe(map(data => data));
  }

  getSymbols() {
    return this._http.get(this.symbol).pipe(map(data => data));
  }

  putSymbols(sas: string, symbol: IListing) {
    return this._http.put(this.getSymbolURL(sas), symbol).pipe(map(data => data));
  }

  patchSymbols(sas: string, symbol: IListing) {
    return this._http.patch(this.getSymbolURL(sas), symbol).pipe(map(data => data));
  }

  deleteSymbols(sas: string) {
    return this._http.delete(this.getSymbolURL(sas)).pipe(map(data => data));
  }

  uploadFile(sas: string, file: File) {
    const formData: FormData = new FormData();
    formData.append('file_name', file, sas);
    return this._http.post(this.getHistoryURL(sas), formData).pipe(map(data => data));
  }

  getHistories(sas: string) {
    return this._http.get(this.getHistoryURL(sas)).pipe(map(data => data));
  }

  freezeBI(sas: string, data) {
    return this._http.patch(this.getFreezeURL(sas), data).pipe(map(data => data));
  }

  // fetchFrozenValues() {
  //   return this._http.get(this.freeze).pipe(map(data => data));
  // }

  resetListing(options) {
    return this._http.delete(this.index).pipe(map(data => data));
  }
}
