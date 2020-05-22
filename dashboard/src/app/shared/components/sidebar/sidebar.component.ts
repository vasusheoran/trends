import { Component, OnInit } from '@angular/core';
import { Router, NavigationStart } from '@angular/router';
import { SharedService } from '../../services/shared.service';
import { Listing } from '../../models/listing';
import { MatDialog } from '@angular/material/dialog';
import { DialogComponent } from '../../widgets/dialog/dialog.component';
import { ConfigService } from '../../services/config.service';
import {MatSnackBar, MatSnackBarRef} from '@angular/material/snack-bar';

export interface FrozenValues{
  data:{
    CP:number;
    HP:number;
    LP:number;
    date:Date;
    bi:number;
  }
}

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.css']
})
export class SidebarComponent implements OnInit {

  currentUrl:string;
  toggleDashBoardTools:boolean;
  toggleHistoricDataTools:boolean;
  listing:Listing;
  snackBarRef:any;
  frozenValues:FrozenValues;
  dialogData;
  
  constructor(private _router : Router,
    private _shared  : SharedService,
    public dialog: MatDialog,
    private _config : ConfigService,
    private _snack : MatSnackBar  ) {
    this.toggleDashBoardTools = true;
    this.toggleHistoricDataTools = false;
    this.dialogData = {};
    
  }

  ngOnInit(): void {
    this._router.events.subscribe((val:NavigationStart) => {
      // see also 
        let currentUrl = this._router.url;  
        if(currentUrl == '/'){
         this.toggleDashBoardTools = true;
         this.toggleHistoricDataTools = false; 
        }else if(currentUrl == '/historica-data'){
          this.toggleHistoricDataTools = true;
          this.toggleDashBoardTools = false;
        }
    });  
    
    this._shared.sharedListing.subscribe((resp) =>{
      if(typeof resp != 'function')
        this.listing = resp;
    });
    this.fetchFreezeValues();
  }

  fetchFreezeValues(){
    this._config.fetchFrozenValues().subscribe((resp) =>{
      this.dialogData = resp['data'];
    });
  }

  updateFreezeTime(result){
    this._config.freezeBI(result).subscribe((res) => {
      this.snackBarRef = this._snack.open("Buy forzen successfully. Updating Values");
      this.fetchFreezeValues();
    },(err) =>{
      this.snackBarRef = this._snack.open('Error in Freezing Buy.','Close',{
        duration:4000
      });
    });
  }

  addRow(result){
    this._config.addNewRow(result).subscribe((res) => {
      this._router.navigate['/'];
      this.snackBarRef = this._snack.open('Row Added successfully.','Close',{
        duration:4000
      });
    },(err) =>{
      this.snackBarRef = this._snack.open('Error in adding row.','Close',{
        duration:4000
      });
    });
  }

  deleteRow(){
    this._config.deleteRow().subscribe((res) => {
      this._router.navigate['/'];
      this.snackBarRef = this._snack.open('Row Deleter successfully.','Close');
    },(err) =>{
      this.snackBarRef = this._snack.open('Error in deleteing row.','Close');
    });
  }

  openDailog(task, isFreeze): void {

    if(this.listing ==null || this.listing.SASSymbol == null){
      this.snackBarRef = this._snack.open('Please set the Stock Listing to continue.','Close',{
        duration:4000
      });
    }else{
      this.dialogData['task'] = task;
      this.dialogData['data'] = this.frozenValues;
      const dialogRef = this.dialog.open(DialogComponent, {
        width: '250px',
        data: this.dialogData
      });

      
      dialogRef.afterClosed().subscribe(result => {
          if(result)
            delete result['tempDate'];
          if(result['CP'] == undefined || result['HP'] == undefined || result['LP'] == undefined){
            this.snackBarRef = this._snack.open('Error. Please set CP/HP/LP.','Close',{
              duration:4000
            });
            return;
          }
          result['index'] = this.listing.SASSymbol;

          if(isFreeze)
            this.updateFreezeTime(result);
          else
            this.addRow(result);
      });
    }
  }
}