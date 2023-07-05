CREATE TABLE IF NOT EXISTS users
(
    id serial int NOT NULL,
    login text NOT NULL UNIQUE,
    password text NOT NULL,
    name text,
    email text,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)