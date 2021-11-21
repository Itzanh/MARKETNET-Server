CREATE TABLE public.shipping_status_history
(
    id integer NOT NULL,
    shipping integer NOT NULL,
    status_id smallint NOT NULL,
    message character varying(255) NOT NULL,
    delivered boolean NOT NULL,
    date_created timestamp(3) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (shipping)
);

ALTER TABLE public.shipping_status_history
    OWNER to postgres;

ALTER TABLE public.shipping
    ADD COLUMN delivered boolean NOT NULL DEFAULT false;

CREATE INDEX shipping_sent_collected_delivered
    ON public.shipping USING btree
    (sent ASC NULLS LAST, collected ASC NULLS LAST, delivered ASC NULLS LAST)
;

ALTER TABLE public.shipping_status_history
    ALTER COLUMN id TYPE bigint;

ALTER TABLE public.shipping_status_history
    ALTER COLUMN shipping TYPE bigint;

CREATE OR REPLACE FUNCTION set_shipping_status_history_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(shipping_status_history.id) END AS id FROM shipping_status_history) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_shipping_status_history_id
before insert on shipping_status_history
for each row execute procedure set_shipping_status_history_id();

ALTER TABLE public.config
    ADD COLUMN cron_sendcloud_tracking character varying(25) NOT NULL DEFAULT '';