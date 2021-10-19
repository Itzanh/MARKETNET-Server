ALTER TABLE public.sales_invoice
    ADD COLUMN amending boolean NOT NULL DEFAULT false;

ALTER TABLE public.sales_invoice
    ADD COLUMN amended_invoice bigint;
ALTER TABLE public.sales_invoice
    ADD CONSTRAINT sales_invoice_amended_sales_invoice FOREIGN KEY (amended_invoice)
    REFERENCES public.sales_invoice (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.sales_invoice_detail
    ALTER COLUMN product DROP NOT NULL;

ALTER TABLE public.sales_invoice_detail
    ADD COLUMN description character varying(150) NOT NULL DEFAULT '';

ALTER TABLE public.purchase_invoice
    ADD COLUMN amending boolean NOT NULL DEFAULT false;

ALTER TABLE public.purchase_invoice
    ADD COLUMN amended_invoice bigint;
ALTER TABLE public.purchase_invoice
    ADD CONSTRAINT purchase_invoice_amended_purchase_invoice FOREIGN KEY (amended_invoice)
    REFERENCES public.purchase_invoice (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.purchase_invoice_details
    ALTER COLUMN product DROP NOT NULL;

ALTER TABLE public.purchase_invoice_details
    ADD COLUMN description character varying(150) NOT NULL DEFAULT '';