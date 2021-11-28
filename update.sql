ALTER TABLE public.config
    ADD COLUMN smtp_identity character varying(50) NOT NULL DEFAULT '';

ALTER TABLE public.config
    ADD COLUMN smtp_username character varying(50) NOT NULL DEFAULT '';

ALTER TABLE public.config
    ADD COLUMN smtp_password character varying(50) NOT NULL DEFAULT '';

ALTER TABLE public.config
    ADD COLUMN smtp_hostname character varying(50) NOT NULL DEFAULT '';