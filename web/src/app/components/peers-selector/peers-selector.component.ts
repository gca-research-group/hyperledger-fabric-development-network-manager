import { NgSelectComponent } from '@ng-select/ng-select';
import { finalize } from 'rxjs';

import { Component, inject, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';
import { Peer } from '@app/models';
import { PeersService } from '@app/services/peers';

import { IconButtonComponent } from '../icon-button';

@Component({
  selector: 'app-peers-selector',
  templateUrl: './peers-selector.component.html',
  styleUrl: './peers-selector.component.scss',
  imports: [
    NgSelectComponent,
    FormsModule,
    ReactiveFormsModule,
    IconButtonComponent,
  ],
})
export class PeersSelectorComponent
  extends CustomControlValueAccessorDirective
  implements OnInit
{
  peers: Peer[] = [];
  loading = false;
  private service = inject(PeersService);

  override ngOnInit() {
    super.ngOnInit();
    this.getAllPeers();
  }

  getAllPeers() {
    this.loading = true;
    this.service
      .findAll()
      .pipe(
        finalize(() => {
          this.loading = false;
        }),
      )
      .subscribe(response => {
        this.peers = response.data;
      });
  }

  addPeer() {
    window.open('/peers/add', '_blank');
  }
}
