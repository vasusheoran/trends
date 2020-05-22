import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-widget-card',
  templateUrl: './card.component.html',
  styleUrls: ['./card.component.css']
})
export class CardComponent implements OnInit {

  @Input() label:string;
  @Input() value:number;
  @Input() percentage:string;
  @Input() bg:string;

  constructor() { }

  ngOnInit(): void {
  }

}
