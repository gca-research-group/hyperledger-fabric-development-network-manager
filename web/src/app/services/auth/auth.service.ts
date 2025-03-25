import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';

import { Login, User } from '@app/models';

import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly apiUrl = environment.apiUrl;
  private readonly httpClient = inject(HttpClient);

  login(data: Login) {
    return this.httpClient.post<User>(`${this.apiUrl}/auth/login`, data);
  }

  refresh() {
    return this.httpClient.post<{ accessToken: string }>(
      `${this.apiUrl}/auth/refresh`,
      null,
    );
  }
}
