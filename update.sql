DROP INDEX public.product_reference;

CREATE INDEX product_reference
    ON public.product USING gin
    (reference COLLATE pg_catalog."default" gin_trgm_ops)
    TABLESPACE pg_default
    WHERE reference::text <> ''::text;

ALTER TABLE public.warehouse_movement DROP COLUMN purchase_invoice_details;
ALTER TABLE public.warehouse_movement DROP COLUMN sales_invoice_detail;
ALTER TABLE public.warehouse_movement DROP COLUMN sales_invoice;
ALTER TABLE public.warehouse_movement DROP COLUMN purchase_invoice;