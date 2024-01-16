import { Component, OnInit } from '@angular/core';
import { ConfigService } from 'src/app/shared/services/config.service';
import { SharedService } from 'src/app/shared/services/shared.service';
import { StockService } from 'src/app/shared/services/stock.service';
import { Router } from '@angular/router';
import { MatSnackBar, MatSnackBarRef, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';
import { ActivatedRoute } from '@angular/router';
import { WebSocketsService } from '../../shared/services/web-sockets.service';
import { environment } from 'src/environments/environment';

import { TickerClient, BidirectionalStream } from '../../generated/ticker_pb_service'
import { SummaryRequest, SummaryResponse } from '../../generated/ticker_pb'

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
  socketClient: WebSocketsService;
  summaryResponseStream1: BidirectionalStream<SummaryRequest, SummaryResponse>;

  constructor(private _config: ConfigService,
    private _shared: SharedService,
    private _route: Router,
    private activatedRoute: ActivatedRoute,
    private _snack: MatSnackBar,
    private _stockHelper: StockService) {
    this.isEnabled = this._stockHelper.isPlotLineEnabled;
    this.socketClient = new WebSocketsService();
  }

  public ngOnDestroy() {
    this.socketClient.close();
  }

  ngOnInit(): void {
    this.activatedRoute.params.subscribe(params => {
      this.sas = params['sas'];

      if (this.sas == null || this.sas == "") {
        this.openSnackBar("Please set the symbol to continue.");
        this._route.navigateByUrl('symbols')
      } else {
        this.socketClient.enable(environment.socketUrl + this.sas);
      }
    });
    var sub: Subscription;
    this.socketClient.getEventListener().subscribe(event => {
      if (event.type == "message") {
        if (event.data != null && event.data != undefined && event.data["summary"] != null) {
          this.cards = event.data["summary"];
        }
      }
      if (event.type == "close") {
        console.info("The handlers connection has been closed");
      }
      if (event.type == "open") {
        console.info("The handlers connection has been opened");
      }
    });

    this._config.fetchIndex(this.sas).subscribe(resp => {
      this.cards = resp['summary']

    }, err => {
      this._snack.open("Failed to fetch symbol. Please upload csv.");
      this._route.navigateByUrl('symbols')
    });
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

  toggleEMA5(greaterThenCloseClr: string, greaterThenAvgClr: string, lessThenCloseClr: string) {
    debugger;

    if (this.cards.close > this.cards.ema_5) {
      return greaterThenCloseClr
    } else if (this.cards.close > this.cards.average) {
      return greaterThenAvgClr
    } else {
      return lessThenCloseClr
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
