import { Component, input } from '@angular/core';
import { MatButton } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-button',
  templateUrl: './button.component.html',
  styleUrl: './button.component.scss',
  host: {
    '[attr.disabled]': 'disabled()',
  },
  imports: [MatButton, TranslateModule],
})
export class ButtonComponent {
  label = input<string>('');
  type = input<string>('button');
  disabled = input<boolean>();
}
