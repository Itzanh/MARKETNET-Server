CREATE TABLE public.email_log
(
    id bigint NOT NULL,
    email_from character varying(100) NOT NULL,
    name_from character varying(100) NOT NULL,
    destination_email character varying(100) NOT NULL,
    destination_name character varying(100) NOT NULL,
    subject character varying(100) NOT NULL,
    content text NOT NULL,
    date_sent timestamp(3) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    enterprise integer NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT email_log_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE public.email_log
    OWNER to postgres;

CREATE INDEX email_log_date_sent
    ON public.email_log USING btree
    (date_sent DESC NULLS LAST)
;

CREATE INDEX email_log_trgm
    ON public.email_log USING gin
    (email_from gin_trgm_ops, name_from gin_trgm_ops, destination_email gin_trgm_ops, destination_name gin_trgm_ops, subject gin_trgm_ops)
;

CREATE OR REPLACE FUNCTION set_email_log_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(email_log.id) END AS id FROM email_log) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_email_log_id
before insert on email_log
for each row execute procedure set_email_log_id();