import { NgSelectComponent } from '@ng-select/ng-select';
import { finalize } from 'rxjs';

import { Component, inject, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';
import { Orderer } from '@app/models';
import { OrdererService } from '@app/services/orderer';

import { IconButtonComponent } from '../icon-button';

@Component({
  selector: 'app-orderer-selector',
  templateUrl: './orderer-selector.component.html',
  styleUrl: './orderer-selector.component.scss',
  imports: [
    NgSelectComponent,
    FormsModule,
    ReactiveFormsModule,
    IconButtonComponent,
  ],
})
export class OrdererSelectorComponent
  extends CustomControlValueAccessorDirective
  implements OnInit
{
  orderer: Orderer[] = [];
  loading = false;
  private service = inject(OrdererService);

  override ngOnInit() {
    super.ngOnInit();
    this.getAllOrderer();
  }

  getAllOrderer() {
    this.loading = true;
    this.service
      .findAll()
      .pipe(
        finalize(() => {
          this.loading = false;
        }),
      )
      .subscribe(response => {
        this.orderer = response.data;
      });
  }

  addOrderer() {
    window.open('/orderer/add', '_blank');
  }
}
