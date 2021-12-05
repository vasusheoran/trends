import { Injectable, EventEmitter } from '@angular/core';

@Injectable()
export class WebSocketsService {

  private socket: WebSocket;
  private listener: EventEmitter<any> = new EventEmitter();

  public constructor() { }

  public enable(uri: string) {
    this.socket = new WebSocket(uri);
    this.socket.onopen = event => {
      this.listener.emit({ "type": "open", "data": event });
    }
    this.socket.onclose = event => {
      this.listener.emit({ "type": "close", "data": event });
    }
    this.socket.onmessage = event => {
      this.listener.emit({ "type": "message", "data": JSON.parse(event.data) });
    }
  }

  public send(data: string) {
    if (this.socket == null || this.socket == undefined) {
      console.error("socket is not enabled");
      return
    }
    this.socket.send(data);
  }

  public close() {
    if (this.socket == null || this.socket == undefined) {
      console.error("socket is not enabled");
      return
    }
    this.socket.close();
  }

  public getEventListener() {
    if (this.socket == null || this.socket == undefined) {
      console.error("socket is not enabled");
      return
    }
    return this.listener;
  }

}
