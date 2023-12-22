import { Component } from '@angular/core';
import { ApiService } from '../api.service';
import { WsService } from '../ws.service';

@Component({
  selector: 'app-code',
  templateUrl: './code.component.html',
  styleUrls: ['./code.component.scss'],
})
export class CodeComponent {
  current_date: string = '';

  constructor(public api: ApiService) {
    this.update();
    setInterval(this.update.bind(this), 1000);
    api.code();
  }

  private update() {
    let d = new Date();
    this.current_date = `${d.getHours()}:${d.getMinutes()}:${d.getSeconds()}`;
  }
}
