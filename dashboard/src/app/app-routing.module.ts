import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { DefaultComponent } from './layouts/default/default.component';
import { DashboardComponent } from './modules/dashboard/dashboard.component';
import { HistoricalDataComponent } from './modules/historical-data/historical-data.component';


const routes: Routes = [{
  path : '',
  component: DefaultComponent,
  children:[{
    path:'',
    component: DashboardComponent
  },{
    path: 'historica-data',
    component: HistoricalDataComponent
  }]
}];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
