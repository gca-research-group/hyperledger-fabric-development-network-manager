CREATE TABLE channel_peers (
    channel_id  INTEGER REFERENCES channels,
    peer_id     INTEGER REFERENCES peers,
    crated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);