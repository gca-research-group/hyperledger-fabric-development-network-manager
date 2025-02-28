import { Component, input } from '@angular/core';
import { MatIcon } from '@angular/material/icon';
import { MatTooltip } from '@angular/material/tooltip';

@Component({
  selector: 'app-icon',
  templateUrl: './icon.component.html',
  styleUrl: './icon.component.scss',
  imports: [MatIcon, MatTooltip],
})
export class IconComponent {
  icon = input<string>();
  tooltip = input<string>();
}
