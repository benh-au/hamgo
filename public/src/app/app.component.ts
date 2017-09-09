import { Component } from '@angular/core';
import { NavEntries } from './top-nav/navEntry';
import { navEntries } from './config/navigation';
import { BreadcrumbService } from 'ng2-breadcrumb/ng2-breadcrumb';

import {ToastyService, ToastyConfig, ToastOptions, ToastData} from 'ng2-toasty';

import { registerBreadcrumbAliases } from './config/routes';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  navEntries: NavEntries = navEntries;

  constructor(private breadcrumbService: BreadcrumbService, private toastyConfig: ToastyConfig) {
    this.toastyConfig.theme = 'bootstrap';
    registerBreadcrumbAliases(breadcrumbService);
  }
}
