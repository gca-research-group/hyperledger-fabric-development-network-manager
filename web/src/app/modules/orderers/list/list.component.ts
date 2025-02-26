import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { TableComponent } from './../../../components/table/table.component';
import { BreadcrumbService } from '@app/services/breadcrumb';
import { OrderersService } from '../services/orderers.service';
import { Orderer } from '@app/models';

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
    field: 'domain',
    label: 'domain',
  },
  {
    field: 'port',
    label: 'port',
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
    label: 'orderers',
    url: '/orderers',
    active: false,
  },
];

@Component({
  selector: 'app-orderers-list',
  imports: [TableComponent],
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
})
export class ListComponent implements OnInit, OnDestroy {
  columns = COLUMNS;

  displayedColumns = COLUMNS.map(column => column.field);

  private breadcrumbService = inject(BreadcrumbService);
  private service = inject(OrderersService);
  data: Orderer[] = [];

  loading = false;
  hasMore = true;
  page = 1;

  constructor() {
    this.breadcrumbService.update(BREADCRUMB);
  }

  ngOnInit(): void {
    this.findAll(this.page);
  }

  ngOnDestroy(): void {
    this.breadcrumbService.reset();
  }

  scroll() {
    if (this.hasMore) {
      this.page++;
      this.findAll(this.page);
    }
  }

  findAll(page: number) {
    this.service.findAll(page).subscribe({
      next: response => {
        this.data = [...this.data, ...response.data];
        this.hasMore = response.hasMore;
      },
      error: error => {
        console.log('[error]', error);
      },
    });
  }
}
