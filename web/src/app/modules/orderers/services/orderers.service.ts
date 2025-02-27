import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Orderer } from '@app/models';
import { environment } from '../../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class OrderersService {
  private readonly http = inject(HttpClient);
  private readonly url = `${environment.apiUrl}/orderer/`;

  findAll(params: object) {
    return this.http.get<{ data: Orderer[]; hasMore: boolean }>(this.url, {
      params: { ...params },
    });
  }

  delete(id: number) {
    return this.http.delete(`${this.url}${id}`);
  }

  save(orderer: Orderer) {
    return this.http.post(this.url, orderer);
  }
}
