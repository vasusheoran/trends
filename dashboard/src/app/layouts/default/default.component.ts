import { Component, OnInit } from '@angular/core';
import { WebSocketsService } from 'src/app/shared/services/web-sockets.service';

@Component({
  selector: 'app-default',
  templateUrl: './default.component.html',
  styleUrls: ['./default.component.css']
})
export class DefaultComponent implements OnInit {

  sideBarOpen = true;
  
  constructor(private _socket:WebSocketsService) { }

  ngOnInit(): void {
    // this._socket.emit('message', "hi");
  }

  sideBarToggler(){
    this.sideBarOpen = !this.sideBarOpen;
    setTimeout(() => {
        window.dispatchEvent(
            new Event('resize')
        );
    }, 300);
  }
}
