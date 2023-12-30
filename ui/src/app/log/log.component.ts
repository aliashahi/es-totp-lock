import { Component } from '@angular/core';
import { ApiService } from '../api.service';
import { MatDialog } from '@angular/material/dialog';
import { UserComponent } from '../user/user.component';

@Component({
  selector: 'app-log',
  templateUrl: './log.component.html',
  styleUrls: ['./log.component.scss'],
})
export class LogComponent {
  constructor(public api: ApiService, private dialog: MatDialog) {}

  onNewUser() {
    this.dialog.open(UserComponent, {
      minWidth: '33vw',
    });
  }
}
