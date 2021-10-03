import { Component, OnInit, Inject } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';
import { DialogComponent } from 'src/app/shared/widgets/dialog/dialog.component';
import { MatDialog, MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { ConfigService } from 'src/app/shared/services/config.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { IHistorical } from 'src/app/shared/models/historical';

@Component({
  selector: 'app-historical-dialog',
  templateUrl: './historical-dialog.component.html',
  styleUrls: ['./historical-dialog.component.css']
})
export class HistoricalDialogComponent implements OnInit {


  cp = new FormControl('', [Validators.required]);
  hp = new FormControl('', [Validators.required]);
  lp = new FormControl('', [Validators.required]);
  dt = new FormControl('', [Validators.required]);

  constructor(public dialogRef: MatDialogRef<DialogComponent>,
    private _config: ConfigService,
    private _snack: MatSnackBar,
    @Inject(MAT_DIALOG_DATA) public data: { history: IHistorical, action: string }) {
    if (data.action == "Add") {
      data.history = {
        Date: null,
        HP: null,
        LP: null,
        CP: null
      }
    }
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

  onClick(): void {
    // if (this.data.action == "Add"){
    //   this._config.postHistories(this.data.history).subscribe(resp => {
    //     this.openSnackBar("Successfully added the new history.")
    //   }, err => {
    //     this.openSnackBar("Error added the new history.")
    //   });
    // }else if(this.data.action == "Edit"){
    //   this._config.putHistories(this.data.history.Date, this.data.history).subscribe(resp => {
    //     this.openSnackBar("Successfully edited the history.")
    //   }, err => {
    //     this.openSnackBar("Error editing the history.")
    //   });

    // }else if(this.data.action == "Delete"){
    //   this._config.deleteHistories(this.data.history.Date).subscribe(resp => {
    //     this.openSnackBar("Successfully deleted the history.")
    //   }, err => {
    //     this.openSnackBar("Error deleting the new history.")
    //   });
    // }
    this.dialogRef.close();
  }

  openSnackBar(msg: string): void {
    this._snack.open(msg, 'Close', {
      duration: 2000
    });

  }

  ngOnInit(): void {
  }

}
