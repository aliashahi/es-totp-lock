import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { WsService } from './ws.service';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  Mode: 'LOGIN' | 'CODE' | 'LOG' = 'LOGIN';
  username: string = '';
  password: string = '';
  currentCode: string = '';
  currentSecret: string = '';
  lastLogId: string = '';
  logs: string[] = [];
  users: any[] = [];

  private _loading: number = 0;

  constructor(private ws: WsService, private http: HttpClient) {}

  get loading(): boolean {
    return this._loading != 0;
  }

  userInfo() {
    this.http.get('/api/userinfo').subscribe({
      next: ({ isAdmin, username }: any) => {
        this.Mode = isAdmin ? 'LOG' : 'CODE';
        this.username = username;
        this.code();
        this.allLogs();
        this.allUsers();
      },
      error: () => {
        this.Mode = 'LOGIN';
      },
    });
  }

  login() {
    if (this.Mode != 'LOGIN') return;
    this._loading++;
    this.http
      .post('/api/login', {
        username: this.username,
        password: this.password,
      })
      .subscribe({
        next: ({ isAdmin }: any) => {
          this._loading--;
          this.Mode = isAdmin ? 'LOG' : 'CODE';
          this.userInfo();
        },
        error: () => {
          this._loading--;
        },
      });
  }

  logout() {
    this._loading++;
    this.http.get('/api/logout').subscribe({
      next: () => {
        this._loading--;
        this.Mode = 'LOGIN';
      },
      error: () => {
        this._loading--;
        this.Mode = 'LOGIN';
      },
    });
  }

  code() {
    if (this.Mode != 'CODE') return;
    this.http.get(`/api/code`).subscribe({
      next: ({ code, secret }: any) => {
        this.currentCode = code;
        this.currentSecret = secret;
      },
      error: () => {
        this.Mode = 'LOGIN';
      },
    });
  }

  create(username: string, password: string, onClose: any) {
    this._loading++;
    this.http
      .post('/api/create', {
        username,
        password,
      })
      .subscribe({
        next: () => {
          this._loading--;
          this.Mode = 'LOG';
          onClose();
        },
        error: () => {
          this._loading--;
        },
      });
  }

  allLogs() {
    if (this.Mode != 'LOG') return;
    this.http.get(`/api/logs`).subscribe({
      next: ({ logs }: any) => {
        this.lastLogId = logs[logs.length - 1].id;
        this.logs = [
          ...logs.map((i: any) => i.message).reverse(),
          ...this.logs,
        ];
        this.ws.log().subscribe((data) => {
          this.logs = [data.data, ...this.logs];
        });
      },
      error: () => {
        this.Mode = 'LOGIN';
      },
    });
  }

  allUsers() {
    if (this.Mode != 'LOG') return;
    this.http.get(`/api/users`).subscribe({
      next: ({ users }: any) => {
        this.users = [...users];
        this.ws.user().subscribe((data) => {
          this.users = [JSON.parse(data.data), ...this.users];
        });
      },
      error: () => {
        this.Mode = 'LOGIN';
      },
    });
  }

  deleteUser(id: string) {
    if (this.Mode != 'LOG') return;
    this.http.delete(`/api/delete/${id}`).subscribe({
      next: ({ users }: any) => {
        this.users = [...users];
      },
      error: () => {
        this.Mode = 'LOGIN';
      },
    });
  }
}
