CREATE TABLE address
(
    id         VARCHAR PRIMARY KEY,
    name       VARCHAR NOT NULL,
owner VARCHAR NOT NULL,
pubkey VARCHAR NOT NULL,
currency VARCHAR NOT NULL,
is_primary BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
