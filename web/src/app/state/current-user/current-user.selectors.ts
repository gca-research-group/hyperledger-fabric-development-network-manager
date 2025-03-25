import { createSelector } from '@ngxs/store';

import { User } from '@app/models';

import { UserState } from './current-user.state';

export const currentUserSelector = createSelector(
  [UserState],
  (state: User) => {
    return state;
  },
);
