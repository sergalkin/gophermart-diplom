CREATE TABLE IF NOT EXISTS balance (
    id bigserial,
    "name" text NOT NULL UNIQUE,
    "balance" float8 NOT NULL DEFAULT 0,
    "withdraw" float8 NOT NULL DEFAULT 0,
    CONSTRAINT balance_fk FOREIGN KEY (name) REFERENCES public.users("name"),
    CONSTRAINT balance_id_pk PRIMARY KEY (id)
);