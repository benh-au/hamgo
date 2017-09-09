import { Component, OnInit } from '@angular/core';
import { ToastyService, ToastOptions } from 'ng2-toasty';
import { Observable } from 'rxjs/Rx';

import { Message } from '../model/message';
import { ApiService } from '../services/api.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  public msgs: Message[];
  public loading: boolean = true;

  public callsign: string = "";
  public ip: string = "";
  private message: string = "";
  public sequence: number = 0;

  public showSend: boolean = false;

  constructor(private apiService: ApiService, private toastyService: ToastyService) { }

  update() {
    this.apiService.getMessages()
      .subscribe((msgs) => {
        this.msgs = msgs;
        this.loading = false;

        console.log(msgs);
      });
  }

  ngOnInit() {
    this.update();

    return Observable
      .interval(5000)
      .subscribe(() => {
        this.update();
      });
  }

  send() {
    console.log("sending message");

    var msg: Message = {
      sequence: this.sequence,
      contact: {
        callsign: this.callsign,
        type: 0,
        ips: [this.ip],
      },
      message: this.message
    };

    this.sequence++;
    this.showSend = false;

    this.apiService.spreadCQ(msg)
      .subscribe(() => {
        this.toastyService.info("Message sent!");
      });
  }

}
