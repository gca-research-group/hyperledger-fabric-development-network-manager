import { Routes } from '@angular/router';

import { isAuthenticatedGuard } from './guards';
import { channelRoutes } from './modules/channel';
import { loginRoutes } from './modules/login';
import { ordererRoutes } from './modules/orderer';
import { peerRoutes } from './modules/peer';

export const routes: Routes = [
  ...loginRoutes,
  {
    path: '',
    children: [...channelRoutes, ...ordererRoutes, ...peerRoutes],
    canActivate: [isAuthenticatedGuard],
  },
];
