-- up
CREATE TABLE  users (
    id SERIAL PRIMARY KEY,
    username TEXT  NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT  NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
