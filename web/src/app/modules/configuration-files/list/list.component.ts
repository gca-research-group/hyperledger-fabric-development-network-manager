import { Component, inject, OnDestroy } from '@angular/core';
import { TableComponent } from './../../../components/table/table.component';
import { BreadcrumbService } from '@app/services/breadcrumb';

const COLUMNS = [
  {
    field: 'id',
    label: 'id',
  },
  {
    field: 'name',
    label: 'name',
  },
  {
    field: 'createdAt',
    label: 'createdAt',
  },
  {
    field: 'updatedAt',
    label: 'updatedAt',
  },
];

const BREADCRUMB = [
  {
    label: 'home',
    url: '/',
    active: false,
  },
  {
    label: 'configuration-files',
    url: '/configuration-files',
    active: false,
  },
];

@Component({
  selector: 'app-configuration-files-list',
  imports: [TableComponent],
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
})
export class ListComponent implements OnDestroy {
  columns = COLUMNS;

  displayedColumns = COLUMNS.map(column => column.field);

  breadcrumbService = inject(BreadcrumbService);

  constructor() {
    this.breadcrumbService.update(BREADCRUMB);
  }

  ngOnDestroy(): void {
    this.breadcrumbService.reset();
  }
}
