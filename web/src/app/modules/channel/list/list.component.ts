import { TranslateModule } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';
import { debounceTime } from 'rxjs';

import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ElementRef,
  inject,
  OnDestroy,
  TemplateRef,
  viewChild,
} from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { Sort } from '@angular/material/sort';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';

import { DeleteDialogComponent } from '@app/components/delete-dialog';
import { IconButtonComponent } from '@app/components/icon-button';
import { InputComponent } from '@app/components/input';
import { TableComponent } from '@app/components/table';
import { Column, ColumnType, Channel, Breadcrumb } from '@app/models';
import { BreadcrumbService } from '@app/services/breadcrumb';
import { ChannelService } from '@app/services/channel';

const COLUMNS: Column[] = [
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
    rowType: ColumnType.DATETIME,
  },
  {
    id: 'updatedAt',
    label: 'updatedAt',
    rowType: ColumnType.DATETIME,
  },
  {
    id: 'actions',
    label: '',
    columnType: ColumnType.TEMPLATE,
    rowType: ColumnType.TEMPLATE,
    width: 150,
    sortable: false,
  },
];

const BREADCRUMB: Breadcrumb[] = [
  {
    label: 'home',
    url: '/',
  },
  {
    label: 'channel',
  },
];

@Component({
  selector: 'app-channel-list',
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
  imports: [
    TableComponent,
    ReactiveFormsModule,
    FormsModule,
    RouterLink,

    TranslateModule,

    InputComponent,
    IconButtonComponent,
  ],
})
export class ListComponent implements AfterViewInit, OnDestroy {
  columns = COLUMNS;

  displayedColumns = COLUMNS.map(column => column.id);

  private breadcrumbService = inject(BreadcrumbService);
  private service = inject(ChannelService);
  private formBuilder = inject(FormBuilder);

  data: Channel[] = [];

  loading = false;
  hasMore = true;

  form!: FormGroup;
  formElement = viewChild<ElementRef<HTMLFormElement>>('filters');

  private filters = viewChild<ElementRef<HTMLFormElement>>('filters');
  tableHeight!: string;
  private cdk = inject(ChangeDetectorRef);

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private actionsColumn = viewChild<TemplateRef<any>>('actionsColumn');
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private actionsRow = viewChild<TemplateRef<any>>('actionsRow');

  readonly dialog = inject(MatDialog);
  private toastr = inject(ToastrService);
  private router = inject(Router);
  private activatedRoute = inject(ActivatedRoute);

  constructor() {
    this.breadcrumbService.update(BREADCRUMB);

    this.form = this.formBuilder.group({
      id: null,
      name: null,
      domain: null,
      page: 1,
      pageSize: 20,
      orderBy: null,
      orderDirection: null,
    });

    const queryParams = this.activatedRoute.snapshot.queryParams;
    this.form.patchValue(queryParams, { emitEvent: false });

    this.form.valueChanges.pipe(debounceTime(300)).subscribe(() => {
      this.search();
      void this.router.navigate([], {
        queryParams: this.form.value as Record<string, string>,
        queryParamsHandling: 'merge',
      });
    });
  }

  ngAfterViewInit(): void {
    const form = this.formElement()?.nativeElement;

    const marginBottom = getComputedStyle(
      this.filters()?.nativeElement as Element,
    ).marginBottom;

    this.tableHeight = `calc(100vh - var(--hfdnm-toolbar-height) - (2 * var(--hfdnm-content-vertical-padding)) - ${form?.offsetHeight}px - ${marginBottom})`;

    this.columns = this.columns.map(column => {
      if (column.id !== 'actions') {
        return column;
      }

      return {
        ...column,
        templateRow: this.actionsRow(),
        templateColumn: this.actionsColumn(),
      };
    });

    this.cdk.detectChanges();
  }

  ngOnDestroy(): void {
    this.breadcrumbService.reset();
  }

  openDialog(item: Channel): void {
    const dialogRef = this.dialog.open(DeleteDialogComponent, {
      data: item.id,
    });

    dialogRef.afterClosed().subscribe((id: number | undefined) => {
      if (id) {
        this.service.delete(id).subscribe({
          next: () => {
            this.data = this.data.filter(item => item.id !== id);
            this.toastr.success('DELETED_SUCCESSFULLY');
          },
          error: error => {
            console.log('[error]', error);
          },
        });
      }
    });
  }

  scroll() {
    if (this.hasMore) {
      this.form.patchValue(
        {
          page: ((this.form.get('page')?.value as unknown as number) ?? 0) + 1,
        },
        { emitEvent: false },
      );
      this.findAll();
    }
  }

  findAll() {
    const _params = this.removeNullFields(this.form.value);
    this.service.findAll(_params).subscribe({
      next: response => {
        this.data = [...this.data, ...response.data];
        this.hasMore = response.hasMore;
      },
      error: error => {
        console.log('[error]', error);
      },
    });
  }

  removeNullFields<T extends object>(obj: T): Partial<T> {
    return Object.fromEntries(
      Object.entries(obj).filter(([_index, value]) => value !== null),
    ) as Partial<T>;
  }

  search() {
    const _params = this.removeNullFields(this.form.value);
    this.service.findAll(_params).subscribe({
      next: response => {
        this.data = response.data;
        this.hasMore = response.hasMore;
      },
      error: error => {
        console.log('[error]', error);
      },
    });
  }

  sort(event: Sort) {
    this.form.patchValue({
      orderBy: event.active,
      orderDirection: event.direction,
      page: 1,
    });
  }
}
