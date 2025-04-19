-- +goose Up

-- +goose StatementBegin
INSERT INTO app(app_name,secret)
VALUES ('first','my_secret')
    ON CONFLICT DO NOTHING;
-- +goose StatementEnd 
