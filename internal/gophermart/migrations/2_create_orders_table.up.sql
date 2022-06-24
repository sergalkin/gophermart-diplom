CREATE TABLE IF NOT EXISTS orders (
    id bigserial,
    "order" text NOT NULL,
    "name" text NOT NULL,
    "status" text NOT NULL DEFAULT 'NEW',
    "uploaded_at" timestamptz NOT NULL,
    "accrual" float8 NOT NULL DEFAULT 0,
    CONSTRAINT orders_fk FOREIGN KEY (name) REFERENCES public.users("name"),
    CONSTRAINT orders_id_pk PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS orders_order_user_idx ON public.orders ("order","name");
CREATE UNIQUE INDEX IF NOT EXISTS orders_order_idx ON public.orders ("order");