import { Component, OnInit, Output } from '@angular/core';
import { EventEmitter } from '@angular/core';
import { Router, NavigationStart } from '@angular/router';
import { SharedService, IListingResponse } from '../../services/shared.service';
import { Listing, IListing } from '../../models/listing';
import { ConfigService } from '../../services/config.service';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {

  @Output() toggleSideBarOutput: EventEmitter<any> = new EventEmitter();

  currentUrl:string;
  currentListing:any;

  constructor(private _router : Router,
    private _shared : SharedService,
    private _config : ConfigService,
    private _snack : MatSnackBar) {
      this.currentListing = null;
      this._shared.sharedListing.subscribe((resp) =>{
        this.currentListing =resp;

      })
  }

  ngOnInit(): void { 
    this._router.events.subscribe((val:NavigationStart) => {
      this.currentUrl = this._router.url;  
    });  
  }

  toggleSideBar(){
    this.toggleSideBarOutput.emit(null);
  }

  resetIndex(){
      this._config.resetListing(null).subscribe(resp => {
        window.location.reload();
      },(err) =>{
        this._snack.open('Error resetting listing.');
        this._router.navigate(['/']);
      });
      
  }

}
