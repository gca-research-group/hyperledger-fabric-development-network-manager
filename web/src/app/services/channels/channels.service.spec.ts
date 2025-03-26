import { TestBed } from '@angular/core/testing';

import { appConfig } from '@app/__tests__/app.config';

import { ChannelsService } from './channels.service';

describe('ChannelsService', () => {
  let service: ChannelsService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [appConfig.providers],
    });
    service = TestBed.inject(ChannelsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
