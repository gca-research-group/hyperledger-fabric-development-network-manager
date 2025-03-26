import { TestBed } from '@angular/core/testing';

import { PeersService } from './peers.service';

describe('PeersService', () => {
  let service: PeersService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(PeersService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
