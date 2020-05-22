import { NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';

import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { MatListModule } from '@angular/material/list';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatDividerModule } from '@angular/material/divider';
import { MatDialogModule } from '@angular/material/dialog';
 
// const config: SocketIoConfig = { url: 'http://192.168.99.101:5000', options: {'origins': '*'} };


import { ReactiveFormsModule, FormsModule } from '@angular/forms';

import { FlexLayoutModule } from '@angular/flex-layout';
import { RouterModule } from '@angular/router';

import { HeaderComponent } from './components/header/header.component';
import { FooterComponent } from './components/footer/footer.component';
import { SidebarComponent } from './components/sidebar/sidebar.component';
import { StockComponent } from './widgets/stock/stock.component';
import { HighchartsChartModule } from 'highcharts-angular';
import { CardComponent } from './widgets/card/card.component';
import { AutocompleteComponent } from './widgets/autocomplete/autocomplete.component';
import { HttpClientModule, HttpClient } from '@angular/common/http';
import { DialogComponent } from './widgets/dialog/dialog.component';
import { CountdownSnackbarComponent } from './widgets/countdown-snackbar/countdown-snackbar.component';

@NgModule({
  declarations: [
    HeaderComponent,
    FooterComponent,
    SidebarComponent,
    StockComponent,
    StockComponent,
    CardComponent,
    AutocompleteComponent,
    DialogComponent,
    CountdownSnackbarComponent
  ],
  imports: [
    CommonModule,
    MatDividerModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    MatInputModule,
    MatSelectModule,
    MatAutocompleteModule,
    MatMenuModule,
    MatListModule,
    FormsModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatDialogModule,   
    FlexLayoutModule,
    RouterModule,
    HighchartsChartModule,
    HttpClientModule
  ],
  exports:[
    HeaderComponent,
    FooterComponent,
    SidebarComponent,
    StockComponent,
    CardComponent,
    AutocompleteComponent,
    DialogComponent,    
    CountdownSnackbarComponent
  ]
})
export class SharedModule { }