import { Component, inject, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';
import { Peer } from '@app/models';
import { PeersService } from '@app/modules/peers/services/peers.service';
import { NgSelectComponent } from '@ng-select/ng-select';

@Component({
  selector: 'app-peers-selector',
  templateUrl: './peers-selector.component.html',
  styleUrl: './peers-selector.component.scss',
  imports: [NgSelectComponent, FormsModule, ReactiveFormsModule],
})
export class PeersSelectorComponent
  extends CustomControlValueAccessorDirective
  implements OnInit
{
  peers: Peer[] = [];
  private service = inject(PeersService);

  override ngOnInit() {
    super.ngOnInit();
    this.getAllPeers();
  }

  getAllPeers() {
    this.service.findAll().subscribe(response => {
      this.peers = response.data;
    });
  }
}
