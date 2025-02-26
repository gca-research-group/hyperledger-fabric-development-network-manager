import { Component, inject, OnDestroy } from '@angular/core';
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

const BREADCRUMB = [
  {
    label: 'home',
    url: '/',
  },
  {
    label: 'orderers',
    url: '/orderers',
  },
  {
    label: 'add',
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
export class FormComponent implements OnDestroy {
  form!: FormGroup;

  private formBuilder = inject(FormBuilder);
  private breadcrumbService = inject(BreadcrumbService);
  private service = inject(OrderersService);
  private location = inject(Location);
  loading = false;

  private toastr = inject(ToastrService);

  constructor() {
    this.form = this.formBuilder.group({
      name: ['Orderer', Validators.required],
      domain: ['orderer.example.com', Validators.required],
      port: [7050, Validators.required],
    });

    this.breadcrumbService.update(BREADCRUMB);
  }

  ngOnDestroy(): void {
    this.breadcrumbService.reset();
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
    this.service.save(this.form.value).subscribe({
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
