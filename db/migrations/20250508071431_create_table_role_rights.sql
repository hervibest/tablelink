-- +goose Up
-- +goose StatementBegin
CREATE TABLE if NOT EXISTS role_rights ( 
    id SERIAL PRIMARY KEY,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    section TEXT NOT NULL,
    route TEXT NOT NULL,
    r_created BOOLEAN NOT NULL DEFAULT FALSE,
    r_read BOOLEAN NOT NULL DEFAULT FALSE,
    r_update BOOLEAN NOT NULL DEFAULT FALSE,
    r_delete  BOOLEAN NOT NULL DEFAULT FALSE

)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS role_rights;
-- +goose StatementEnd
