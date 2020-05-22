import { TestBed } from '@angular/core/testing';

import { ConfigLogService } from './config-log.service';

describe('ConfigLogService', () => {
  let service: ConfigLogService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ConfigLogService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
