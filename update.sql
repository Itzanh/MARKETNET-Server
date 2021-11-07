ALTER TABLE public.config
    ADD COLUMN invoice_delete_policy smallint NOT NULL DEFAULT 1;