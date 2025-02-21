import { Component, computed, input, output } from '@angular/core';
import { MatTableModule } from '@angular/material/table';

import { InfiniteScrollDirective } from '@app/directives/infinite-scroll';
import { TranslateModule } from '@ngx-translate/core';
import { Column, SmartContract } from '@app/models';
import { IconButtonComponent } from '../icon-button';
import { RouterLink } from '@angular/router';

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
export class TableComponent {
  dataSource = input<SmartContract[]>([]);

  displayedColumns = input<string[]>([]);
  _displayedColumns = computed(() => [
    ...(this.displayedColumns() ?? []),
    'add',
  ]);

  columns = input<Column[]>([]);

  loadMore = output();
}
