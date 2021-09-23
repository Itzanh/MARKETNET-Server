ALTER TABLE public.sales_invoice
    ADD COLUMN simplified_invoice boolean NOT NULL DEFAULT false;