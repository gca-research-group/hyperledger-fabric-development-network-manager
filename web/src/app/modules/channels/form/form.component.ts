import { TranslateModule } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';
import { finalize } from 'rxjs';

import { Location } from '@angular/common';
import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ActivatedRoute } from '@angular/router';

import { ButtonComponent } from '@app/components/button';
import { PeersSelectorComponent } from '@app/components/peers-selector';
import { BreadcrumbService } from '@app/services/breadcrumb';

import { InputComponent } from '../../../components/input/input.component';
import { ChannelsService } from '../services/channels.service';

const BREADCRUMB = [
  {
    label: 'home',
    url: '/',
  },
  {
    label: 'channels',
    url: '/channels',
  },
];

@Component({
  selector: 'app-channels-form',
  templateUrl: './form.component.html',
  styleUrl: './form.component.scss',
  host: { class: 'd-md-flex d-sm-block justify-content-center' },
  imports: [
    ReactiveFormsModule,
    FormsModule,
    TranslateModule,
    InputComponent,
    ButtonComponent,
    PeersSelectorComponent,
  ],
})
export class FormComponent implements OnInit, OnDestroy {
  form!: FormGroup;

  private formBuilder = inject(FormBuilder);
  private breadcrumbService = inject(BreadcrumbService);
  private service = inject(ChannelsService);
  private location = inject(Location);
  private activatedRoute = inject(ActivatedRoute);
  loading = false;

  private toastr = inject(ToastrService);

  constructor() {
    this.form = this.formBuilder.group({
      id: null,
      name: [null, Validators.required],
      peers: [[]],
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
      next: channel => {
        this.form.patchValue({
          ...channel,
          peers: channel.peers.map(peer => peer.id),
        });
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
          const message = this.form.value.id
            ? 'RECORD_UPDATED_SUCCESSFULLY'
            : 'RECORD_CREATED_SUCCESSFULLY';

          this.toastr.success(message, undefined, {
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
