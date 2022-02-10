ALTER TABLE public.charges
    ALTER COLUMN concept TYPE character varying(140) COLLATE pg_catalog."default";

ALTER TABLE public.payments
    ALTER COLUMN concept TYPE character varying(140) COLLATE pg_catalog."default";