import { NgOptionComponent, NgSelectComponent } from '@ng-select/ng-select';
import { TranslateModule } from '@ngx-translate/core';

import { Component } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';

@Component({
  selector: 'app-blockchain-platform-selector',
  templateUrl: './blockchain-platform-selector.component.html',
  styleUrl: './blockchain-platform-selector.component.scss',
  imports: [
    FormsModule,
    ReactiveFormsModule,
    NgSelectComponent,
    NgOptionComponent,
    TranslateModule,
  ],
})
export class BlockchainPlatformSelectorComponent extends CustomControlValueAccessorDirective {
  blockchainPlatforms = [
    {
      id: 'HYPERLEDGER_FABRIC',
      name: 'Hyperledger Fabric',
      image: '/images/hyperledger.png',
    },
    {
      id: 'ETHEREUM',
      name: 'Ethereum',
      image: '/images/ethereum.png',
    },
  ];
}
