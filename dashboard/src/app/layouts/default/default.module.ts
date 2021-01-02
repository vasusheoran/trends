import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DefaultComponent } from './default.component';
import { DashboardComponent } from 'src/app/modules/dashboard/dashboard.component';
import { RouterModule } from '@angular/router';
import { SharedModule } from 'src/app/shared/shared.module';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatDividerModule } from '@angular/material/divider';
import { MatCardModule } from '@angular/material/card';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatMenuModule } from '@angular/material/menu';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatDialogModule } from '@angular/material/dialog';
import { FormsModule, ReactiveFormsModule  } from '@angular/forms';
import { MatSnackBarModule, MAT_SNACK_BAR_DEFAULT_OPTIONS } from '@angular/material/snack-bar';
import { FlexLayoutModule } from '@angular/flex-layout';
import { HistoricalDataComponent } from 'src/app/modules/historical-data/historical-data.component';
import { SettingsComponent } from 'src/app/modules/settings/settings.component';
import { SymbolsComponent } from 'src/app/modules/symbols/symbols.component';

const materialModules = [
  MatMenuModule,
  MatSidenavModule,
  MatDividerModule,
  MatCardModule,
  MatTableModule,
  MatPaginatorModule,
  MatSnackBarModule,
  MatIconModule,
  MatButtonModule,
  MatInputModule,
  MatSelectModule,
  MatDatepickerModule,
  MatNativeDateModule,
  FormsModule,
  ReactiveFormsModule,
  MatFormFieldModule,
  MatDialogModule,
]

@NgModule({
  declarations: [
    DefaultComponent,
    DashboardComponent,
    HistoricalDataComponent,
    SettingsComponent,
    SymbolsComponent
  ],
  imports: [
    CommonModule,
    RouterModule,
    SharedModule,
    FlexLayoutModule,
    ...materialModules
  ],
  exports:[
    ...materialModules
  ],
  providers: [
    {provide: MAT_SNACK_BAR_DEFAULT_OPTIONS, useValue: {duration: 2500}},
  ]
})
export class DefaultModule { }
