ALTER TABLE public.purchase_invoice
    ADD COLUMN income_tax boolean NOT NULL DEFAULT false;

ALTER TABLE public.purchase_invoice
    ADD COLUMN income_tax_base real NOT NULL DEFAULT 0;

ALTER TABLE public.purchase_invoice
    ADD COLUMN income_tax_percentage real NOT NULL DEFAULT 0;

ALTER TABLE public.purchase_invoice
    ADD COLUMN income_tax_value real NOT NULL DEFAULT 0;

ALTER TABLE public.purchase_invoice
    ADD COLUMN rent boolean NOT NULL DEFAULT false;

ALTER TABLE public.purchase_invoice
    ADD COLUMN rent_base real NOT NULL DEFAULT 0;

ALTER TABLE public.purchase_invoice
    ADD COLUMN rent_percentage real NOT NULL DEFAULT 0;

ALTER TABLE public.purchase_invoice
    ADD COLUMN rent_value real NOT NULL DEFAULT 0;

ALTER TABLE public.purchase_invoice_details
    ADD COLUMN income_tax boolean NOT NULL DEFAULT false;

ALTER TABLE public.purchase_invoice_details
    ADD COLUMN rent boolean NOT NULL DEFAULT false;