import { NgSelectComponent } from '@ng-select/ng-select';
import { finalize } from 'rxjs';

import { Component, inject, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';
import { Orderer } from '@app/models';
import { OrderersService } from '@app/services/orderers';

import { IconButtonComponent } from '../icon-button';

@Component({
  selector: 'app-orderers-selector',
  templateUrl: './orderers-selector.component.html',
  styleUrl: './orderers-selector.component.scss',
  imports: [
    NgSelectComponent,
    FormsModule,
    ReactiveFormsModule,
    IconButtonComponent,
  ],
})
export class OrderersSelectorComponent
  extends CustomControlValueAccessorDirective
  implements OnInit
{
  orderers: Orderer[] = [];
  loading = false;
  private service = inject(OrderersService);

  override ngOnInit() {
    super.ngOnInit();
    this.getAllOrderers();
  }

  getAllOrderers() {
    this.loading = true;
    this.service
      .findAll()
      .pipe(
        finalize(() => {
          this.loading = false;
        }),
      )
      .subscribe(response => {
        this.orderers = response.data;
      });
  }

  addOrderer() {
    window.open('/orderers/add', '_blank');
  }
}
