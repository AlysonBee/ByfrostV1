import { TestBed } from '@angular/core/testing';

import { ByfrostTemplateService } from './byfrost-template.service';

describe('ByfrostTemplateService', () => {
  let service: ByfrostTemplateService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ByfrostTemplateService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
