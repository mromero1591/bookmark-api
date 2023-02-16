CREATE TABLE user_accounts (
    id   UUID PRIMARY KEY,
    email text NOT NULL UNIQUE,
    pwd_hash text NOT NULL,
    name text,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE category (
    id UUID PRIMARY KEY,
    name text NOT NULL,
    user_id UUID NOT NULL REFERENCES user_accounts(id),
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE bookmark (
    id UUID PRIMARY KEY,
    url text NOT NULL,
    name text NOT NULL,
    logo text,
    category_id UUID NOT NULL REFERENCES category(id),
    user_id UUID NOT NULL REFERENCES user_accounts(id),
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);