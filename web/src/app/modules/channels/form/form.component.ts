import { TranslateModule } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';
import { finalize } from 'rxjs';

import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { ButtonComponent } from '@app/components/button';
import { InputComponent } from '@app/components/input';
import { OrderersSelectorComponent } from '@app/components/orderers-selector/orderers-selector.component';
import { PeersSelectorComponent } from '@app/components/peers-selector';
import { Channel } from '@app/models';
import { BreadcrumbService } from '@app/services/breadcrumb';
import { ChannelsService } from '@app/services/channels';

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
    OrderersSelectorComponent,
  ],
})
export class FormComponent implements OnInit, OnDestroy {
  form!: FormGroup<{
    id: FormControl<number | null>;
    name: FormControl<string | null>;
    peers: FormControl<number[] | null>;
    orderers: FormControl<number[] | null>;
  }>;

  private formBuilder = inject(FormBuilder);
  private breadcrumbService = inject(BreadcrumbService);
  private service = inject(ChannelsService);
  private toastr = inject(ToastrService);
  private activatedRoute = inject(ActivatedRoute);
  private router = inject(Router);
  loading = false;

  constructor() {
    this.form = this.formBuilder.group({
      id: new FormControl(),
      name: new FormControl('', Validators.required),
      peers: new FormControl<number[]>([]),
      orderers: new FormControl<number[]>([]),
    });

    this.breadcrumbService.update([
      ...BREADCRUMB,
      {
        label: 'add',
      },
    ]);
  }

  ngOnInit(): void {
    const id = this.activatedRoute.snapshot.params['id'] as unknown as number;
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
          orderers: channel.orderers.map(orderer => orderer.id),
        });
      },
      error: (error: { error?: { message: string } }) => {
        this.toastr.error(error.error?.message ?? 'INTERNAL_SERVER_ERROR');
      },
    });
  }

  save() {
    if (this.form.invalid) {
      this.toastr.warning('INVALID_FORM');
      return;
    }

    this.loading = true;
    this.service
      .save({ ...this.form.value } as unknown as Channel)
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

          this.toastr.success(message);
          void this.router.navigate(['./..'], {
            relativeTo: this.activatedRoute,
          });
        },
        error: (error: { error?: { message: string } }) => {
          this.toastr.error(error.error?.message ?? 'INTERNAL_SERVER_ERROR');
        },
      });
  }
}
