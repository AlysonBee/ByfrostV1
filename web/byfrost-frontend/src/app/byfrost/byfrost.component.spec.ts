import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ByfrostComponent } from './byfrost.component';

describe('ByfrostComponent', () => {
  let component: ByfrostComponent;
  let fixture: ComponentFixture<ByfrostComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ByfrostComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ByfrostComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
