ALTER TABLE public.purchase_order_detail
    ADD COLUMN cancelled boolean NOT NULL DEFAULT false;