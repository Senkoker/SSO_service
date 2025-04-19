-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_info (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    first_name VARCHAR(64) NOT NULL,
    second_name VARCHAR(64) NOT NULL,
    img_url TEXT NOT NULL,
    birth_date DATE,
    education TEXT,
    country VARCHAR(128),
    city VARCHAR,
    CONSTRAINT user_id_fr FOREIGN KEY (user_id) REFERENCES users(id)
    );
-- +goose StatementEnd 

-- +goose Down

-- +goose StatementBegin
DROP TABLE users_info;
-- +goose StatementEnd