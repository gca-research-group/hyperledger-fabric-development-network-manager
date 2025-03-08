import { NgOptionComponent, NgSelectComponent } from '@ng-select/ng-select';

import { Component, model, output } from '@angular/core';
import { FormsModule } from '@angular/forms';

import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';

@Component({
  selector: 'app-blockchain-platform-selector',
  templateUrl: './blockchain-platform-selector.component.html',
  styleUrl: './blockchain-platform-selector.component.scss',
  imports: [FormsModule, NgSelectComponent, NgOptionComponent],
})
export class BlockchainPlatformSelectorComponent extends CustomControlValueAccessorDirective {
  current = model('hyperledger-fabric');

  select = output<string>();

  blockchainPlatforms = [
    {
      id: 'hyperledger-fabric',
      name: 'Hyperledger Fabric',
      image: '/images/hyperledger.png',
    },
    {
      id: 'ethereum',
      name: 'Ethereum',
      image: '/images/ethereum.png',
    },
  ];
}
