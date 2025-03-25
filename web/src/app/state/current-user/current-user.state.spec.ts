import { provideStore, Store } from '@ngxs/store';

import { TestBed } from '@angular/core/testing';

import { User } from '@app/models';

import { AddCurrentUserAction } from './current-user.actions';
import { UserState } from './current-user.state';

describe('User store', () => {
  let store: Store;
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideStore([UserState])],
    });

    store = TestBed.inject(Store);
  });

  it('should create an action and add an item', () => {
    const expected: User = {
      name: 'Name',
      email: 'Email',
      isSuper: false,
      isAuthenticated: false,
    };
    store.dispatch(new AddCurrentUserAction(expected));
  });
});
