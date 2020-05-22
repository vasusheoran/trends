import { Component, OnInit, Inject, Input } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA} from '@angular/material/dialog';
import { FrozenValues } from '../../models/frozen-values';
import {FormControl, Validators} from '@angular/forms';

export interface DialogData {
  CP:     number;
  Date:   string;
  HP:     number;
  LP:     number;
  bi:     number;
  tempDate : Date;
  task:string;
}


@Component({
  selector: 'app-widget-dialog',
  templateUrl: './dialog.component.html',
  styleUrls: ['./dialog.component.css']
})
export class DialogComponent{

  close = new FormControl('', [Validators.required]);
  high = new FormControl('', [Validators.required]);
  low = new FormControl('', [Validators.required]);

  constructor(
    public dialogRef: MatDialogRef<DialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: DialogData) {

      console.log(this.data);
      let time = new Date();
      data.tempDate = time;
      if(!data.Date){
        data.Date = time.getMonth() + ':' + time.getDate() + ':' + time.getFullYear() + " " +
        time.getHours() + ':' + time.getMinutes() + ':' + time.getSeconds();
      }
    }

  onNoClick(): void {
    this.dialogRef.close();
  }
}