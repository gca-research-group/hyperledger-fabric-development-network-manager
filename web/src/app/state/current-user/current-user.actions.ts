import { User } from '@app/models';

export class AddCurrentUserAction {
  static readonly type = '[User] Add user';
  constructor(readonly payload: User) {}
}

export class UpdateCurrentUserAccessTokenAction {
  static readonly type = "[User] Add the user's access token";
  constructor(readonly payload: string) {}
}

export class RemoveCurrentUserAction {
  static readonly type = '[User] Remove the current user';
}
