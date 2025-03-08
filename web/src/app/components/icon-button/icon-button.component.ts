import { TranslateModule } from '@ngx-translate/core';

import { Component, input } from '@angular/core';
import { MatIconButton } from '@angular/material/button';
import { MatIcon } from '@angular/material/icon';
import { MatTooltip, TooltipPosition } from '@angular/material/tooltip';

@Component({
  selector: 'app-icon-button',
  templateUrl: './icon-button.component.html',
  styleUrl: './icon-button.component.scss',
  imports: [MatIconButton, MatIcon, MatTooltip, TranslateModule],
})
export class IconButtonComponent {
  icon = input<string>();
  tooltip = input<string>('');
  color = input<string>();
  tooltipPosition = input<TooltipPosition>('below');
}
