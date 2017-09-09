import { Component, OnInit, OnChanges, Input } from '@angular/core';

@Component({
  selector: 'app-progress',
  templateUrl: './progress.component.html',
  styleUrls: ['./progress.component.css']
})
export class ProgressComponent implements OnInit, OnChanges {

  @Input() min: number = 0;
  @Input() max: number = 100;
  @Input() value: number = 25;

  @Input() type: string = 'success';

  constructor() { }

  ngOnChanges() {
    let percent = (this.value - this.min) / this.max;

    if (percent > 0.9) {
      this.type = 'danger';
    } else if (percent > 0.7) {
      this.type = 'warning';
    } else if (percent > 0.5) {
      this.type = 'info';
    } else {
      this.type = 'success';
    }
  }

  ngOnInit() {
  }

}
