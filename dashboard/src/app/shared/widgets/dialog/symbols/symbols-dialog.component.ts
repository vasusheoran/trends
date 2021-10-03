import { Component, Inject } from '@angular/core';
import { MatDialog, MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { IListing } from 'src/app/shared/models/listing';
import { FormControl, Validators } from '@angular/forms';
import { DialogComponent } from 'src/app/shared/widgets/dialog/dialog.component';
import { NONE_TYPE } from '@angular/compiler';
import { ConfigService } from 'src/app/shared/services/config.service';
import { MatSnackBar } from '@angular/material/snack-bar';



@Component({
  selector: 'symbols-dialog',
  templateUrl: 'symbols-dialog.component.html',
})

export class SymbolsDataDialog {


  Company = new FormControl('', [Validators.required]);
  SAS = new FormControl('', [Validators.required]);
  Series = new FormControl('', [Validators.required]);
  Symbol = new FormControl('', [Validators.required]);

  constructor(public dialogRef: MatDialogRef<DialogComponent>,
    private _config: ConfigService,
    private _snack: MatSnackBar,
    @Inject(MAT_DIALOG_DATA) public data: { symbol: IListing, action: string }) {
    if (data.action == "Add") {
      data.symbol = {
        Company: null,
        Symbol: null,
        Series: null,
        SAS: null
      }
    }
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

  onClick(): void {
    if (this.data.action == "Add") {
      this._config.putSymbols(this.data.symbol.SAS, this.data.symbol).subscribe(resp => {
        this.openSnackBar("Successfully added the new symbol.")
      }, err => {
        this.openSnackBar("Error added the new symbol.")
      });
    } else if (this.data.action == "Edit") {
      this._config.patchSymbols(this.data.symbol.SAS, this.data.symbol).subscribe(resp => {
        this.openSnackBar("Successfully edited the symbol.")
      }, err => {
        this.openSnackBar("Error editing the symbol.")
      });
    } else if (this.data.action == "Delete") {
      this._config.deleteSymbols(this.data.symbol.SAS).subscribe(resp => {
        this.openSnackBar("Successfully deleted the symbol.")
      }, err => {
        this.openSnackBar("Error deleting the new symbol.")
      });
    }
    this.dialogRef.close();
  }

  openSnackBar(msg: string): void {
    this._snack.open(msg, 'Close', {
      duration: 2000
    });

  }

}