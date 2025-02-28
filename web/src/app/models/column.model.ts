import { TemplateRef } from '@angular/core';

export enum ColumnType {
  DATETIME = 'DATETIME',
  TEMPLATE = 'TEMPLATE',
}
export interface Column {
  label: string;
  id: string;
  rowType?: ColumnType;
  columnType?: ColumnType;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  templateColumn?: TemplateRef<any>;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  templateRow?: TemplateRef<any>;
  width?: number;
  sortable?: boolean;
}
