import { EventEmitter, Input } from '@angular/core';
import { Component, OnInit, Output,  } from '@angular/core';
import { FormControl } from '@angular/forms';
import {Observable} from 'rxjs';
import {map, startWith} from 'rxjs/operators';
import {IListing, Listing} from 'src/app/shared/models/listing';
import { ConfigService } from '../../services/config.service';
import { SharedService } from '../../services/shared.service';
import { MatSnackBar, MatSnackBarVerticalPosition, MatSnackBarHorizontalPosition } from '@angular/material/snack-bar';

export interface User {
  Company: string;
  Series: string;
  SAS: string;
  Symbol: string;
}

@Component({
  selector: 'app-widget-autocomplete',
  templateUrl: './autocomplete.component.html',
  styleUrls: ['./autocomplete.component.css']
})
export class AutocompleteComponent implements OnInit {

  @Output() selectedListing:EventEmitter<any> = new EventEmitter();
  @Input() date:Date;
  
  myControl = new FormControl();

  options: IListing[];
  
  filteredOptions: Observable<User[]>;

  constructor(private _config: ConfigService, 
    private _shared : SharedService,
    private _snack : MatSnackBar){ }

  ngOnInit() {
    this.openSnackBar('Fetching Symbols...')
    this._config.getSymbols().subscribe((resp:Array<User>) =>{
      this.options = resp;        
      this.filteredOptions = this.myControl.valueChanges
        .pipe(
          startWith(''),
          map(value => typeof value === 'string' ? value : value.name),
          map(name => name ? this._filter(name) : this.options.slice())
        );
    },(err) => {   
      if(err.status == 200)
        this.openSnackBar(err.statusText, "Go to Settings");
      else
        this.openSnackBar('Unable to fetch Symbols. Server Unavailable.')
    });
  }

  public getLisiting(option:typeof Listing){
    this.selectedListing.emit(option);
    // this._shared.nextListing(option);
  }  

  displayFn(user: User): string {
    return user && user.Company ? user.Company : '';
  }

  private _filter(name: string): User[] {
    const filterValue = name.toLowerCase();
    return this.options.filter(option => option.Company.toLowerCase().indexOf(filterValue) === 0);
  }
    
  horizontalPosition: MatSnackBarHorizontalPosition = 'end';
  verticalPosition: MatSnackBarVerticalPosition = 'bottom';
  
  openSnackBar(msg:string, actionName?:string) {
    if (!msg)
      msg = "Unknown Error.";

    this._snack.open(msg, actionName, {
      duration: 1000,
      horizontalPosition: this.horizontalPosition,
      verticalPosition: this.verticalPosition,
    });
  }
}
