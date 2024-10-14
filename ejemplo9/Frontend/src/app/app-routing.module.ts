import { Component, NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { InicioComponent } from './component/inicio/inicio.component';

const routes: Routes = [
  {
    path: 'inicio',
    component: InicioComponent
  },
  {
    path: '**',
    redirectTo: 'inicio'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
