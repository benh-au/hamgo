import { OpaqueToken } from '@angular/core';
import { AppConfig } from './config.interfaces';
import { environment } from '../../environments/environment';

export let APP_CONFIG = new OpaqueToken('app.config');

export const NS_APP_CONFIG: AppConfig = {
  apiEndpoint: environment.production ? "/api" : "//127.0.0.1:9125/api",
  title: "HamGO"
};
