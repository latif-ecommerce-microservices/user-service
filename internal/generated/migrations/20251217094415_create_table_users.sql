-- +goose Up
CREATE TABLE users (
    id              UUID primary key      default gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,
    phone_number    VARCHAR(255)
);

-- +goose Down
DROP TABLE users;
