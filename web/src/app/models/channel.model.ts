import { Orderer } from './orderer.model';
import { Peer } from './peer.model';

export type Channel = {
  id: number;
  name: string;
  peers: Peer[];
  orderers: Orderer[];
};
