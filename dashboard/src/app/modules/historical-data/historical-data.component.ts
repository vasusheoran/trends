import { Component, OnInit, ViewChild } from '@angular/core';
import {MatPaginator} from '@angular/material/paginator';
import {MatTableDataSource} from '@angular/material/table'
import { ConfigService } from 'src/app/shared/services/config.service';
import { SharedService } from 'src/app/shared/services/shared.service';
import { MatSnackBar, MatSnackBarRef, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';
import { Router } from '@angular/router';

@Component({
  selector: 'app-historical-data',
  templateUrl: './historical-data.component.html',
  styleUrls: ['./historical-data.component.css']
})
export class HistoricalDataComponent implements OnInit {

  page:number;
  size:number;
  displayedColumns: string[] = ['Date', 'CP', 'HP', 'LP'];
  dataSource:MatTableDataSource<ResponseData>;
  // dataSource = new MatTableDataSource<PeriodicElement>(ELEMENT_DATA);

  @ViewChild(MatPaginator, {static: true}) paginator: MatPaginator;

  ngOnInit() { 
    this._config.fetchHistoricalData(1, 10).subscribe((resp:ResponseData[]) => {

      if(resp.length ==0){
        this.openSnackBar("Please set the symbol to continue.");
        this._route.navigateByUrl('symbols');
      }
      
      this.dataSource = new MatTableDataSource<ResponseData>(resp);
      this.dataSource.paginator = this.paginator;
    
    },err =>{
      if(err.status == 200){
          this.openSnackBar(err.statusText);
      }
      else
          this.openSnackBar("Server unavailable...");  
          
      this._route.navigateByUrl('symbols');
    });
  }

  constructor(private _config : ConfigService,
    private _route : Router,
    private _snack : MatSnackBar) { 
      this.dataSource = null; 
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


export interface ResponseData {
  Date:Date;
  CP: number;
  HP: number;
  LP: number;
}