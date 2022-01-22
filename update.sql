CREATE TABLE public.hs_codes
(
    id character varying(8) COLLATE pg_catalog."default" NOT NULL,
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT hs_codes_pkey PRIMARY KEY (id)
);

ALTER TABLE public.product
    ADD COLUMN origin_country character varying(2) NOT NULL DEFAULT '';

ALTER TABLE public.product
    ADD COLUMN hs_code character varying(8);