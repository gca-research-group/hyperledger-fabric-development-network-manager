import { Component, Input } from '@angular/core';
import { MatIcon } from '@angular/material/icon';
import { MatTooltip } from '@angular/material/tooltip';

@Component({
  selector: 'app-icon',
  templateUrl: './icon.component.html',
  styleUrl: './icon.component.scss',
  imports: [MatIcon, MatTooltip],
})
export class IconComponent {
  @Input() icon!: string;
  @Input() tooltip!: string;
}
