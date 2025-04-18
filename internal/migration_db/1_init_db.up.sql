CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email varchar(128) NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS users_hash  (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email varchar(128) NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL UNIQUE,
    code varchar(138) NOT NULL,
    time_reg TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS app (
    id SERIAL PRIMARY KEY NOT NULL,
    app_name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);