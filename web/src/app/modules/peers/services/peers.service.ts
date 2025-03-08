import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';

import { Peer } from '@app/models';

import { environment } from '../../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class PeersService {
  private readonly http = inject(HttpClient);
  private readonly url = `${environment.apiUrl}/peer/`;

  findAll(params?: object) {
    return this.http.get<{ data: Peer[]; hasMore: boolean }>(this.url, {
      params: { ...params },
    });
  }

  findById(id: number) {
    return this.http.get(`${this.url}${id}`);
  }

  delete(id: number) {
    return this.http.delete(`${this.url}${id}`);
  }

  save(peer: Peer) {
    if (peer.id) {
      return this.http.put(`${this.url}`, peer);
    }

    return this.http.post(this.url, peer);
  }
}
