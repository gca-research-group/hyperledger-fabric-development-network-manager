import { Peer } from './peer.model';

export type Channel = {
  id: number;
  name: string;
  peers: Peer[];
};
