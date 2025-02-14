import { Component } from '@angular/core';
import { Input } from '@angular/core';
import { MatIconButton } from '@angular/material/button';
import { MatIcon } from '@angular/material/icon';
import { MatTooltip } from '@angular/material/tooltip';

@Component({
  selector: 'app-icon-button',
  templateUrl: './icon-button.component.html',
  styleUrl: './icon-button.component.scss',
  imports: [MatIconButton, MatIcon, MatTooltip],
})
export class IconButtonComponent {
  @Input() icon!: string;
  @Input() tooltip!: string;
}
