import { TestBed } from '@angular/core/testing';

import { appConfig } from '@app/__tests__/app.config';

import { PeersService } from './peers.service';

describe('PeersService', () => {
  let service: PeersService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [appConfig.providers],
    });
    service = TestBed.inject(PeersService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
