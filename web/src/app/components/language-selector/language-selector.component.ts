import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { NgOptionComponent, NgSelectComponent } from '@ng-select/ng-select';

@Component({
  selector: 'app-language-selector',
  templateUrl: './language-selector.component.html',
  styleUrl: './language-selector.component.scss',
  imports: [FormsModule, NgSelectComponent, NgOptionComponent],
})
export class LanguageSelectorComponent {
  @Input()
  current = 'pt';

  @Output()
  select = new EventEmitter<string>();

  languages = [
    { id: 'pt', name: 'Português', flag: 'https://flagcdn.com/br.svg' },
    { id: 'es', name: 'Español', flag: 'https://flagcdn.com/es.svg' },
    { id: 'en', name: 'English', flag: 'https://flagcdn.com/gb.svg' },
  ];
}
