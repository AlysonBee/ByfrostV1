import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ByfrostInfoComponent } from './byfrost-info.component';

describe('ByfrostInfoComponent', () => {
  let component: ByfrostInfoComponent;
  let fixture: ComponentFixture<ByfrostInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ByfrostInfoComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ByfrostInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
