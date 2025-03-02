import { Routes } from '@angular/router';
import { orderersRoutes } from './modules/orderers';
import { peersRoutes } from './modules/peers';
import { channelsRoutes } from './modules/channels';

export const routes: Routes = [
  ...orderersRoutes,
  ...peersRoutes,
  ...channelsRoutes,
];
