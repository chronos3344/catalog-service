CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL NOT NULL UNIQUE,
    guid UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

--bun:split

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL NOT NULL UNIQUE,
    guid UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    price DECIMAL(12,2) NOT NULL,
    category_guid UUID NOT NULL REFERENCES categories(guid) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
