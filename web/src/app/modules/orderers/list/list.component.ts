import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ElementRef,
  inject,
  OnDestroy,
  OnInit,
  viewChild,
} from '@angular/core';
import { BreadcrumbService } from '@app/services/breadcrumb';
import { OrderersService } from '../services/orderers.service';
import { ColumnType, Orderer } from '@app/models';
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
    type: ColumnType.DATETIME,
  },
  {
    field: 'updatedAt',
    label: 'updatedAt',
    type: ColumnType.DATETIME,
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
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
  imports: [
    TableComponent,
    ReactiveFormsModule,
    FormsModule,
    TranslateModule,
    InputComponent,
  ],
})
export class ListComponent implements OnInit, AfterViewInit, OnDestroy {
  columns = COLUMNS;

  displayedColumns = COLUMNS.map(column => column.field);

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

  constructor() {
    this.breadcrumbService.update(BREADCRUMB);

    this.form = this.formBuilder.group({
      id: null,
      name: null,
      domain: null,
      port: null,
      page: 1,
      pageSize: 20,
    });

    this.form.valueChanges.pipe(debounceTime(300)).subscribe(value => {
      this.search(value);
    });
  }

  ngOnInit(): void {
    this.search(this.form.value);
  }

  ngAfterViewInit(): void {
    const form = this.formElement()?.nativeElement;

    const marginBottom = getComputedStyle(
      this.filters()?.nativeElement as Element,
    ).marginBottom;

    this.tableHeight = `calc(100vh - var(--hfdnm-toolbar-height) - (2 * var(--hfdnm-content-vertical-padding)) - ${form?.offsetHeight}px - ${marginBottom})`;
    this.cdk.detectChanges();
  }

  ngOnDestroy(): void {
    this.breadcrumbService.reset();
  }

  scroll() {
    if (this.hasMore) {
      this.form.patchValue({ page: (this.form.get('page')?.value ?? 0) + 1 });
      this.findAll();
    }
  }

  findAll() {
    this.service.findAll(this.form.value).subscribe({
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
      Object.entries(obj).filter(([_, value]) => value !== null),
    ) as Partial<T>;
  }

  search(params: object) {
    const _params = this.removeNullFields(params);
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
}
