import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { LoginComponent } from './login/login.component';
import { CodeComponent } from './code/code.component';
import { UserComponent } from './user/user.component';
import { CommonModule } from '@angular/common';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MAT_FORM_FIELD_DEFAULT_OPTIONS } from '@angular/material/form-field';
import { HttpClientModule } from '@angular/common/http';
import { LogComponent } from './log/log.component';
import { SocketIoConfig, SocketIoModule } from 'ngx-socket-io';
import { WsService } from './ws.service';

const MATERIALS = [
  MatInputModule,
  MatButtonModule,
  MatProgressSpinnerModule,
  MatIconModule,
];

const config: SocketIoConfig = {
  url: '',
  options: {
    path: '/api/ws',
    autoConnect: true,
    upgrade: true,
    rememberUpgrade: true,
  },
};

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    CodeComponent,
    UserComponent,
    LogComponent,
  ],
  imports: [
    BrowserModule,
    // AppRoutingModule,
    BrowserAnimationsModule,
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    HttpClientModule,
    SocketIoModule.forRoot(config),
    ...MATERIALS,
  ],
  providers: [
    {
      provide: MAT_FORM_FIELD_DEFAULT_OPTIONS,
      useValue: { appearance: 'outline' },
    },
    WsService,
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
