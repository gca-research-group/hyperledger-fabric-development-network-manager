import { TestBed } from '@angular/core/testing';

import { appConfig } from '@app/__tests__/app.config';

import { OrdererService } from './orderer.service';

describe('OrdererService', () => {
  let service: OrdererService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [appConfig.providers],
    });
    service = TestBed.inject(OrdererService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
