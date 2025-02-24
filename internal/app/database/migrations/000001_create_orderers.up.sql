CREATE TABLE orderers (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255) NOT NULL,
    domain  VARCHAR(255) NOT NULL,
    port    integer NOT NULL
);