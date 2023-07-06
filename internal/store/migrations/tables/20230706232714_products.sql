-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products
(
    barcode TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    cost INT NOT NULL,
    user_id INT,
    file_location TEXT DEFAULT 'empty',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT products_pkey PRIMARY KEY (barcode),
    CONSTRAINT users_fkey FOREIGN KEY (user_id)
    REFERENCES users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products
-- +goose StatementEnd
