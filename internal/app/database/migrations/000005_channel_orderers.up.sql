CREATE TABLE channel_orderers (
  channel_id INTEGER REFERENCES channels,
  orderer_id INTEGER REFERENCES orderers,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT unique_channel_orderer UNIQUE (channel_id, orderer_id),
  updated_at TIMESTAMP
);