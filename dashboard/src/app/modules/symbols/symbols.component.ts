import { Component, OnInit, ViewChild } from '@angular/core';
import { ConfigService } from 'src/app/shared/services/config.service';
import { FormBuilder, Validators  } from '@angular/forms';
import { Router } from '@angular/router';
import { IListing } from 'src/app/shared/models/listing';
import {MatPaginator} from '@angular/material/paginator';
import {MatTableDataSource} from '@angular/material/table'
import { SelectionModel} from '@angular/cdk/collections';
import { MatSnackBar, MatSnackBarRef, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';
import { MatDialog } from '@angular/material/dialog';
import { SymbolsDataDialog } from 'src/app/shared/widgets/dialog/symbols/symbols-dialog.component';
import { MatSort } from '@angular/material/sort';
import * as _moment from 'moment';

interface Period {
  value: string;
  viewValue: string;
}

interface Instrument {
  value: string;
  viewValue: string;
}
@Component({
  selector: 'app-symbols',
  templateUrl: './symbols.component.html',
  styleUrls: ['./symbols.component.css']
})
export class SymbolsComponent implements OnInit {

  isSymbolSelected:boolean;
  page:number;
  size:number;
  displayedColumns: string[] = ['Company', 'Symbol', 'SAS', 'Series', 'Actions'];
  dataSource:MatTableDataSource<IListing>;
  selection = new SelectionModel<IListing>(false, []);
  
  horizontalPosition: MatSnackBarHorizontalPosition = 'end';
  verticalPosition: MatSnackBarVerticalPosition = 'bottom';

  @ViewChild(MatPaginator, {static: true}) paginator: MatPaginator;
  @ViewChild(MatSort, {static: true}) sort: MatSort;

  constructor(private _config : ConfigService,
    private formBuilder : FormBuilder,
    private _route : Router,
    private _snack : MatSnackBar,
    public dialog: MatDialog) { 
      this.isSymbolSelected = false;

      this.symbolForm = this.formBuilder.group({
        period: [{value: ''}, Validators.required],
        expiry: [{value: '', disabled: true}, Validators.required],
        instrument: [{value: ''}, Validators.required],
        option: [{value: '', disabled: true}, Validators.required],
        strikePrice: [{value: '', disabled: true}, Validators.required],
      });
    }

  ngOnInit(): void {
    this.refreshSymbols("Done fetching symbols",)
  }

  refreshSymbols(msg): void {
    this._config.getSymbols().subscribe((resp:IListing[]) => {
      this.dataSource = new MatTableDataSource<IListing>(resp);
      // setTimeout(() => this.dataSource.paginator = this.paginator);
      this.dataSource.paginator = this.paginator;
      this.dataSource.sort = this.sort;
      this.openSnackBar(msg)
    }, (err) => {
      this.openSnackBar("Unable to fetch symbols...")
    })
  }

  applyFilter(filterValue: string) {
    filterValue = filterValue.trim(); // Remove whitespace
    filterValue = filterValue.toLowerCase(); // Datasource defaults to lowercase matches
    this.dataSource.filter = filterValue;
  }
  
  openDailog(action, element) {
    console.log(action + " : : " + element)
    this.dialog.open(SymbolsDataDialog, {
      data: {
        'symbol': element,
        'action': action
      }
    });
  }

  openSnackBar(msg?:string, actionName?:string) {
    if (!msg)
      msg = "Unknown Error.";

    this._snack.open(msg, actionName, {
      duration: 3000,
      horizontalPosition: this.horizontalPosition,
      verticalPosition: this.verticalPosition,
    });
  }

  // Setting symbol
  listing:any;
  isDerivative:boolean = false;
  symbolForm;

  periods: Period[] = [
    {value: 'day', viewValue: '1 Day'},
    {value: '7days', viewValue: '7 Days'},
    {value: 'week', viewValue: '1 Week'},
    {value: '2weeks', viewValue: '2 Weeks'},
    {value: 'month', viewValue: '1 Month'},
    {value: '3months', viewValue: '3 Months'}
  ];

  intstruments: Instrument[] = [
    {value: 'FUTCUR', viewValue: 'Future Currency'},
    {value: 'OPTCUR', viewValue: 'Options Currency'}
  ];

  expiries;
  isExpiryEnabled:boolean = false;

  options = [
    {value: 'CE', viewValue: 'Call'},
    {value: 'PE', viewValue: 'Put'}];

  strikePrices;


  setSelectedListing(event){
    this.listing = event;

    if(this.listing.Series == "Derivative"){
      this.isDerivative = true;
      
      this.symbolForm.get('period').setValue('3months');
      this.symbolForm.get('instrument').setValue('FUTCUR');
      this.onInstrumentChange({'value' : 'FUTCUR'})

    }else{
      this.openSnackBar('Please Wait...');
      this._config.setListing(this.listing).subscribe(resp =>{
        this.openSnackBar(resp['msg']);
        this._route.navigateByUrl('dashboard');
      },err =>{
        console.log(err)
        this.openSnackBar(err.message);
      });
    }
  }

  onSubmit(data) {
    this.openSnackBar('Please Wait...');
    // Setting symbol as SAS
    data['symbol'] = this.listing['Symbol']
    this.listing['options'] = data;
    this._config.setListing(this.listing).subscribe(resp =>{
      this.openSnackBar(resp['msg']);
      this._route.navigateByUrl('dashboard');
    },err =>{
      this.openSnackBar(err.error.message);
    });
  }

  onInstrumentChange(event){
    
    this.symbolForm.get('expiry').setValue(null);
    this.symbolForm.get('expiry').disable();
    this._config.fetchExpiry(event.value, this.listing['Symbol']).subscribe((resp) => {
        this.expiries = resp;
        this.symbolForm.get('expiry').enable();
    },err =>{
      this._snack.open(err.message);
    });

  }

  onExpiryChange(){
    this.symbolForm.get('option').setValue(null);
    if(this.symbolForm.get('instrument').value == "OPTCUR"){      
      this.symbolForm.get('option').enable();
    }else{
      this.symbolForm.get('option').disable();
    }
  }

  onOptionTypeChange(event){ 
    this._snack.open("Fetching strike prices ...");
    this.symbolForm.get('strikePrice').setValue(null);
    this.symbolForm.get('strikePrice').disable();  
    var options = {
      "optionType" : event.value,
      "expiry" : this.symbolForm.get('expiry').value,
      "instrument" : this.symbolForm.get('instrument').value,
      "symbol" : this.listing['Symbol']
    } 
    this._config.fetchStrikePrices(this.listing['Symbol'], this.symbolForm.get('instrument').value, this.symbolForm.get('expiry').value, event.value).subscribe(resp=> {
      this.strikePrices = resp;
      this.symbolForm.get('strikePrice').enable();
      this._snack.open("Done.");
    },err =>{
      this._snack.open(err.message);
    });
  }

  selectSymbol(element){
    console.log(element)
    this.isSymbolSelected = true;
    this.setSelectedListing(element);
  }

}
