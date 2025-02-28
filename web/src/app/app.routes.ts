import { Routes } from '@angular/router';
import { configurationFilesRoutes } from './modules/configuration-files';
import { orderersRoutes } from './modules/orderers';
import { peersRoutes } from './modules/peers';

export const routes: Routes = [
  { path: '', redirectTo: 'configuration-files', pathMatch: 'full' },
  ...configurationFilesRoutes,
  ...orderersRoutes,
  ...peersRoutes,
];
