-- +goose Up
-- +goose StatementBegin

CREATE TABLE if NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role_id INT NOT NULL REFERENCES roles(id),
    last_access TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd