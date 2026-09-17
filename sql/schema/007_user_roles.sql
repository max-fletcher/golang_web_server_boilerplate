-- *IMPORTANT: SQLC generates type-safe code from SQL
-- You write queries in SQL.
-- You run sqlc to generate code with type-safe interfaces to those queries.
-- You write application code that calls the generated code.

-- *IMPORTANT: Goose is a database migration tool. Both a CLI and a library.
-- It is used to create schemas for easier migrations and database management.
-- To do database migration Up, go to the sql/schema folder and use "goose postgres postgres://postgres:@localhost:5432/golang_web_server up"
-- To do database migration Down, go to the sql/schema folder and use "goose postgres postgres://postgres:@localhost:5432/golang_web_server down"

-- +goose Up
CREATE TABLE user_roles(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    UNIQUE (user_id, role_id)
);

-- +goose Down
DROP TABLE user_roles;