ALTER TABLE public.config
    ADD COLUMN email_send_error_ecommerce character varying(150) NOT NULL DEFAULT '';

ALTER TABLE public.config
    ADD COLUMN email_send_error_sendcloud character varying(150) NOT NULL DEFAULT '';