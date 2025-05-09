-- +goose Up
-- +goose StatementBegin
CREATE TABLE if NOT EXISTS role_rights ( 
    id SERIAL PRIMARY KEY,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    section TEXT NOT NULL,
    route TEXT NOT NULL,
    r_create SMALLINT NOT NULL DEFAULT 0,
    r_read SMALLINT NOT NULL DEFAULT 0,
    r_update SMALLINT NOT NULL DEFAULT 0,
    r_delete  SMALLINT NOT NULL DEFAULT 0
);

INSERT INTO role_rights (role_id, section, route, r_create, r_read, r_update, r_delete) VALUES 
    (1, 'be','UsersService/CreateUser', 1, 0, 0, 0),
    (1, 'be','UsersService/GetAllUser', 0, 1, 0, 0),
    (1, 'be','UsersService/UpdateUser', 0, 0, 1, 0),
    (1, 'be','UsersService/DeleteUser', 0, 0, 0, 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS role_rights;
-- +goose StatementEnd
