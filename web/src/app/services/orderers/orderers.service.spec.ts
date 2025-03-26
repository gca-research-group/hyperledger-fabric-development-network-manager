import { TestBed } from '@angular/core/testing';

import { appConfig } from '@app/__tests__/app.config';

import { OrderersService } from './orderers.service';

describe('OrderersService', () => {
  let service: OrderersService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [appConfig.providers],
    });
    service = TestBed.inject(OrderersService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
