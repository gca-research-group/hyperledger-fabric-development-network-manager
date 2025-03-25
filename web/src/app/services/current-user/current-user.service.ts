import { Store } from '@ngxs/store';

import { inject, Injectable, Signal } from '@angular/core';

import { User } from '@app/models';
import {
  AddCurrentUserAction,
  currentUserSelector,
  RemoveCurrentUserAction,
  UpdateCurrentUserAccessTokenAction,
} from '@app/state/current-user';

@Injectable({
  providedIn: 'root',
})
export class CurrentUserService {
  private store = inject(Store);

  currentUser$: Signal<User | undefined> =
    this.store.selectSignal(currentUserSelector);

  add(user: User) {
    this.store.dispatch(new AddCurrentUserAction(user));
  }

  remove() {
    this.store.dispatch(new RemoveCurrentUserAction());
  }

  updateAccessToken(accessToken: string) {
    this.store.dispatch(new UpdateCurrentUserAccessTokenAction(accessToken));
  }
}
