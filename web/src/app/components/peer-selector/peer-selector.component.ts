import { NgSelectComponent } from '@ng-select/ng-select';
import { finalize } from 'rxjs';

import { Component, inject, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';
import { Peer } from '@app/models';
import { PeerService } from '@app/services/peer';

import { IconButtonComponent } from '../icon-button';

@Component({
  selector: 'app-peer-selector',
  templateUrl: './peer-selector.component.html',
  styleUrl: './peer-selector.component.scss',
  imports: [
    NgSelectComponent,
    FormsModule,
    ReactiveFormsModule,
    IconButtonComponent,
  ],
})
export class PeerSelectorComponent
  extends CustomControlValueAccessorDirective
  implements OnInit
{
  peer: Peer[] = [];
  loading = false;
  private service = inject(PeerService);

  override ngOnInit() {
    super.ngOnInit();
    this.getAllPeer();
  }

  getAllPeer() {
    this.loading = true;
    this.service
      .findAll()
      .pipe(
        finalize(() => {
          this.loading = false;
        }),
      )
      .subscribe(response => {
        this.peer = response.data;
      });
  }

  addPeer() {
    window.open('/peer/add', '_blank');
  }
}
