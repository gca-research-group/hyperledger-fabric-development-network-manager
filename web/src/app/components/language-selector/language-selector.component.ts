import { Component, effect, model } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';
import { NgOptionComponent, NgSelectComponent } from '@ng-select/ng-select';

@Component({
  selector: 'app-language-selector',
  templateUrl: './language-selector.component.html',
  styleUrl: './language-selector.component.scss',
  imports: [
    FormsModule,
    ReactiveFormsModule,
    NgSelectComponent,
    NgOptionComponent,
  ],
})
export class LanguageSelectorComponent extends CustomControlValueAccessorDirective<string> {
  language = model('pt');

  languages = [
    { id: 'pt', name: 'Português', flag: 'https://flagcdn.com/br.svg' },
    { id: 'es', name: 'Español', flag: 'https://flagcdn.com/es.svg' },
    { id: 'en', name: 'English', flag: 'https://flagcdn.com/gb.svg' },
  ];

  constructor() {
    super();
    effect(() => {
      this.formControl.patchValue(this.language());
    });
  }

  change(language: string) {
    this.language.update(() => language);
  }
}
