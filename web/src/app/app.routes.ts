import { Routes } from '@angular/router';

import { channelsRoutes } from './modules/channels';
import { orderersRoutes } from './modules/orderers';
import { peersRoutes } from './modules/peers';

export const routes: Routes = [
  ...orderersRoutes,
  ...peersRoutes,
  ...channelsRoutes,
];
