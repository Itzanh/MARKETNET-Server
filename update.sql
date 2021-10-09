ALTER TABLE public.stock DROP CONSTRAINT stock_warehouse;

ALTER TABLE public.stock
    ADD CONSTRAINT stock_warehouse FOREIGN KEY (warehouse, enterprise)
    REFERENCES public.warehouse (id, enterprise) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;

ALTER TABLE public.stock DROP CONSTRAINT stock_product;

ALTER TABLE public.stock
    ADD CONSTRAINT stock_product FOREIGN KEY (product, enterprise)
    REFERENCES public.product (id, enterprise) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;
