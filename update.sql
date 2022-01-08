ALTER TABLE public.sales_order_discount
    ADD COLUMN sales_invoice_detail integer;
ALTER TABLE public.sales_order_discount
    ADD CONSTRAINT sales_order_discount_sales_invoice_detail FOREIGN KEY (sales_invoice_detail, enterprise)
    REFERENCES public.sales_invoice_detail (id, enterprise) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;