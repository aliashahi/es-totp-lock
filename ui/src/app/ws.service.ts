import { Injectable } from '@angular/core';
import { Observable, Observer } from 'rxjs';

@Injectable()
export class WsService {
  private readonly _url: string = `ws://${location.host}`;
  log(): Observable<any> {
    let ws = new WebSocket(this._url + '/api/ws-logs');

    const observable = new Observable((obs: Observer<MessageEvent>) => {
      ws.onmessage = obs.next.bind(obs);
      ws.onerror = obs.error.bind(obs);
      ws.onclose = obs.complete.bind(obs);

      return ws.close.bind(ws);
    });

    return observable;
  }

  user(): Observable<any> {
    let ws = new WebSocket(this._url + '/api/ws-users');

    const observable = new Observable((obs: Observer<MessageEvent>) => {
      ws.onmessage = obs.next.bind(obs);
      ws.onerror = obs.error.bind(obs);
      ws.onclose = obs.complete.bind(obs);

      return ws.close.bind(ws);
    });

    return observable;
  }
}
