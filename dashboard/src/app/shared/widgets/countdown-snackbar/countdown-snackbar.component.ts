import { Component, OnInit } from '@angular/core';
import { CountdownTimerService } from '../../services/countdown-timer.service';

@Component({
  selector: 'app-countdown-snackbar',
  templateUrl: './countdown-snackbar.component.html',
  styleUrls: ['./countdown-snackbar.component.css']
})
export class CountdownSnackbarComponent {

  timeLeft$ = this.countdown.timeLeft();

  constructor(private countdown: CountdownTimerService) { }

}
