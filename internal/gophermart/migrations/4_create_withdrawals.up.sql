CREATE TABLE IF NOT EXISTS withdrawals (
    id bigserial,
    "name" text NOT NULL,
    "order" text NOT NULL,
    "processed_at"  timestamptz,
    "withdraw" float8 NOT NULL DEFAULT 0,
    CONSTRAINT withdrawals_fk FOREIGN KEY (name) REFERENCES users("name")
);

CREATE UNIQUE INDEX IF NOT EXISTS withdrawals_user_order_idx ON withdrawals ("order", "name");
CREATE UNIQUE INDEX IF NOT EXISTS withdrawals_order_idx ON withdrawals ("order");