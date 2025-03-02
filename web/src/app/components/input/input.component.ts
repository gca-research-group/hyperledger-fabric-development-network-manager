import { Component, input, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { CustomControlValueAccessorDirective } from '@app/directives/custom-control-value-accessor';
import { TranslateModule } from '@ngx-translate/core';

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
  ],
})
export class InputComponent
  extends CustomControlValueAccessorDirective
  implements OnInit
{
  label = input('');
  placeholder = input('');
  type = input('text');
  min = input<number>();
}
