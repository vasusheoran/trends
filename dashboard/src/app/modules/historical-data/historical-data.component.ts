import { Component, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table'
import { ConfigService } from 'src/app/shared/services/config.service';
import { MatSnackBar, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { MatDialog } from '@angular/material/dialog';
import { FormBuilder, Validators, FormControl } from '@angular/forms';
import { HistoricalDialogComponent } from 'src/app/shared/widgets/dialog/historical/historical-dialog.component';
import { SharedService } from 'src/app/shared/services/shared.service';
import { Listing, IListing } from 'src/app/shared/models/listing';


@Component({
  selector: 'app-historical-data',
  templateUrl: './historical-data.component.html',
  styleUrls: ['./historical-data.component.css']
})
export class HistoricalDataComponent implements OnInit {

  dialogData;
  page: number;
  size: number;
  listing: any;
  displayedColumns: string[] = ['Date', 'CP', 'HP', 'LP', 'Actions'];
  dataSource: MatTableDataSource<ResponseData>;

  dt = new FormControl('', [Validators.required]);
  megedDate: any;

  @ViewChild(MatPaginator, { static: true }) paginator: MatPaginator;

  ngOnInit() {
    this._shared.sharedListing.subscribe(listing => {
      this.listing = listing;
      this.getHistories();
    }, err => {
      this._route.navigateByUrl('symbol');
    })
  }

  getHistories() {
    this._config.getHistories(this.listing['SAS']).subscribe((resp: ResponseData[]) => {

      if (resp.length == 0) {
        this.openSnackBar("Please set the symbol to continue.");
        this._route.navigateByUrl('symbols');
      }

      this.dataSource = new MatTableDataSource<ResponseData>(resp);
      this.dataSource.paginator = this.paginator;

    }, err => {
      if (err.status == 200) {
        this.openSnackBar(err.statusText);
      }
      else
        this.openSnackBar("Server unavailable...");

      this._route.navigateByUrl('symbols');
    });

  }

  constructor(private _config: ConfigService,
    private _route: Router,
    private _snack: MatSnackBar,
    public dialog: MatDialog,
    private _shared: SharedService,
    private formBuilder: FormBuilder) {
    this.dataSource = null;
  }

  horizontalPosition: MatSnackBarHorizontalPosition = 'end';
  verticalPosition: MatSnackBarVerticalPosition = 'bottom';

  openSnackBar(msg?: string, actionName?: string) {
    if (!msg)
      msg = "Unknown Error.";

    this._snack.open(msg, actionName, {
      duration: 3000,
      horizontalPosition: this.horizontalPosition,
      verticalPosition: this.verticalPosition,
    });
  }

  openDailog(action, element) {
    console.log(action + " : : " + element)
    this.dialog.open(HistoricalDialogComponent, {
      data: {
        'history': element,
        'action': action
      }
    });
  }

  mergeHistories() {
    console.log(this.megedDate)
  }
}


export interface ResponseData {
  Date: Date;
  CP: number;
  HP: number;
  LP: number;
}