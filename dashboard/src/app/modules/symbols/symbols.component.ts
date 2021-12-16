import { Component, OnInit, ViewChild } from '@angular/core';
import { ConfigService } from 'src/app/shared/services/config.service';
import { FormBuilder, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { IListing, Listing } from 'src/app/shared/models/listing';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table'
import { SelectionModel } from '@angular/cdk/collections';
import { MatSnackBar, MatSnackBarRef, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';
import { MatDialog } from '@angular/material/dialog';
import { SymbolsDataDialog } from 'src/app/shared/widgets/dialog/symbols/symbols-dialog.component';
import { MatSort } from '@angular/material/sort';
import * as _moment from 'moment';
import { SharedService } from 'src/app/shared/services/shared.service';

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

  page: number;
  size: number;
  displayedColumns: string[] = ['Company', 'Symbol', 'SAS', 'Series', 'Actions'];
  dataSource: MatTableDataSource<IListing>;
  selection = new SelectionModel<IListing>(false, []);

  horizontalPosition: MatSnackBarHorizontalPosition = 'end';
  verticalPosition: MatSnackBarVerticalPosition = 'bottom';

  @ViewChild(MatPaginator, { static: true }) paginator: MatPaginator;
  @ViewChild(MatSort, { static: true }) sort: MatSort;

  constructor(private _config: ConfigService,
    private formBuilder: FormBuilder,
    private _route: Router,
    private _snack: MatSnackBar,
    private _shared: SharedService,
    public dialog: MatDialog) {

    this.symbolForm = this.formBuilder.group({
      period: [{ value: '' }, Validators.required],
      expiry: [{ value: '', disabled: true }, Validators.required],
      instrument: [{ value: '' }, Validators.required],
      option: [{ value: '', disabled: true }, Validators.required],
      strikePrice: [{ value: '', disabled: true }, Validators.required],
    });
  }

  ngOnInit(): void {
    this.refreshSymbols()
  }

  refreshSymbols(): void {
    this._config.getSymbols().subscribe((resp: IListing[]) => {
      this.dataSource = new MatTableDataSource<IListing>(resp);
      // setTimeout(() => this.dataSource.paginator = this.paginator);
      this.dataSource.paginator = this.paginator;
      this.dataSource.sort = this.sort;
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

  handleFileSelect(sas: string, files: FileList) {
    this._config.uploadFile(sas, files[0]).subscribe(resp => {
      this.openSnackBar("Uploaded file successfully");
    }, err => {
      console.log(err.error);
      this.openSnackBar("Failed to upload file");
    })
  }

  openSnackBar(msg?: string, actionName?: string) {
    if (!msg)
      msg = "Unknown Error 1.";

    this._snack.open(msg, actionName, {
      duration: 3000,
      horizontalPosition: this.horizontalPosition,
      verticalPosition: this.verticalPosition,
    });
  }

  // Setting symbol
  listing: Listing;
  isDerivative: boolean = false;
  symbolForm;


  setSelectedListing(event) {
    this.listing = event;
    this.openSnackBar('Please Wait...');
    this._config.setListing(this.listing).subscribe(resp => {
      this.openSnackBar("Symbol '" + this.listing.SAS + "' set successfully.");
      this._shared.nextListing(this.listing);
      this._route.navigateByUrl('dashboard/' + this.listing.SAS);
    }, err => {
      console.log(err)
      this.openSnackBar(err.error);
    });
  }

  selectSymbol(element) {
    this.setSelectedListing(element);
  }
}
