import { Routes } from '@angular/router';

import { isAuthenticatedGuard } from './guards';
import { channelsRoutes } from './modules/channels';
import { configRoutes } from './modules/config';
import { loginRoutes } from './modules/login';
import { orderersRoutes } from './modules/orderers';
import { peersRoutes } from './modules/peers';

export const routes: Routes = [
  ...loginRoutes,
  {
    path: '',
    children: [
      ...orderersRoutes,
      ...peersRoutes,
      ...channelsRoutes,
      ...configRoutes,
    ],
    canActivate: [isAuthenticatedGuard],
  },
];
