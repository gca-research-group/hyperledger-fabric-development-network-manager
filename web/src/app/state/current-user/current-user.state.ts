import { State, Action, StateContext } from '@ngxs/store';

import { Injectable } from '@angular/core';

import { User } from '@app/models';

import {
  AddCurrentUserAction,
  RemoveCurrentUserAction,
  UpdateCurrentUserAccessTokenAction,
} from './current-user.actions';

@State<User>({
  name: 'user',
})
@Injectable()
export class UserState {
  @Action(AddCurrentUserAction)
  add(ctx: StateContext<User>, { payload }: AddCurrentUserAction) {
    let state = ctx.getState();
    state = { ...state, ...payload };
    ctx.setState(state);
  }

  @Action(RemoveCurrentUserAction)
  remove(ctx: StateContext<User | object>) {
    let state = ctx.getState();
    state = {};
    ctx.setState(state);
  }

  @Action(UpdateCurrentUserAccessTokenAction)
  updateUserAccessToken(
    ctx: StateContext<User>,
    { payload }: UpdateCurrentUserAccessTokenAction,
  ) {
    let state = ctx.getState();
    state = { ...state, accessToken: payload };
    ctx.setState(state);
  }
}
