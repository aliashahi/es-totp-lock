import { Component } from '@angular/core';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss'],
})
export class UserComponent {
  new_username: string = '';
  new_password: string = '';

  constructor(public api: ApiService) {}
}
