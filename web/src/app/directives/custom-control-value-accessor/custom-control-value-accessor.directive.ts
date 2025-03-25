import {
  ChangeDetectorRef,
  Directive,
  inject,
  Input,
  OnInit,
} from '@angular/core';
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
  private cdr = inject(ChangeDetectorRef);

  @Input()
  set required(required: boolean) {
    if (required) {
      this.formControl.setValidators([control => Validators.required(control)]);
      return;
    }

    this.formControl.removeValidators(control => Validators.required(control));
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

  writeValue(event: unknown): void {
    const value =
      (event as { target: { value: unknown } })?.target?.value || event;

    if (value !== this.formControl.value) {
      setTimeout(() => {
        this.formControl.setValue(value, { emitEvent: false });
        this.formControl.updateValueAndValidity();
        this.cdr.detectChanges();
      });
    }
  }

  registerOnChange(fn: (_value: T) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: (_value: T) => void): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    if (!this.formControl) return;
    if (this.formControl.disabled === isDisabled) return;
    const key = isDisabled ? 'disable' : 'enable';
    this.formControl[key]({ emitEvent: false });
  }

  protected onChange = (_value: T) => {
    this.updateFormControl();
  };

  protected onTouched?: (_value: T) => void;
}
