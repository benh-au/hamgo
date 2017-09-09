import { Component, Input, Inject } from '@angular/core';
import { NavEntries } from './navEntry';
import { APP_CONFIG } from '../config/app-config';
import { AppConfig } from '../config/config.interfaces';

@Component({
  selector: 'app-top-nav',
  templateUrl: './top-nav.component.html',
  styleUrls: ['./top-nav.component.css']
})
export class TopNavComponent {
  @Input() routes?: NavEntries;
  title: string;

  constructor( @Inject(APP_CONFIG) private config: AppConfig) {
    this.title = config.title;
  }
}
