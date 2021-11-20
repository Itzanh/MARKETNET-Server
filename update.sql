CREATE TABLE public.transactional_log
(
    id bigint NOT NULL,
    enterprise integer NOT NULL,
    "table" character varying(150) COLLATE pg_catalog."default" NOT NULL,
    register jsonb NOT NULL,
    date_created timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT transactional_log_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.transactional_log
    OWNER to postgres;

GRANT ALL ON TABLE public.transactional_log TO marketnet;

GRANT ALL ON TABLE public.transactional_log TO postgres;

CREATE OR REPLACE FUNCTION set_transactional_log_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(transactional_log.id) END AS id FROM transactional_log) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_transactional_log_id
before insert on transactional_log
for each row execute procedure set_transactional_log_id();

ALTER TABLE public.transactional_log
    ADD COLUMN register_id bigint NOT NULL;
ALTER TABLE public.transactional_log
    ADD COLUMN "user" integer;

ALTER TABLE public.transactional_log
    ADD COLUMN mode character(1) NOT NULL;
ALTER TABLE public.transactional_log
    ADD CONSTRAINT transactional_log_user FOREIGN KEY ("user", enterprise)
    REFERENCES public."user" (id, config) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.config
    ADD COLUMN transaction_log boolean NOT NULL DEFAULT true;