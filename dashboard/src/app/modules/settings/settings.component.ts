import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators  } from '@angular/forms';
import { ConfigService } from 'src/app/shared/services/config.service';
import { Router } from '@angular/router';
import { MatSnackBar, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';
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
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.css']
})
export class SettingsComponent implements OnInit {

  listing:any;
  isDerivative:boolean = false;
  symbolForm;

  selectedValue: string;
  selectedCar: string;

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
  



  constructor(private formBuilder : FormBuilder,
    private _config : ConfigService,
    private _route : Router,
    private _snack : MatSnackBar) { 
    this.symbolForm = this.formBuilder.group({
      period: [{value: ''}, Validators.required],
      expiry: [{value: '', disabled: true}, Validators.required],
      instrument: [{value: ''}, Validators.required],
      option: [{value: '', disabled: true}, Validators.required],
      strikePrice: [{value: '', disabled: true}, Validators.required],
    });
  }

  ngOnInit() {  }  

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
    
  horizontalPosition: MatSnackBarHorizontalPosition = 'end';
  verticalPosition: MatSnackBarVerticalPosition = 'bottom';
  
  openSnackBar(msg?:string, actionName?:string) {
    if (!msg)
      msg = "Unknown Error.";

    this._snack.open(msg, actionName, {
      duration: 3000,
      horizontalPosition: this.horizontalPosition,
      verticalPosition: this.verticalPosition,
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

}
