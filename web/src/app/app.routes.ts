import { Routes } from '@angular/router';
import { configurationFilesRoutes } from './modules/configuration-files';

export const routes: Routes = [
  { path: '', redirectTo: 'configuration-files', pathMatch: 'full' },
  ...configurationFilesRoutes,
];
