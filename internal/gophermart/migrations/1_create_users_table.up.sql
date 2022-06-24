create table if not exists users(
    id bigserial,
    name text NOT NULL UNIQUE,
    password text NOT NULL,
    CONSTRAINT users_id_pk PRIMARY KEY (id)
)