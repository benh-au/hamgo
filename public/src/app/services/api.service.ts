import { Injectable, Inject } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from "@angular/http";
import { Observable } from 'rxjs/Rx';

import { APP_CONFIG } from '../config/app-config';
import { AppConfig } from '../config/config.interfaces';

import { Message } from '../model/message';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

@Injectable()
export class ApiService {

  constructor(private http: Http,
    @Inject(APP_CONFIG) private config: AppConfig) {
  }

  getMessages(): Observable<Message[]> {
    return this.http.get(this.config.apiEndpoint + '/cache')
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error));
  }

  spreadCQ(msg: Message): Observable<void> {
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let options = new RequestOptions({ headers: headers });

    return this.http.post(this.config.apiEndpoint + '/spread/cq', msg, options)
      .map((res: Response) => null)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

}
