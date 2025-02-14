import { Component } from '@angular/core';
import { TableComponent } from './../../../components/table/table.component';

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

@Component({
  selector: 'app-list',
  imports: [TableComponent],
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
})
export class ListComponent {
  columns = COLUMNS;

  displayedColumns = COLUMNS.map(column => column.field);
}
