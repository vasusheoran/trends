import { Component, Input, forwardRef } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import { MatDatepickerInputEvent } from '@angular/material/datepicker';
import * as _moment from 'moment';
const moment = _moment;

// <mat-form-field>
//     <mat-label>Choose a date</mat-label>
//     <input matInput [matDatepicker]="picker" [value]="dateValue" (dateInput)="addEvent('input', $event)" [placeholder]="title">
//     <mat-datepicker-toggle matSuffix [for]="picker"></mat-datepicker-toggle>
//     <mat-datepicker #picker></mat-datepicker>
// </mat-form-field>

@Component({
  selector: 'cva-date',
  template: `
  <mat-form-field>
    <mat-label>Choose a date</mat-label>
    <input matInput [matDatepicker]="picker" [value]="dateValue" (dateInput)="addEvent('input', $event)" [placeholder]="title">
    <mat-datepicker-toggle matSuffix [for]="picker"></mat-datepicker-toggle>
    <mat-datepicker #picker></mat-datepicker>
    </mat-form-field>
  `,
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => CvaDateComponent),
      multi: true
    }
  ]
})
export class CvaDateComponent implements ControlValueAccessor {

  @Input()
  _dateValue; // notice the '_'

  @Input() title: string;

  get dateValue() {
    return moment(this._dateValue, 'YYYY/MM/DD HH:mm:ss');
  }

  set dateValue(val) {
    this._dateValue = moment(val).format('YYYY/MM/DD HH:mm:ss');
    this.propagateChange(this._dateValue);
  }

  addEvent(type: string, event: MatDatepickerInputEvent<Date>) {
    console.log(event.value);
    this._dateValue = moment(event.value, 'YYYY/MMM/DD HH:mm:ss');
  }

  writeValue(value: any) {
    if (value !== undefined) {
      this._dateValue = moment(value, 'YYYY/MM/DD HH:mm:ss');
    }
  }
  propagateChange = (_: any) => { };

  registerOnChange(fn) {
    this.propagateChange = fn;
  }

  registerOnTouched() { }
}