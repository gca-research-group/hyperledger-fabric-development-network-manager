import { Component, computed, input, output } from '@angular/core';
import { MatTableModule } from '@angular/material/table';

import { InfiniteScrollDirective } from '@app/directives/infinite-scroll';
import { TranslateModule } from '@ngx-translate/core';
import { Column } from '@app/models';
import { IconButtonComponent } from '../icon-button';
import { RouterLink } from '@angular/router';
import { debounceTime, Subject } from 'rxjs';

@Component({
  selector: 'app-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss'],
  imports: [
    MatTableModule,
    TranslateModule,
    RouterLink,
    InfiniteScrollDirective,
    IconButtonComponent,
  ],
})
export class TableComponent<T> {
  dataSource = input<T[]>([]);

  displayedColumns = input<string[]>([]);
  _displayedColumns = computed(() => [
    ...(this.displayedColumns() ?? []),
    'add',
  ]);

  columns = input<Column[]>([]);

  loadMore = output();

  scrollEvents = new Subject<void>();

  constructor() {
    this.scrollEvents.pipe(debounceTime(300)).subscribe(() => {
      this.loadMore.emit();
    });
  }
}
