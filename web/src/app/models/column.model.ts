export enum ColumnType {
  DATETIME = 'DATETIME',
}
export interface Column {
  label: string;
  field: string;
  type?: ColumnType;
}
