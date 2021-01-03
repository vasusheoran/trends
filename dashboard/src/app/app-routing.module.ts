import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { DefaultComponent } from './layouts/default/default.component';
import { DashboardComponent } from './modules/dashboard/dashboard.component';
import { HistoricalDataComponent } from './modules/historical-data/historical-data.component';
// import { SettingsComponent } from './modules/settings/settings.component';
import { SymbolsComponent } from 'src/app/modules/symbols/symbols.component';


const routes: Routes = [{
  path : '',
  component: DefaultComponent,
  children:[{
    path: '',
    component: SymbolsComponent
  },{
    path: 'historical',
    component: HistoricalDataComponent
  },{
    path: 'symbols',
    component: SymbolsComponent
  },{
    path: 'dashboard',
    component: DashboardComponent
  }]
}];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
