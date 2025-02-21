import { Component, inject, OnDestroy } from '@angular/core';
import {
  FormArray,
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { InputComponent } from '../../../components/input/input.component';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { BreadcrumbService } from '@app/services/breadcrumb';
import { IconButtonComponent } from '@app/components/icon-button';
import { ToastrService } from 'ngx-toastr';
import { NgSelectComponent } from '@ng-select/ng-select';

const BREADCRUMB = [
  {
    label: 'home',
    url: '/',
  },
  {
    label: 'configuration-files',
    url: '/configuration-files',
  },
  {
    label: 'add',
  },
];

type Orderer = {
  name: string;
  domain: string;
  port: string;
};

type Org = {
  name: string;
  domain: string;
  port: string;
};

type Channel = {
  name: string;
  orgs: string[][];
};

@Component({
  selector: 'app-configuration-files-form',
  templateUrl: './form.component.html',
  styleUrl: './form.component.scss',
  host: { class: 'd-md-flex d-sm-block justify-content-center' },
  imports: [
    ReactiveFormsModule,
    FormsModule,
    TranslateModule,
    NgSelectComponent,
    InputComponent,
    IconButtonComponent,
  ],
})
export class FormComponent implements OnDestroy {
  form!: FormGroup;
  orderers!: FormArray;
  orgs!: FormArray;
  channels!: FormArray;

  formBuilder = inject(FormBuilder);
  breadcrumbService = inject(BreadcrumbService);

  private toastr = inject(ToastrService);
  private translateService = inject(TranslateService);

  constructor() {
    this.form = this.formBuilder.group({
      orderers: this.formBuilder.array([
        this.formBuilder.group<Orderer>({
          name: 'Orderer',
          domain: 'orderer.example.com',
          port: '7050',
        }),
      ]),
      orgs: this.formBuilder.array([
        this.formBuilder.group<Org>({
          name: 'Org1',
          domain: 'org1.example.com',
          port: '7051',
        }),
      ]),
      channels: this.formBuilder.array([
        this.formBuilder.group<Channel>({
          name: 'Channel1',
          orgs: [['Org1']],
        }),
      ]),
    });

    this.orderers = this.form.get('orderers') as FormArray;
    this.orgs = this.form.get('orgs') as FormArray;
    this.channels = this.form.get('channels') as FormArray;

    this.breadcrumbService.update(BREADCRUMB);
  }

  ngOnDestroy(): void {
    this.breadcrumbService.reset();
  }

  addOrderer() {
    this.orderers.push(
      this.formBuilder.group({
        name: 'orderer',
        domain: 'orderer.example.com',
        port: 7050,
      }),
    );
  }

  removeOrderer(index: number) {
    if (this.orderers.length === 1) {
      this.toastr.error(
        this.translateService.instant('YOU_MUST_HAVE_AT_LEAST_ONE_ORDERER'),
        undefined,
        {
          closeButton: true,
          progressBar: true,
        },
      );
      return;
    }
    this.orderers.removeAt(index);
  }

  addOrg() {
    this.orgs.push(
      this.formBuilder.group({
        name: `Org${this.orgs.length + 1}`,
        domain: `org${this.orgs.length + 1}.example.com`,
        port: 7051,
      }),
    );
  }

  removeOrg(index: number) {
    if (this.orgs.length === 1) {
      this.toastr.error(
        this.translateService.instant('YOU_MUST_HAVE_AT_LEAST_ONE_ORG'),
        undefined,
        {
          closeButton: true,
          progressBar: true,
        },
      );
      return;
    }
    this.orgs.removeAt(index);
  }

  addChannel() {
    const orgs = (this.orgs.value as Org[]).map(item => item.name);

    this.channels.push(
      this.formBuilder.group({
        name: `Channel${this.channels.length + 1}`,
        orgs: [orgs],
      }),
    );
  }

  removeChannel(index: number) {
    if (this.channels.length === 1) {
      this.toastr.error(
        this.translateService.instant('YOU_MUST_HAVE_AT_LEAST_ONE_CHANNEL'),
        undefined,
        {
          closeButton: true,
          progressBar: true,
        },
      );
      return;
    }
    this.channels.removeAt(index);
  }
}
