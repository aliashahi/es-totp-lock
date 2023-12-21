import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  Mode: 'LOGIN' | 'CODE' | 'USER' | 'LOG' = 'LOGIN';
  username: string = '';
  password: string = '';
  currentCode: string = '123 123';
  lastLogId: string = '';
  logs: string[] = [];
  users: any[] = [];

  private _loading: number = 0;

  constructor(private http: HttpClient) {
    this.userInfo();
    setInterval(() => this.code(), 1 * 1000);
    setInterval(() => this.allLogs(), 5 * 1000);
    setInterval(() => this.getUsers(), 5 * 1000);
  }

  get loading(): boolean {
    return this._loading != 0;
  }

  userInfo() {
    this.http.get('/api/userinfo').subscribe({
      next: ({ isAdmin }: any) => {
        this.Mode = isAdmin ? 'LOG' : 'CODE';
        this.code();
        this.allLogs();
        this.getUsers();
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
        },
        error: () => {
          this._loading--;
        },
      });
  }

  code() {
    if (this.Mode != 'CODE') return;
    this.http.get(`/api/code`).subscribe({
      next: ({ code }: any) => {
        this.currentCode = code;
      },
      error: () => {
        this.Mode = 'LOGIN';
      },
    });
  }

  create(username: string, password: string) {
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
        },
        error: () => {
          this._loading--;
        },
      });
  }

  allLogs() {
    if (this.Mode != 'LOG') return;
    this.http.get(`/api/logs?lastId=${this.lastLogId}`).subscribe({
      next: ({ logs }: any) => {
        this.lastLogId = logs[logs.length - 1].id;
        this.logs = [
          ...logs.map((i: any) => i.message).reverse(),
          ...this.logs,
        ];
      },
      error: () => {
        this.Mode = 'LOGIN';
      },
    });
  }

  getUsers() {
    if (this.Mode != 'LOG') return;
    this.http.get(`/api/users`).subscribe({
      next: ({ users }: any) => {
        this.users = [...users];
      },
      error: () => {
        this.Mode = 'LOGIN';
      },
    });
  }
}
