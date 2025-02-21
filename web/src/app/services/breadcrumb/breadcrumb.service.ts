import { Injectable } from '@angular/core';
import { Breadcrumb } from '@app/models';
import { BehaviorSubject } from 'rxjs';

const BASE_BREADCRUMB: Breadcrumb[] = [
  {
    label: 'home',
    url: '/',
  },
];

@Injectable({
  providedIn: 'root',
})
export class BreadcrumbService {
  private breadcrumb = new BehaviorSubject<Breadcrumb[]>(BASE_BREADCRUMB);

  breadcrumb$ = this.breadcrumb.asObservable();

  update(breadcrumb: Breadcrumb[]) {
    this.breadcrumb.next(breadcrumb);
  }

  reset() {
    this.breadcrumb.next(BASE_BREADCRUMB);
  }
}
