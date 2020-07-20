import { Component, OnInit, ViewChild } from '@angular/core';
import { ConfigService } from 'src/app/shared/services/config.service';
import { IListing } from 'src/app/shared/models/listing';
import {MatPaginator} from '@angular/material/paginator';
import {MatTableDataSource} from '@angular/material/table'
import { MatSnackBar, MatSnackBarRef, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';

@Component({
  selector: 'app-symbols',
  templateUrl: './symbols.component.html',
  styleUrls: ['./symbols.component.css']
})
export class SymbolsComponent implements OnInit {

  page:number;
  size:number;
  displayedColumns: string[] = ['Company', 'Symbol', 'SAS', 'Series', 'Actions'];
  dataSource:MatTableDataSource<IListing>;
  
  horizontalPosition: MatSnackBarHorizontalPosition = 'end';
  verticalPosition: MatSnackBarVerticalPosition = 'bottom';

  @ViewChild(MatPaginator, {static: true}) paginator: MatPaginator;

  constructor(private _config : ConfigService,
    private _snack : MatSnackBar) { }

  ngOnInit(): void {
    this._config.getSymbols().subscribe((resp:IListing[]) => {
      this.dataSource = new MatTableDataSource<IListing>(resp);
      setTimeout(() => this.dataSource.paginator = this.paginator);
    }, (err) => {
      this.openSnackBar("Unable to fetch symbols...")
    })
  }

  applyFilter(filterValue: string) {
    filterValue = filterValue.trim(); // Remove whitespace
    filterValue = filterValue.toLowerCase(); // Datasource defaults to lowercase matches
    this.dataSource.filter = filterValue;
  }

  openDailog(action, element){
    console.log(action + " : : " + element)
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
}
