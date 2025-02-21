import { TestBed } from '@angular/core/testing';
import { CustomControlValueAccessorDirective } from './custom-control-value-accessor.directive';

describe('CustomControlValueAccessorDirective', () => {
  let directive: CustomControlValueAccessorDirective;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [CustomControlValueAccessorDirective],
    });

    directive = TestBed.inject(CustomControlValueAccessorDirective);
  });

  it('should create an instance', () => {
    expect(directive).toBeTruthy();
  });
});
