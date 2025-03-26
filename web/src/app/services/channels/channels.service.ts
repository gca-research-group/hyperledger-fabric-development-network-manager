import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';

import { Channel } from '@app/models';

import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class ChannelsService {
  private readonly http = inject(HttpClient);
  private readonly url = `${environment.apiUrl}/channel/`;

  findAll(params?: object) {
    return this.http.get<{ data: Channel[]; hasMore: boolean }>(this.url, {
      params: { ...params },
    });
  }

  findById(id: number) {
    return this.http.get<Channel>(`${this.url}${id}`);
  }

  delete(id: number) {
    return this.http.delete(`${this.url}${id}`);
  }

  save(channel: Channel) {
    if (channel.id) {
      return this.http.put(`${this.url}`, channel);
    }

    return this.http.post(this.url, channel);
  }
}
