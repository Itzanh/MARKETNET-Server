ALTER TABLE public.manufacturing_order
    ADD COLUMN warehouse character(2) NOT NULL DEFAULT 'W1';

ALTER TABLE public.manufacturing_order
    ADD COLUMN warehouse_movement bigint;

ALTER TABLE public.manufacturing_order
    ADD COLUMN quantity_manufactured integer NOT NULL DEFAULT 1;
ALTER TABLE public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_warehouse FOREIGN KEY (warehouse, enterprise)
    REFERENCES public.warehouse (id, enterprise) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_warehouse_movement FOREIGN KEY (warehouse_movement)
    REFERENCES public.warehouse_movement (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.manufacturing_order_type
    ADD COLUMN quantity_manufactured integer NOT NULL DEFAULT 1;