-- +goose Up

-- +goose StatementBegin
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email varchar(128) NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL UNIQUE);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_hash  (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email varchar(128) NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL UNIQUE,
    code varchar(138) NOT NULL,
    time_reg TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS app (
    id SERIAL PRIMARY KEY NOT NULL,
    app_name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TABLE users_hash;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE app;
-- +goose StatementEnd

