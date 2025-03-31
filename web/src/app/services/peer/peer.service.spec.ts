import { TestBed } from '@angular/core/testing';

import { appConfig } from '@app/__tests__/app.config';

import { PeerService } from './peer.service';

describe('PeerService', () => {
  let service: PeerService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [appConfig.providers],
    });
    service = TestBed.inject(PeerService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
