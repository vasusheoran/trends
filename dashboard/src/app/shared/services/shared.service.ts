import { Injectable } from '@angular/core';
import { Listing, IListing } from '../models/listing';
import { BehaviorSubject } from 'rxjs';
import { IUpdateResponse } from '../models/listing-response';

export class ListingResponse implements IListingResponse{
  data: number[];
  values: IUpdatedValues;

  constructor(data: Partial<ListingResponse>){
    Object.assign(this, data);
  }
}

export interface IListingResponse {
  data:   number[];
  values: IUpdatedValues;
}

export interface IUpdatedValues {
  buy:         ValuesData;
  close:       ValuesData;
  sell:         ValuesData;
  support:       ValuesData;
  date:        ValuesData;
  "ema.CP.10": ValuesData;
  "ema.CP.5":  ValuesData;
  "ema.CP.50": ValuesData;
  high:        ValuesData;
  low:         ValuesData;
  "min.HP.2":  ValuesData;
  open:        ValuesData;
  pe:          ValuesData;
  rsi:         ValuesData;
}

export interface ValuesData{
  name?:  string;
  value?: number;
}

@Injectable({
  providedIn: 'root'
})
export class SharedService {

  private resp:ListingResponse;

  private listing = new BehaviorSubject(Listing);
  sharedListing = this.listing.asObservable();

  private updateResponse = new BehaviorSubject(Object);
  sharedUpdateResponse = this.updateResponse.asObservable();

  private reset:BehaviorSubject<Boolean>  = new BehaviorSubject(null);
  sharedResetListing = this.reset.asObservable();

  nextListing(listing){
    this.listing.next(listing);
  }

  nextUpdateResponse(resp){
    this.updateResponse.next(resp);
  }

  resetListing(val){
    this.reset.next(val);
  }

  constructor() { }
}
