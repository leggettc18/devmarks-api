CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    email text UNIQUE NOT NULL,
    hashed_password bytea NOT NULL,
    name text,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);