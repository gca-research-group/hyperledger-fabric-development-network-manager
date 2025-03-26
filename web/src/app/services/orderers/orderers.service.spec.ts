import { TestBed } from '@angular/core/testing';

import { OrderersService } from './orderers.service';

describe('OrderersService', () => {
  let service: OrderersService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(OrderersService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
