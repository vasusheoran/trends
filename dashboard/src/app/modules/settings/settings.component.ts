import { Component, OnInit } from '@angular/core';
import { FormBuilder  } from '@angular/forms';
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
  options;any;
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


  constructor(private formBuilder : FormBuilder,
    private _config : ConfigService,
    private _route : Router,
    private _snack : MatSnackBar) { 
    this.symbolForm = this.formBuilder.group({
      period: '',
      expiry: '',
      instrument: ''
    });
  }

  ngOnInit() {  }  

  setSelectedListing(event){
    this.listing = event;

    if(this.listing.Series == "Derivative"){
      this.isDerivative = true;
      
      this.symbolForm.get('period').setValue('3months');
      this.symbolForm.get('instrument').setValue('FUTCUR');

    }else{
      this.openSnackBar('Please Wait...');
      this._config.setListing(this.listing).subscribe(resp =>{
        this.openSnackBar(resp['msg']);
        this._route.navigateByUrl('dashboard');
      },err =>{
        this.openSnackBar(err.statusText);
      });
    }
  }

  onSubmit(data) {
    debugger;
    this.openSnackBar('Please Wait...');
    data['expiry'] = _moment(data['expiry']).format("DDMMMYYYY")
    this.listing['options'] = data;
    this._config.setListing(this.listing).subscribe(resp =>{
      this.openSnackBar(resp['msg']);
      this._route.navigateByUrl('dashboard');
    },err =>{
      this.openSnackBar(err.statusText);
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

}
