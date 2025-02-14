import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'smart-contracts',
    children: [
      {
        path: '',
        loadComponent: () =>
          import('./modules/smart-contract/list/list.component').then(
            m => m.ListComponent,
          ),
      },
      {
        path: 'add',
        loadComponent: () =>
          import('./modules/smart-contract/form/form.component').then(
            m => m.FormComponent,
          ),
      },
      {
        path: ':id',
        loadComponent: () =>
          import('./modules/smart-contract/form/form.component').then(
            m => m.FormComponent,
          ),
      },
    ],
  },
];
