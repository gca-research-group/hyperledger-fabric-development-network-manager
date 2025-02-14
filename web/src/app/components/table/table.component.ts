import { Component, EventEmitter, Input, Output } from '@angular/core';
import { MatTableModule } from '@angular/material/table';

import { InfiniteScrollDirective } from '@app/directives/infinite-scroll';
import { TranslateModule } from '@ngx-translate/core';
import { Column, SmartContract } from '@app/models';

@Component({
  selector: 'app-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss'],
  imports: [MatTableModule, InfiniteScrollDirective, TranslateModule],
})
export class TableComponent {
  @Input()
  dataSource: SmartContract[] = [];

  @Input()
  displayedColumns: string[] = [];

  @Input()
  columns: Column[] = [];

  @Output()
  loadMore = new EventEmitter();

  loadData() {
    this.loadMore.emit();
  }
}
