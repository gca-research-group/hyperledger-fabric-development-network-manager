import { Component, inject, OnDestroy } from '@angular/core';
import { TableComponent } from './../../../components/table/table.component';
import { BreadcrumbService } from '@app/services/breadcrumb';

const COLUMNS = [
  {
    id: 'id',
    label: 'id',
  },
  {
    id: 'name',
    label: 'name',
  },
  {
    id: 'createdAt',
    label: 'createdAt',
  },
  {
    id: 'updatedAt',
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
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
  imports: [TableComponent],
})
export class ListComponent implements OnDestroy {
  columns = COLUMNS;

  displayedColumns = COLUMNS.map(column => column.id);

  breadcrumbService = inject(BreadcrumbService);

  constructor() {
    this.breadcrumbService.update(BREADCRUMB);
  }

  ngOnDestroy(): void {
    this.breadcrumbService.reset();
  }
}
