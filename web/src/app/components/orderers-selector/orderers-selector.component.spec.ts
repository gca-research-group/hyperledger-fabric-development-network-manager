import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OrderersSelectorComponent } from './orderers-selector.component';

describe('OrderersSelectorComponent', () => {
  let component: OrderersSelectorComponent;
  let fixture: ComponentFixture<OrderersSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OrderersSelectorComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(OrderersSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
