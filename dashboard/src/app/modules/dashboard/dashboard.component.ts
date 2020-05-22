import { Component, OnInit } from '@angular/core';
import { IListing } from 'src/app/shared/models/listing';
import { UpdateResponse, IUpdateResponse } from 'src/app/shared/models/listing-response';
import { WebSocketsService } from 'src/app/shared/services/web-sockets.service';
import { SharedService, ListingResponse, IUpdatedValues, ValuesData } from 'src/app/shared/services/shared.service';
import { StockService } from 'src/app/share/widgets/stock/stock.service';


export interface PeriodicElement {
  name: string;
  position: number;
  weight: number;
  symbol: string;
}

const ELEMENT_DATA: PeriodicElement[] = [
  {position: 1, name: 'Hydrogen', weight: 1.0079, symbol: 'H'},
  {position: 2, name: 'Helium', weight: 4.0026, symbol: 'He'},
  {position: 3, name: 'Lithium', weight: 6.941, symbol: 'Li'},
  {position: 4, name: 'Beryllium', weight: 9.0122, symbol: 'Be'},
  {position: 5, name: 'Boron', weight: 10.811, symbol: 'B'},
  {position: 6, name: 'Carbon', weight: 12.0107, symbol: 'C'},
  {position: 7, name: 'Nitrogen', weight: 14.0067, symbol: 'N'},
  {position: 8, name: 'Oxygen', weight: 15.9994, symbol: 'O'},
  {position: 9, name: 'Fluorine', weight: 18.9984, symbol: 'F'},
  {position: 10, name: 'Neon', weight: 20.1797, symbol: 'Ne'},
];

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  
  listing:any;
  cards:any;
  values:any;
  isEnabled:any;
  plotLines:[];

  
  displayedColumns: string[] = ['position', 'name', 'weight', 'symbol'];
  dataSource = ELEMENT_DATA;

  constructor(private _socket : WebSocketsService,
    private _shared : SharedService,
    private _stockHelper : StockService) {
      this.isEnabled = this._stockHelper.isPlotLineEnabled;
    }

  ngOnInit(): void { 
    this._shared.resetListing(resp => {
        if (resp){
            this.listing = null;
        }
    }); 

    this._shared.sharedUpdateResponse.subscribe(resp =>{
      // debugger;
      this.cards = resp['cards'];
      this.values = resp['table'];
    });

    this._socket.emit('message', "Connected.");
  }  

  setSelectedListing(event){
    this.listing = event;
  }

  toggleEnable(card, key){
    this.isEnabled[key] = !this.isEnabled[key];
    this._stockHelper.toggleClickableFields(key, this.isEnabled);
  }

}
