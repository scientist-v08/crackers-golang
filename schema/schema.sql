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

CREATE TABLE IF NOT EXISTS temp_billings (
    id              SMALLSERIAL PRIMARY KEY,
    "user"          TEXT        NOT NULL DEFAULT '',
    mobile          TEXT        NOT NULL DEFAULT '',
    comments        TEXT        NOT NULL DEFAULT '',
    grand_total     DOUBLE PRECISION NOT NULL DEFAULT 0,
    finalized_amt   DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS temp_purchases (
    id              SMALLSERIAL PRIMARY KEY,
    billing_id      SMALLINT    NOT NULL
                        REFERENCES temp_billings(id)
                        ON UPDATE CASCADE
                        ON DELETE CASCADE,
    mobile          TEXT        NOT NULL DEFAULT '',
    mrp_or_net      INTEGER     NOT NULL DEFAULT 0,
    item            TEXT        NOT NULL DEFAULT '',
    quantity        INTEGER     NOT NULL DEFAULT 0,
    discount        TEXT        NOT NULL DEFAULT '',
    sub_total       DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_temp_purchases_billing_id ON temp_purchases (billing_id);
CREATE INDEX IF NOT EXISTS idx_temp_purchases_mobile     ON temp_purchases (mobile);