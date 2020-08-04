CREATE TABLE IF NOT EXISTS bookmarks(
    id serial PRIMARY KEY,
    name text NOT NULL,
    url text NOT NULL,
    color text,
    owner_id int NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,

    CONSTRAINT bookmarks_owner_id_fkey FOREIGN KEY (owner_id)
    REFERENCES users(id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE CASCADE
);