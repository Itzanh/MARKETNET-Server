ALTER TABLE public.carrier
    ADD COLUMN sendcloud_url character varying(75) NOT NULL DEFAULT '';

ALTER TABLE public.carrier
    ADD COLUMN sendcloud_key character varying(32) NOT NULL DEFAULT '';

ALTER TABLE public.carrier
    ADD COLUMN sendcloud_secret character varying(32) NOT NULL DEFAULT '';

ALTER TABLE public.carrier
    ADD COLUMN sendcloud_shipping_method integer NOT NULL DEFAULT 0;

ALTER TABLE public.carrier
    ADD COLUMN sendcloud_sender_address bigint NOT NULL DEFAULT 0;

ALTER TABLE public.config
    ADD COLUMN cron_clear_labels character varying(25) NOT NULL DEFAULT '@midnight';