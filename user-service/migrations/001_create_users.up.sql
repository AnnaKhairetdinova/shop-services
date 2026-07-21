CREATE TABLE IF NOT EXISTS users
(
    uuid       UUID PRIMARY KEY,
    name       TEXT      NOT NULL,
    email      TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL
);
