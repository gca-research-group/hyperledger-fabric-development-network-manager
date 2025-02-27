import {
  Component,
  computed,
  inject,
  input,
  OnDestroy,
  output,
} from '@angular/core';
import { MatTableModule } from '@angular/material/table';

import { InfiniteScrollDirective } from '@app/directives/infinite-scroll';
import { TranslateModule } from '@ngx-translate/core';
import { Column } from '@app/models';
import { IconButtonComponent } from '../icon-button';
import { RouterLink } from '@angular/router';
import { debounceTime, Subject, takeUntil } from 'rxjs';
import { NgStyle } from '@angular/common';
import { CustomDatePipe } from '@app/pipes';
import { LanguageService } from '@app/services/language';

@Component({
  selector: 'app-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss'],
  imports: [
    RouterLink,
    NgStyle,
    CustomDatePipe,

    MatTableModule,
    TranslateModule,

    InfiniteScrollDirective,
    IconButtonComponent,
  ],
})
export class TableComponent<T> implements OnDestroy {
  dataSource = input<T[]>([]);

  displayedColumns = input<string[]>([]);
  _displayedColumns = computed(() => [
    ...(this.displayedColumns() ?? []),
    'add',
  ]);

  columns = input<Column[]>([]);
  height = input<string>();

  loadMore = output();

  scrollEvents = new Subject<void>();

  language = 'en';
  private languageService = inject(LanguageService);

  onDestroy$ = new Subject();

  constructor() {
    this.scrollEvents
      .pipe(debounceTime(300), takeUntil(this.onDestroy$))
      .subscribe(() => {
        this.loadMore.emit();
      });

    this.languageService.language$.subscribe(language => {
      this.language = language;
    });
  }

  ngOnDestroy(): void {
    this.onDestroy$.complete();
    this.onDestroy$.unsubscribe();
  }
}
