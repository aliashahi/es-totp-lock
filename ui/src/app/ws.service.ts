import { Injectable } from '@angular/core';
import { Observable, Observer, Subject, filter } from 'rxjs';

@Injectable()
export class WsService {
  log(): Observable<any> {
    let ws = new WebSocket('ws://' + location.host + '/api/ws-logs');

    const observable = new Observable((obs: Observer<MessageEvent>) => {
      ws.onmessage = obs.next.bind(obs);
      ws.onerror = obs.error.bind(obs);
      ws.onclose = obs.complete.bind(obs);

      return ws.close.bind(ws);
    });

    // const observer = {
    //   next: (data: Object) => {
    //     if (ws.readyState === WebSocket.OPEN) {
    //       ws.send(JSON.stringify(data));
    //     }
    //   },
    // };

    return observable;
  }

  user(): Observable<any> {
    let ws = new WebSocket('ws://' + location.host + '/api/ws-users');

    const observable = new Observable((obs: Observer<MessageEvent>) => {
      ws.onmessage = obs.next.bind(obs);
      ws.onerror = obs.error.bind(obs);
      ws.onclose = obs.complete.bind(obs);

      return ws.close.bind(ws);
    });

    return observable;
  }
}
