import { loginGuard } from '@app/guards';

export const loginRoutes = [
  {
    path: 'login',
    loadComponent: () =>
      import('./login.component').then(m => m.LoginComponent),
    canActivate: [loginGuard],
  },
];
