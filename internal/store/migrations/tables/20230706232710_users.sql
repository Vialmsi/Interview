-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id SERIAL NOT NULL,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT,
    email TEXT,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users
-- +goose StatementEnd
