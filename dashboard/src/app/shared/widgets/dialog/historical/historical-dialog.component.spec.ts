import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HistoricalDialogComponent } from './historical-dialog.component';

describe('HistoricalDialogComponent', () => {
  let component: HistoricalDialogComponent;
  let fixture: ComponentFixture<HistoricalDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HistoricalDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HistoricalDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
