CREATE TABLE public.manufacturing_order_type_components
(
    id integer NOT NULL,
    manufacturing_order_type integer NOT NULL,
    type character(1) COLLATE pg_catalog."default" NOT NULL,
    product integer NOT NULL,
    quantity integer NOT NULL,
    enterprise integer NOT NULL,
    CONSTRAINT manufacturing_order_type_components_pkey PRIMARY KEY (id),
    CONSTRAINT manufacturing_order_type_components_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT manufacturing_order_type_components_manufacturing_order_type FOREIGN KEY (enterprise, manufacturing_order_type)
        REFERENCES public.manufacturing_order_type (enterprise, id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT manufacturing_order_type_components_product FOREIGN KEY (enterprise, product)
        REFERENCES public.product (enterprise, id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE public.manufacturing_order_type_components
    OWNER to postgres;

GRANT ALL ON TABLE public.manufacturing_order_type_components TO postgres;

CREATE OR REPLACE FUNCTION set_manufacturing_order_type_components_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(manufacturing_order_type_components.id) END AS id FROM manufacturing_order_type_components) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_manufacturing_order_type_components_id
before insert on manufacturing_order_type_components
for each row execute procedure set_manufacturing_order_type_components_id();

ALTER TABLE public.manufacturing_order_type
    ADD COLUMN complex boolean NOT NULL DEFAULT false;

CREATE UNIQUE INDEX manufacturing_order_type_components_manufacturing_order_type_type_product
    ON public.manufacturing_order_type_components USING btree
    (manufacturing_order_type ASC NULLS LAST, type ASC NULLS LAST, product ASC NULLS LAST)
;

ALTER TABLE public.manufacturing_order
    ALTER COLUMN user_manufactured TYPE integer;

CREATE TABLE public.complex_manufacturing_order
(
    id bigint NOT NULL,
    type integer NOT NULL,
    manufactured boolean NOT NULL DEFAULT false,
    date_manufactured timestamp(3) with time zone,
    user_manufactured integer,
    enterprise integer NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT complex_manufacturing_order_type FOREIGN KEY (type, enterprise)
        REFERENCES public.manufacturing_order_type (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT complex_manufacturing_order_user_manufactured FOREIGN KEY (user_manufactured, enterprise)
        REFERENCES public."user" (id, config) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT complex_manufacturing_order_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE public.complex_manufacturing_order
    OWNER to postgres;

ALTER TABLE public.complex_manufacturing_order
    ADD COLUMN quantity_pending_manufacture integer NOT NULL DEFAULT 0;

ALTER TABLE public.complex_manufacturing_order
    ADD COLUMN quantity_manufactured integer NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX manufacturing_order_id_enterprise
    ON public.manufacturing_order USING btree
    (id ASC NULLS LAST, enterprise ASC NULLS LAST)
;

CREATE UNIQUE INDEX complex_manufacturing_order_id_enterprise
    ON public.complex_manufacturing_order USING btree
    (id ASC NULLS LAST, enterprise ASC NULLS LAST)
;

CREATE TABLE public.complex_manufacturing_order_manufacturing_order
(
    id bigint NOT NULL,
    manufacturing_order bigint NOT NULL,
    type character(1) NOT NULL,
    complex_manufacturing_order bigint NOT NULL,
    enterprise integer NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT complex_manufacturing_order_manufacturing_order FOREIGN KEY (manufacturing_order, enterprise)
        REFERENCES public.manufacturing_order (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT complex_manufacturing_order_complex_manufacturing_order FOREIGN KEY (complex_manufacturing_order, enterprise)
        REFERENCES public.complex_manufacturing_order (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT complex_manufacturing_order_manufacturing_order_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    OWNER to postgres;

CREATE UNIQUE INDEX complex_manufacturing_order_complex_manufacturing_order_manufac
    ON public.complex_manufacturing_order_manufacturing_order USING btree
    (complex_manufacturing_order ASC NULLS LAST, manufacturing_order ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE INDEX complex_manufacturing_order_complex_manufacturing_order_type
    ON public.complex_manufacturing_order_manufacturing_order USING btree
    (complex_manufacturing_order ASC NULLS LAST, type ASC NULLS LAST)
;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ALTER COLUMN manufacturing_order DROP NOT NULL;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD COLUMN warehouse_movement bigint;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD COLUMN manufactured boolean NOT NULL DEFAULT false;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD COLUMN product integer NOT NULL;

ALTER TABLE public.complex_manufacturing_order
    ADD COLUMN warehouse character(2) NOT NULL;
ALTER TABLE public.complex_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_warehouse FOREIGN KEY (warehouse, enterprise)
    REFERENCES public.warehouse (id, enterprise) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.complex_manufacturing_order
    VALIDATE CONSTRAINT complex_manufacturing_order_warehouse;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_warehouse_movement FOREIGN KEY (warehouse_movement)
    REFERENCES public.warehouse_movement (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_product FOREIGN KEY (product, enterprise)
    REFERENCES public.product (id, enterprise) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

CREATE UNIQUE INDEX manufacturing_order_type_components_id_enterprise
    ON public.manufacturing_order_type_components USING btree
    (id ASC NULLS LAST, enterprise ASC NULLS LAST)
;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD COLUMN manufacturing_order_type_component integer NOT NULL;
ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_manufacturing_order_type_components FOREIGN KEY (manufacturing_order_type_component, enterprise)
    REFERENCES public.manufacturing_order_type_components (id, enterprise) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD COLUMN purchase_order_detail integer;
ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_purchase_order_detail FOREIGN KEY (purchase_order_detail, enterprise)
    REFERENCES public.purchase_order_detail (id, enterprise) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.manufacturing_order
    ADD COLUMN complex boolean NOT NULL DEFAULT false;

CREATE INDEX manufacturing_order_for_stock_pending
    ON public.manufacturing_order USING btree
    (enterprise ASC NULLS LAST, product ASC NULLS LAST, manufactured ASC NULLS LAST, order_detail ASC NULLS LAST, complex ASC NULLS LAST)

    WHERE NOT manufactured AND order_detail IS NULL AND NOT complex;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ALTER COLUMN purchase_order_detail TYPE bigint;

CREATE OR REPLACE FUNCTION set_complex_manufacturing_order_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(complex_manufacturing_order.id) END AS id FROM complex_manufacturing_order) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_complex_manufacturing_order_id
before insert on complex_manufacturing_order
for each row execute procedure set_complex_manufacturing_order_id();

CREATE OR REPLACE FUNCTION set_complex_manufacturing_order_manufacturing_order_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(complex_manufacturing_order_manufacturing_order.id) END AS id FROM complex_manufacturing_order_manufacturing_order) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_complex_manufacturing_order_manufacturing_order_id
before insert on complex_manufacturing_order_manufacturing_order
for each row execute procedure set_complex_manufacturing_order_manufacturing_order_id();

ALTER TABLE public.complex_manufacturing_order
    ADD COLUMN date_created timestamp(3) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

ALTER TABLE public.complex_manufacturing_order
    ADD COLUMN uuid uuid NOT NULL;

ALTER TABLE public.complex_manufacturing_order
    ADD COLUMN user_created integer NOT NULL;

ALTER TABLE public.complex_manufacturing_order
    ADD COLUMN tag_printed boolean NOT NULL DEFAULT false;

ALTER TABLE public.complex_manufacturing_order
    ADD COLUMN date_tag_printed timestamp(3) with time zone;

ALTER TABLE public.complex_manufacturing_order
    ADD COLUMN user_tag_printed integer;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD COLUMN sale_order_detail bigint;
ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_sale_order_detail FOREIGN KEY (sale_order_detail, enterprise)
    REFERENCES public.sales_order_detail (id, enterprise) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

CREATE INDEX complex_manufacturing_order_manufacturing_order_pending_product_output_no_sale
    ON public.complex_manufacturing_order_manufacturing_order USING btree
    (product ASC NULLS LAST, manufactured ASC NULLS LAST, type ASC NULLS LAST, sale_order_detail ASC NULLS LAST)

    WHERE (NOT manufactured) AND (type = 'O') AND (sale_order_detail IS NULL);