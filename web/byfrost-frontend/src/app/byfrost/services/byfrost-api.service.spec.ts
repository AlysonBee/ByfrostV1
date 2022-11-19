import { TestBed } from '@angular/core/testing';

import { ByfrostApiService } from './byfrost-api.service';

describe('ByfrostApiService', () => {
  let service: ByfrostApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ByfrostApiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
