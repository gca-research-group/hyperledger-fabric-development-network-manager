import { DatePipe } from '@angular/common';
import { TestBed } from '@angular/core/testing';

import { CustomDatePipe } from './date.pipe';

describe('CustomDatePipe', () => {
  let pipe: CustomDatePipe;
  beforeAll(() => {
    TestBed.configureTestingModule({
      providers: [CustomDatePipe, DatePipe],
    });

    pipe = TestBed.inject(CustomDatePipe);
  });

  it('create an instance', () => {
    expect(pipe).toBeTruthy();
  });
});
