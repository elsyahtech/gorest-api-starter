CREATE TABLE IF NOT EXISTS users (
    id text PRIMARY KEY,
    email text,
    password text,
    is_active boolean,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
