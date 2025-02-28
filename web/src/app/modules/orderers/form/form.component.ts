import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { InputComponent } from '../../../components/input/input.component';
import { TranslateModule } from '@ngx-translate/core';
import { BreadcrumbService } from '@app/services/breadcrumb';
import { ToastrService } from 'ngx-toastr';
import { ButtonComponent } from '@app/components/button';
import { OrderersService } from '../services/orderers.service';
import { Location } from '@angular/common';
import { finalize } from 'rxjs';
import { ActivatedRoute } from '@angular/router';

const BREADCRUMB = [
  {
    label: 'home',
    url: '/',
  },
  {
    label: 'orderers',
    url: '/orderers',
  },
];

@Component({
  selector: 'app-orderers-form',
  templateUrl: './form.component.html',
  styleUrl: './form.component.scss',
  host: { class: 'd-md-flex d-sm-block justify-content-center' },
  imports: [
    ReactiveFormsModule,
    FormsModule,
    TranslateModule,
    InputComponent,
    ButtonComponent,
  ],
})
export class FormComponent implements OnInit, OnDestroy {
  form!: FormGroup;

  private formBuilder = inject(FormBuilder);
  private breadcrumbService = inject(BreadcrumbService);
  private service = inject(OrderersService);
  private location = inject(Location);
  private activatedRoute = inject(ActivatedRoute);
  loading = false;

  private toastr = inject(ToastrService);

  constructor() {
    this.form = this.formBuilder.group({
      id: null,
      name: [null, Validators.required],
      domain: [null, Validators.required],
      port: [null, Validators.required],
    });

    this.breadcrumbService.update([
      ...BREADCRUMB,
      {
        label: 'add',
      },
    ]);
  }

  ngOnInit(): void {
    const id = this.activatedRoute.snapshot.params['id'];
    if (id) {
      this.find(id);
      this.breadcrumbService.update([...BREADCRUMB, { label: 'edit' }]);
    }
  }

  ngOnDestroy(): void {
    this.breadcrumbService.reset();
  }

  find(id: number) {
    this.service.findById(id).subscribe({
      next: orderer => {
        this.form.patchValue(orderer);
      },
      error: error => {
        this.toastr.error(error.message, undefined, {
          closeButton: true,
          progressBar: true,
        });
      },
    });
  }

  save() {
    if (this.form.invalid) {
      this.toastr.warning('INVALID_FORM', undefined, {
        closeButton: true,
        progressBar: true,
      });
      return;
    }

    this.loading = true;
    this.service
      .save({ ...this.form.value, port: +this.form.value.port })
      .pipe(
        finalize(() => {
          this.loading = false;
        }),
      )
      .subscribe({
        next: () => {
          this.toastr.success('RECORD_CREATED_SUCCESSFULLY', undefined, {
            closeButton: true,
            progressBar: true,
          });
          this.location.back();
        },
        error: error => {
          this.toastr.error(error.message, undefined, {
            closeButton: true,
            progressBar: true,
          });
        },
      });
  }
}
