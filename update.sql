ALTER TABLE public.product
    ADD COLUMN purchase_price real NOT NULL DEFAULT 0;

UPDATE product SET purchase_price = price WHERE manufacturing = false;