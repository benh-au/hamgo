import { RouterModule, Routes } from '@angular/router';
import { BreadcrumbService } from 'ng2-breadcrumb/ng2-breadcrumb';

import { HomeComponent } from '../home/home.component';

export const appRoutes: Routes = [
  {
    path: '',
    component: HomeComponent
  },
];

export function registerBreadcrumbAliases(bs: BreadcrumbService) {
//  bs.addFriendlyNameForRoute('/subnets', 'Subnets');
}
