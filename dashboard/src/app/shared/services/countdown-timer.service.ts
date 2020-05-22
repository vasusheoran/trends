import { Injectable, OnInit } from '@angular/core';
import { timer, Subject, Observable, interval, fromEvent, merge, empty } from 'rxjs';
import { switchMap, scan, take, mapTo, startWith, takeWhile, map } from 'rxjs/operators';
import { CountdownSnackbarComponent } from '../widgets/countdown-snackbar/countdown-snackbar.component';
// const COUNTDOWN_SECONDS = 10;
let COUNTDOWN_SECONDS;
// elem refs
const interval$ = interval(1000).pipe(mapTo(-1));

@Injectable({
  providedIn: 'root'
})
export class CountdownTimerService {
  startTimer = new Subject<number>();
  // timer:any;
  s:number;
  // ngAfterViewInit():void{
  //   this.timer = this.startTimer.pipe(
  //     startWith(true),
      
  //   switchMap(val => (val ? interval$ : empty())),
  //   scan((acc, curr) => (curr ? curr + acc : acc), COUNTDOWN_SECONDS),
  //   takeWhile(v => v >= 0)
  //   );
  // }
  timer = this.startTimer.pipe(
    startWith(true),
    
  switchMap(val => (val ? interval$ : empty())),
  scan((acc, curr) => (curr ? curr + acc : acc), COUNTDOWN_SECONDS),
  takeWhile(v => v >= 0)
  );
  

  start(time: number) {
    const seconds = Math.floor(time / 1000);
    this.startTimer.next(seconds);
  }

  timeLeft(): Observable<number> {
    return this.timer;
  }
}
