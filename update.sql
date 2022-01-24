CREATE TABLE public.inventory
(
    id integer NOT NULL,
    enterprise integer NOT NULL,
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    date_created timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    finished boolean NOT NULL DEFAULT false,
    date_finished timestamp(3) without time zone,
    warehouse character(2) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT inventory_pkey PRIMARY KEY (id),
    CONSTRAINT inventory_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT inventory_warehouse FOREIGN KEY (enterprise, warehouse)
        REFERENCES public.warehouse (enterprise, id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TABLE public.inventory_products
(
    inventory integer NOT NULL,
    product integer NOT NULL,
    enterprise integer NOT NULL,
    quantity integer NOT NULL,
    warehouse_movement bigint,
    CONSTRAINT inventory_products_pkey PRIMARY KEY (inventory, product),
    CONSTRAINT inventory_products_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT inventory_products_inventory FOREIGN KEY (inventory)
        REFERENCES public.inventory (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT inventory_products_product FOREIGN KEY (product)
        REFERENCES public.product (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT inventory_products_warehouse_movement FOREIGN KEY (warehouse_movement)
        REFERENCES public.warehouse_movement (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE OR REPLACE FUNCTION set_inventory_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(inventory.id) END AS id FROM inventory) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_inventory_id
before insert on inventory
for each row execute procedure set_inventory_id();