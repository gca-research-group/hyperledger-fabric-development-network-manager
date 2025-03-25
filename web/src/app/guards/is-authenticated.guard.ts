import { inject } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  CanActivateFn,
  Router,
  RouterStateSnapshot,
} from '@angular/router';

import { CurrentUserService } from '@app/services/current-user';

export const isAuthenticatedGuard: CanActivateFn = (
  _route: ActivatedRouteSnapshot,
  _state: RouterStateSnapshot,
) => {
  const router = inject(Router);

  const currentUserService = inject(CurrentUserService);
  const currentUser$ = currentUserService.currentUser$();
  const isAuthenticated = !!currentUser$?.isAuthenticated;

  if (!isAuthenticated) {
    void router.navigate(['/login']);
  }

  return isAuthenticated;
};
