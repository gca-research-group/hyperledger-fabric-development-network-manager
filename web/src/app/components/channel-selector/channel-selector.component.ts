import { NgSelectComponent } from '@ng-select/ng-select';
import { finalize } from 'rxjs';

import { Component, inject, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';
import { Channel } from '@app/models';
import { ChannelService } from '@app/services/channel';

import { IconButtonComponent } from '../icon-button';

@Component({
  selector: 'app-channel-selector',
  templateUrl: './channel-selector.component.html',
  styleUrl: './channel-selector.component.scss',
  imports: [
    NgSelectComponent,
    FormsModule,
    ReactiveFormsModule,
    IconButtonComponent,
  ],
})
export class ChannelSelectorComponent
  extends CustomControlValueAccessorDirective
  implements OnInit
{
  channel: Channel[] = [];
  loading = false;
  private service = inject(ChannelService);

  override ngOnInit() {
    super.ngOnInit();
    this.getAllChannel();
  }

  getAllChannel() {
    this.loading = true;
    this.service
      .findAll()
      .pipe(
        finalize(() => {
          this.loading = false;
        }),
      )
      .subscribe(response => {
        this.channel = response.data;
      });
  }

  addChannel() {
    window.open('/channel/add', '_blank');
  }
}
