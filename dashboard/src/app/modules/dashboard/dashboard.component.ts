import { Component, OnInit } from '@angular/core';
import { IListing } from 'src/app/shared/models/listing';
import { UpdateResponse, IUpdateResponse } from 'src/app/shared/models/listing-response';
import { WebSocketsService } from 'src/app/shared/services/web-sockets.service';
import { SharedService } from 'src/app/shared/services/shared.service';
import { StockService } from 'src/app/shared/services/stock.service';

import { Subscription } from 'rxjs';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  cards:any;
  isEnabled:any;
  plotLines:[];
  subscription: Subscription;

  constructor(private _socket : WebSocketsService,
    private _shared : SharedService,
    private _stockHelper : StockService) {
      this.isEnabled = this._stockHelper.isPlotLineEnabled;
    }

  ngOnInit(): void { 

    this.subscription=  this._shared.sharedUpdateResponse.subscribe(resp =>{
      // debugger;
      if(typeof resp != "function"){
        this.cards = resp['table'];
      }
    });

  }  

  toggleEnable(card, key){
    this.isEnabled[key] = !this.isEnabled[key];
    this._stockHelper.toggleClickableFields(key, this.isEnabled);
  }

  toggleClass(color1, color2, value){
      if (this.cards.Close.value > value){
        return color1
      }
      else{
        return color2
      }
  }

}
