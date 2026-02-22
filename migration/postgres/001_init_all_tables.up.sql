CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL NOT NULL UNIQUE,
    guid UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

--bun:split

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL NOT NULL UNIQUE,
    guid UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT DEFAULT NULL,
    price DECIMAL(12,3) NOT NULL,
    category_guid UUID NOT NULL REFERENCES categories(guid) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(name, category_guid)
    );

CREATE INDEX idx_products_category ON products(category_guid);
CREATE INDEX idx_products_name ON products(name);