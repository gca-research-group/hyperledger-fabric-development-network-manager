import { Component, inject } from '@angular/core';
import {
  FormBuilder,
  FormControl,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { Router } from '@angular/router';

import { ButtonComponent } from '@app/components/button';
import { InputComponent } from '@app/components/input';
import { Login } from '@app/models';
import { AuthService } from '@app/services/auth';
import { CurrentUserService } from '@app/services/current-user';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
  imports: [ReactiveFormsModule, InputComponent, ButtonComponent],
})
export class LoginComponent {
  form!: FormGroup<{
    email: FormControl<string | null>;
    password: FormControl<string | null>;
  }>;

  private formBuilder = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);
  private currentUserService = inject(CurrentUserService);

  constructor() {
    this.form = this.formBuilder.group({
      email: new FormControl('', [Validators.required]),
      password: new FormControl('', [Validators.required]),
    });
  }

  login() {
    if (this.form.invalid) {
      return;
    }

    this.authService.login(this.form.value as unknown as Login).subscribe({
      next: response => {
        this.currentUserService.add({ ...response, isAuthenticated: true });
        void this.router.navigate(['/']);
      },
      error: error => {
        console.error('[error]', error);
      },
    });
  }
}
