ALTER TABLE public.sales_order_detail
    ADD COLUMN cancelled boolean NOT NULL DEFAULT false;