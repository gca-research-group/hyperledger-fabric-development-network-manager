import { Pipe, PipeTransform } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';

@Pipe({
  name: 'isRequired',
})
export class IsRequiredPipe implements PipeTransform {
  transform(control: FormControl): boolean {
    return !!control?.hasValidator(Validators.required);
  }
}
