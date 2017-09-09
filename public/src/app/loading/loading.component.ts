import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-loading',
  templateUrl: './loading.component.html',
  styleUrls: ['./loading.component.css']
})
export class LoadingComponent {

  @Input('state') public loading: boolean = false;
  @Input('error') public error: string = null;

  constructor() { }
}
