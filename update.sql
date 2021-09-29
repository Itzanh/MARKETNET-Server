ALTER TABLE public.config
    ADD COLUMN password_minimum_length smallint NOT NULL DEFAULT 8;

ALTER TABLE public.config
    ADD COLUMN password_minumum_complexity character(1) NOT NULL DEFAULT 'B';

CREATE TABLE public.pwd_blacklist
(
    pwd character varying(255) NOT NULL,
    PRIMARY KEY (pwd)
);

CREATE TABLE public.pwd_sha1_blacklist
(
    hash bytea NOT NULL,
    PRIMARY KEY (hash)
);
