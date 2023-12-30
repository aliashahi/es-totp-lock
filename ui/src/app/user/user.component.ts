import { Component } from '@angular/core';
import { ApiService } from '../api.service';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss'],
})
export class UserComponent {
  new_username: string = '';
  new_password: string = '';

  constructor(
    public api: ApiService,
    public dialogRef: MatDialogRef<UserComponent>
  ) {}

  onSubmit() {
    this.api.create(this.new_username, this.new_password, () => {
      this.dialogRef.close();
    });
  }
}
