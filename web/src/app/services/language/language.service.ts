import { BehaviorSubject } from 'rxjs';

import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class LanguageService {
  private _language = new BehaviorSubject(
    localStorage.getItem('language') ?? 'en',
  );

  language$ = this._language.asObservable();
  language = this._language.value;

  update(language: string) {
    localStorage.setItem('language', language);
    this._language.next(language);
  }
}
