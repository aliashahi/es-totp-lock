import { Component } from '@angular/core';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-log',
  templateUrl: './log.component.html',
  styleUrls: ['./log.component.scss'],
})
export class LogComponent {
  constructor(public api: ApiService) {}
}
