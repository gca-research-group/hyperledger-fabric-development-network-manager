import { DatePipe } from '@angular/common';
import { inject, Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'date',
})
export class CustomDatePipe implements PipeTransform {
  private datePipe = inject(DatePipe);
  transform(value: string, format: string, locale: string): string | null {
    return this.datePipe.transform(value, format, undefined, locale);
  }
}
