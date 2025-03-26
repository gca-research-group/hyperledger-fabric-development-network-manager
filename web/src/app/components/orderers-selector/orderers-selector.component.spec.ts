import { ComponentFixture, TestBed } from '@angular/core/testing';

import { appConfig } from '@app/__tests__/app.config';

import { OrderersSelectorComponent } from './orderers-selector.component';

describe('OrderersSelectorComponent', () => {
  let component: OrderersSelectorComponent;
  let fixture: ComponentFixture<OrderersSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OrderersSelectorComponent],
      providers: [appConfig.providers],
    }).compileComponents();

    fixture = TestBed.createComponent(OrderersSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
