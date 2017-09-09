import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { RouterModule } from '@angular/router';

import { ToastyModule } from 'ng2-toasty';

import { APP_CONFIG, NS_APP_CONFIG } from './config/app-config';
import { appRoutes } from './config/routes';

import { AlertModule, ProgressbarModule, ButtonsModule, ModalModule } from 'ng2-bootstrap';

import { AppComponent } from './app.component';
import { TopNavComponent } from './top-nav/top-nav.component';
import { ApiService } from './services/api.service';
import { ProgressComponent } from './progress/progress.component';
import { LoadingComponent } from './loading/loading.component';

import { Ng2BreadcrumbModule } from 'ng2-breadcrumb/ng2-breadcrumb';
import { HomeComponent } from './home/home.component';

@NgModule({
  declarations: [
    AppComponent,
    TopNavComponent,
    ProgressComponent,
    LoadingComponent,
    HomeComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    AlertModule.forRoot(),
    ProgressbarModule.forRoot(),
    ButtonsModule.forRoot(),
    ModalModule.forRoot(),
    RouterModule.forRoot(appRoutes),
    Ng2BreadcrumbModule.forRoot(),
    ToastyModule.forRoot()
  ],
  providers: [
    { provide: APP_CONFIG, useValue: NS_APP_CONFIG },
    ApiService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
