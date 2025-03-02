import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ElementRef,
  inject,
  OnDestroy,
  OnInit,
  TemplateRef,
  viewChild,
} from '@angular/core';
import { BreadcrumbService } from '@app/services/breadcrumb';
import { OrderersService } from '../services/orderers.service';
import { Breadcrumb, Column, ColumnType, Orderer } from '@app/models';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { InputComponent } from '@app/components/input';
import { debounceTime } from 'rxjs';
import { TableComponent } from '@app/components/table';
import { IconButtonComponent } from '@app/components/icon-button';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { DeleteDialogComponent } from '@app/components/delete-dialog';
import { MatDialog } from '@angular/material/dialog';
import { ToastrService } from 'ngx-toastr';
import { Sort } from '@angular/material/sort';

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
    id: 'domain',
    label: 'domain',
  },
  {
    id: 'port',
    label: 'port',
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
    label: 'orderers',
  },
];

@Component({
  selector: 'app-orderers-list',
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
export class ListComponent implements OnInit, AfterViewInit, OnDestroy {
  columns = COLUMNS;

  displayedColumns = COLUMNS.map(column => column.id);

  private breadcrumbService = inject(BreadcrumbService);
  private service = inject(OrderersService);
  private formBuilder = inject(FormBuilder);

  data: Orderer[] = [];

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
      port: null,
      page: 1,
      pageSize: 20,
      orderBy: null,
      orderDirection: null,
    });

    const queryParams = this.activatedRoute.snapshot.queryParams;
    this.form.patchValue(queryParams, { emitEvent: false });

    this.form.valueChanges.pipe(debounceTime(300)).subscribe(() => {
      this.search();
      this.router.navigate([], {
        queryParams: this.form.value,
        queryParamsHandling: 'merge',
      });
    });
  }

  ngOnInit(): void {
    this.search();
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

  openDialog(item: Orderer): void {
    const dialogRef = this.dialog.open(DeleteDialogComponent, {
      data: item.id,
    });

    dialogRef.afterClosed().subscribe((id: number | undefined) => {
      if (id) {
        this.service.delete(id).subscribe({
          next: () => {
            this.data = this.data.filter(item => item.id !== id);
            this.toastr.success('DELETED_SUCCESSFULLY', undefined, {
              closeButton: true,
              progressBar: true,
            });
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
        { page: (this.form.get('page')?.value ?? 0) + 1 },
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
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
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
