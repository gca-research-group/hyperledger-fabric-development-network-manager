import { TranslateModule } from '@ngx-translate/core';

import { Component, input, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';

import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';
import { IsRequiredPipe } from '@app/pipes';

import { IconComponent } from '../icon';

@Component({
  selector: 'app-input',
  templateUrl: './input.component.html',
  styleUrl: './input.component.scss',
  imports: [
    MatFormFieldModule,
    MatInputModule,
    FormsModule,
    ReactiveFormsModule,
    TranslateModule,
    IconComponent,
    IsRequiredPipe,
  ],
})
export class InputComponent
  extends CustomControlValueAccessorDirective
  implements OnInit
{
  icon = input('');
  label = input('');
  placeholder = input('');
  type = input('text');
  min = input<number>();
}
