CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    email TEXT UNIQUE,
    password TEXT,
    roles TEXT
);

CREATE INDEX idx_users_deleted_at
    ON users(deleted_at);


CREATE TABLE routes (
    id BIGSERIAL PRIMARY KEY,
    route TEXT NOT NULL,
    heading TEXT NOT NULL,
    role TEXT NOT NULL
);


CREATE TABLE billings (
    id BIGSERIAL PRIMARY KEY,
    "user" TEXT,
    mobile TEXT,
    grand_total INTEGER,
    finalized_amt INTEGER
);


CREATE TABLE purchases (
    id BIGSERIAL PRIMARY KEY,
    billing_id BIGINT NOT NULL,
    mobile TEXT,
    mrp_or_net INTEGER,
    item TEXT,
    quantity INTEGER,
    discount TEXT,
    sub_total INTEGER,

    CONSTRAINT fk_billings_purchases
        FOREIGN KEY (billing_id)
        REFERENCES billings(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE INDEX idx_purchases_billing_id
    ON purchases(billing_id);

CREATE INDEX idx_purchases_mobile
    ON purchases(mobile);


CREATE TABLE inventories (
    id BIGSERIAL PRIMARY KEY,
    brand_or_company TEXT,
    state VARCHAR(20),
    item TEXT,
    num_of_boxes INTEGER,
    num_of_cartons INTEGER,
    price_per_carton INTEGER,
    sub_total INTEGER
);

CREATE INDEX idx_inventories_item
    ON inventories(item);

CREATE INDEX idx_inventories_state
    ON inventories(state);


CREATE TABLE expenses (
    id BIGSERIAL PRIMARY KEY,
    reason_for_expense TEXT,
    amount INTEGER
);


CREATE TABLE price_lists (
    id BIGSERIAL PRIMARY KEY,
    price BIGINT,
    item TEXT,
    brand TEXT
);