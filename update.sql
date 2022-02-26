DROP INDEX public.currency_sign;
ALTER TABLE public.currency
    ALTER COLUMN name TYPE character varying(150) COLLATE pg_catalog."default";
