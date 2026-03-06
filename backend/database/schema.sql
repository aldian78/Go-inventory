-- Create Database (Optional, if not exists)
-- CREATE DATABASE inventory_db;

-- Table: products
CREATE TABLE IF NOT EXISTS products (
    id bigserial PRIMARY KEY,
    created_at timestamptz,
    updated_at timestamptz,
    deleted_at timestamptz,
    name text,
    sku text UNIQUE,
    customer text,
    physical_stock bigint,
    available_stock bigint
);

CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);

-- Table: stock_ins
CREATE TABLE IF NOT EXISTS stock_ins (
    id bigserial PRIMARY KEY,
    created_at timestamptz,
    updated_at timestamptz,
    deleted_at timestamptz,
    product_id bigint,
    quantity bigint,
    status text DEFAULT 'CREATED',
    notes text,
    CONSTRAINT fk_stock_ins_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_stock_ins_deleted_at ON stock_ins(deleted_at);

-- Table: stock_outs
CREATE TABLE IF NOT EXISTS stock_outs (
    id bigserial PRIMARY KEY,
    created_at timestamptz,
    updated_at timestamptz,
    deleted_at timestamptz,
    product_id bigint,
    quantity bigint,
    status text DEFAULT 'DRAFT',
    notes text,
    CONSTRAINT fk_stock_outs_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_stock_outs_deleted_at ON stock_outs(deleted_at);

-- Table: stock_logs
CREATE TABLE IF NOT EXISTS stock_logs (
    id bigserial PRIMARY KEY,
    created_at timestamptz,
    updated_at timestamptz,
    deleted_at timestamptz,
    product_id bigint,
    transaction_id bigint,
    transaction_type text,
    quantity bigint,
    previous_stock bigint,
    current_stock bigint,
    notes text,
    CONSTRAINT fk_stock_logs_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_stock_logs_deleted_at ON stock_logs(deleted_at);
