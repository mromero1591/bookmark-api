CREATE TABLE user_accounts (
    id   UUID PRIMARY KEY,
    email text NOT NULL UNIQUE,
    pwd_hash text NOT NULL,
    name text,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);