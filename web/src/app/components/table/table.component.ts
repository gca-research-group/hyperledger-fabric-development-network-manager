import { Component, inject, input, OnDestroy, output } from '@angular/core';
import { MatTableModule } from '@angular/material/table';

import { InfiniteScrollDirective } from '@app/directives/infinite-scroll';
import { TranslateModule } from '@ngx-translate/core';
import { Column } from '@app/models';
import { debounceTime, Subject, takeUntil } from 'rxjs';
import { NgStyle, NgTemplateOutlet } from '@angular/common';
import { CustomDatePipe } from '@app/pipes';
import { LanguageService } from '@app/services/language';
import { MatSortModule, Sort } from '@angular/material/sort';

@Component({
  selector: 'app-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss'],
  imports: [
    NgStyle,
    NgTemplateOutlet,
    CustomDatePipe,

    MatTableModule,
    MatSortModule,
    TranslateModule,

    InfiniteScrollDirective,
  ],
})
export class TableComponent<T> implements OnDestroy {
  dataSource = input<T[]>([]);

  displayedColumns = input<string[]>([]);

  columns = input<Column[]>([]);
  height = input<string>();

  loadMore = output();
  sort = output<Sort>();

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

  sortFn(sort: Sort) {
    this.sort.emit(sort);
  }
}
