import { Injectable } from '@angular/core';
import { Socket } from 'ngx-socket-io';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class WsService {
  constructor(private socket: Socket) {
    this.getMessage().subscribe((data) => console.log(data));
  }

  sendMessage(msg: string) {
    this.socket.emit('code', msg);
  }
  getMessage() {
    return this.socket.fromEvent('code').pipe(map((data: any) => data.code));
  }
}
