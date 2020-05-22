import { Injectable } from '@angular/core';
import * as io from 'socket.io-client';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WebSocketsService {
 
  constructor() { 
    this.uri = environment.apiUrl;
    this.socket = io(this.uri);
  }

  socket:any;
  readonly uri:string = "";

  
  public listen(eventName:string){
    return new Observable(sub => {
      this.socket.on(eventName, (message) => {
        sub.next(message);
      });
    });
  }

  public emit(eventName: string, data:any){
    this.socket.emit(eventName, data)
  }
}
