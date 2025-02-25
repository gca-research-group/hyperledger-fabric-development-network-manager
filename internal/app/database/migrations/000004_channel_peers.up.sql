CREATE TABLE channel_peers (
    channel_id  INTEGER REFERENCES channels,
    peer_id     INTEGER REFERENCES peers,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_channel_peer UNIQUE (channel_id, peer_id)
);