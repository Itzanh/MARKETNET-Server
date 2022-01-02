ALTER TABLE public.sales_order
    ALTER COLUMN reference TYPE character varying(15) COLLATE pg_catalog."default";

CREATE UNIQUE INDEX sales_invoice_detail_invoice_product
    ON public.sales_invoice_detail USING btree
    (invoice ASC NULLS LAST, product ASC NULLS LAST)
;