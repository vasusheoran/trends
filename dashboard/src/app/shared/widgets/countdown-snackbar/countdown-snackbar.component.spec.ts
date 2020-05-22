import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CountdownSnackbarComponent } from './countdown-snackbar.component';

describe('CountdownSnackbarComponent', () => {
  let component: CountdownSnackbarComponent;
  let fixture: ComponentFixture<CountdownSnackbarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CountdownSnackbarComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CountdownSnackbarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
