import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { DefaultComponent } from './layouts/default/default.component';
import { DashboardComponent } from './modules/dashboard/dashboard.component';
import { HistoricalDataComponent } from './modules/historical-data/historical-data.component';
import { SettingsComponent } from './modules/settings/settings.component';


const routes: Routes = [{
  path : '',
  component: DefaultComponent,
  children:[{
    path:'',
    component: SettingsComponent
  },{
    path: 'historica-data',
    component: HistoricalDataComponent
  },{
    path: 'dash',
    component: DashboardComponent
  }]
}];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
