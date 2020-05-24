import { Component, OnInit } from '@angular/core';
import { IListing } from 'src/app/shared/models/listing';
import { UpdateResponse, IUpdateResponse } from 'src/app/shared/models/listing-response';
import { WebSocketsService } from 'src/app/shared/services/web-sockets.service';
import { SharedService } from 'src/app/shared/services/shared.service';
import { StockService } from 'src/app/share/widgets/stock/stock.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  cards:any;
  values:any;
  isEnabled:any;
  plotLines:[];

  constructor(private _socket : WebSocketsService,
    private _shared : SharedService,
    private _stockHelper : StockService) {
      this.isEnabled = this._stockHelper.isPlotLineEnabled;
    }

  ngOnInit(): void { 
    this._shared.sharedUpdateResponse.subscribe(resp =>{
      // debugger;
      if(typeof resp != "function"){
        this.cards = resp['cards'];
        this.values = resp['table'];
      }
    });

    this._socket.emit('updateui', "event name : updateui");
  }  

  toggleEnable(card, key){
    this.isEnabled[key] = !this.isEnabled[key];
    this._stockHelper.toggleClickableFields(key, this.isEnabled);
  }

}
