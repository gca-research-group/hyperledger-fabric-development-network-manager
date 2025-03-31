import { ComponentFixture, TestBed } from '@angular/core/testing';

import { appConfig } from '@app/__tests__/app.config';

import { OrdererSelectorComponent } from './orderer-selector.component';

describe('OrdererSelectorComponent', () => {
  let component: OrdererSelectorComponent;
  let fixture: ComponentFixture<OrdererSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OrdererSelectorComponent],
      providers: [appConfig.providers],
    }).compileComponents();

    fixture = TestBed.createComponent(OrdererSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
