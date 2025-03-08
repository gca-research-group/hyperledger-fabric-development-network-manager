import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { InjectionToken } from '@angular/core';

export const IS_MOBILE = new InjectionToken<boolean>('IsMobile', {
  providedIn: 'root',
  factory: () => {
    const breakpointObserver = new BreakpointObserver();
    return breakpointObserver.isMatched(Breakpoints.Handset);
  },
});
