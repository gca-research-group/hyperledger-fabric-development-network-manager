import { Directive, inject, Input, OnInit } from '@angular/core';
import {
  ControlValueAccessor,
  FormControl,
  FormControlDirective,
  FormControlName,
  NgControl,
  Validators,
} from '@angular/forms';

@Directive()
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export class CustomControlValueAccessorDirective<T = any>
  implements ControlValueAccessor, Validators, OnInit
{
  @Input()
  set required(required: boolean) {
    if (required) {
      this.formControl.setValidators([Validators.required]);
      return;
    }

    this.formControl.removeValidators(Validators.required);
  }

  formControl = new FormControl();
  public ngControl = inject(NgControl, { optional: true, self: true });

  constructor() {
    if (this.ngControl) {
      this.ngControl.valueAccessor = this;
    }
  }

  ngOnInit(): void {
    this.updateFormControl();
  }

  updateFormControl() {
    if (this.ngControl instanceof FormControlDirective && this.ngControl.form) {
      this.formControl = this.ngControl.form;
    }

    if (this.ngControl instanceof FormControlName && this.ngControl.control) {
      this.formControl = this.ngControl.control;
    }
  }
  writeValue(value: T): void {
    if (value !== this.formControl.value) {
      this.formControl.setValue(value);
    }
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  registerOnTouched(fn: any): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    if (!this.formControl) return;
    if (this.formControl.disabled === isDisabled) return;
    const key = isDisabled ? 'disable' : 'enable';
    this.formControl[key]({ emitEvent: false });
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  protected onChange = (value: T) => {
    this.updateFormControl();
  };

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  protected onTouched: any;
}
