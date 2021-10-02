import { Component, OnInit } from '@angular/core';
import { ConfigService } from 'src/app/shared/services/config.service';
import { SharedService } from 'src/app/shared/services/shared.service';
import { StockService } from 'src/app/shared/services/stock.service';
import { Router } from '@angular/router';
import { MatSnackBar, MatSnackBarRef, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';
import { ActivatedRoute } from '@angular/router';
import { WebSocketsService } from '../../shared/services/web-sockets.service';

import { Observable, Subscription } from 'rxjs';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  sas: string
  cards: any;
  isEnabled: any;
  plotLines: [];
  subscription: Subscription;
  updateUI: Observable<any>;
  updateUISub: Subscription;

  constructor(private _config: ConfigService,
    private _shared: SharedService,
    private activatedRoute: ActivatedRoute,
    private _snack: MatSnackBar,
    private _socket: WebSocketsService,
    private _stockHelper: StockService) {
    this.isEnabled = this._stockHelper.isPlotLineEnabled;
  }

  ngOnInit(): void {
    this.activatedRoute.params.subscribe(params => {
      this.sas = params['sas'];

      debugger;
      if (this.sas == null || this.sas == "") {
        this.openSnackBar("Please set the symbol to continue.");
      } else {
        // this._socket.enable()
      }
    })

    this._config.fetchIndex(this.sas).subscribe(resp => {
      this.cards = resp['summary']
    })


    // this._socket.listen('updateui')

    // this._shared.sharedIsChartEnabled.subscribe(resp => {
    //   if (this.updateUISub != undefined) {
    //     this.updateUISub.unsubscribe();
    //   }

    //   if (resp) {
    //     console.log("Updating cards");
    //     this.updateUISub = this.updateUI.subscribe((resp) => {
    //       this._shared.nextUpdateResponse(resp['dashboard']);
    //     });
    //   }
    // });

  }

  toggleEnable(card, key) {
    this.isEnabled[key] = !this.isEnabled[key];
    this._stockHelper.toggleClickableFields(key, this.isEnabled);
  }

  toggleClass(color1, color2, value) {
    if (this.cards.close > value) {
      return color1
    }
    else {
      return color2
    }
  }

  toggleRSIClass(value) {
    if (value >= 70) {
      return 'darkblue'
    } else if (value >= 50) {
      return 'green'
    } else {
      return 'red'
    }
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

}
