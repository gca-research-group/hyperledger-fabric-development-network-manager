CREATE TABLE peers (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    domain      VARCHAR(255) NOT NULL,
    port        INTEGER NOT NULL,
    crated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);