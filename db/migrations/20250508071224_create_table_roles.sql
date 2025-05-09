-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO roles (name) VALUES 
    ('SUPER_ADMIN'),
    ('ADMIN'),
    ('STAFF');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
