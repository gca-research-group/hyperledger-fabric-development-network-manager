CREATE TABLE channels (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    crated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);