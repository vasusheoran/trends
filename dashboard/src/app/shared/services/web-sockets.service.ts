import { Injectable } from '@angular/core';
import * as io from 'socket.io-client';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WebSocketsService {


  socket: any;
  readonly uri: string = "";

  constructor() {
    this.uri = environment.apiUrl;
  }

  public enable() {
    this.socket = io(this.uri);
  }

  public listen(eventName: string) {
    return new Observable(sub => {
      this.socket.on(eventName, (message) => {
        sub.next(message);
      });
    });
  }

  public emit(eventName: string, data: any) {
    this.socket.emit(eventName, data)
  }

  public disconnet() {
    this.socket.disconnect();
  }

  public connect() {
    this.socket.connect();
  }
}
