import { InjectionToken } from '@angular/core';
import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';

export const IS_MOBILE = new InjectionToken<boolean>('IsMobile', {
  providedIn: 'root',
  factory: () => {
    const breakpointObserver = new BreakpointObserver();
    return breakpointObserver.isMatched(Breakpoints.Handset);
  },
});
