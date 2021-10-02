import { Component, OnInit } from '@angular/core';
import { ConfigService } from 'src/app/shared/services/config.service';
import { SharedService } from 'src/app/shared/services/shared.service';
import { StockService } from 'src/app/shared/services/stock.service';
import { Router } from '@angular/router';
import { MatSnackBar, MatSnackBarRef, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';
import { ActivatedRoute } from '@angular/router';
import { WebSocketsService } from '../../shared/services/web-sockets.service';

import { TickerClient, ResponseStream } from '../../generated/ticker_pb_service'
import { SummaryRequest, SummaryReply } from '../../generated/ticker_pb'

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
  tickerClient: TickerClient;
  summaryResponseStream: ResponseStream<SummaryReply>;

  constructor(private _config: ConfigService,
    private _shared: SharedService,
    private _route: Router,
    private activatedRoute: ActivatedRoute,
    private _snack: MatSnackBar,
    private _stockHelper: StockService) {
    this.isEnabled = this._stockHelper.isPlotLineEnabled;
    this.tickerClient = new TickerClient("http://localhost:8080");
  }

  ngOnInit(): void {
    this.activatedRoute.params.subscribe(params => {
      this.sas = params['sas'];

      if (this.sas == null || this.sas == "") {
        this.openSnackBar("Please set the symbol to continue.");
        this._route.navigateByUrl('symbols')
      } else {
        // this._socket.enable()Ticker
        var req = new SummaryRequest();
        // req.setSas("1");
        req.setSas(this.sas);
        this.summaryResponseStream = this.tickerClient.getSummary(req);
      }
    })

    this._config.fetchIndex(this.sas).subscribe(resp => {
      this.cards = resp['summary']
    })

    this.summaryResponseStream.on("data", (message) => {
      console.log(message)
      this.cards = message.toObject()
    })
    // this.summaryResponseStream.on('data', (message: SummaryReply) => {
    //   console.log(message)
    // })

  }
  getSummaries() {
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
