CREATE TABLE public.transfer_between_warehouses
(
    id bigint NOT NULL,
    warehouse_origin character(2) COLLATE pg_catalog."default" NOT NULL,
    warehouse_destination character(2) COLLATE pg_catalog."default" NOT NULL,
    enterprise integer NOT NULL,
    date_created timestamp(3) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    date_finished timestamp(3) without time zone,
    finished boolean NOT NULL DEFAULT false,
    lines_transfered integer NOT NULL,
    lines_total integer NOT NULL,
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT transfer_between_warehouses_pkey PRIMARY KEY (id),
    CONSTRAINT transfer_between_warehouses_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transfer_between_warehouses_warehouse_destination FOREIGN KEY (enterprise, warehouse_destination)
        REFERENCES public.warehouse (enterprise, id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transfer_between_warehouses_warehouse_origin FOREIGN KEY (enterprise, warehouse_origin)
        REFERENCES public.warehouse (enterprise, id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE INDEX transfer_between_warehouses_enterprise_finished_date_created
    ON public.transfer_between_warehouses USING btree
    (enterprise ASC NULLS LAST, finished ASC NULLS LAST, date_created ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE UNIQUE INDEX transfer_between_warehouses_id_enterprise
    ON public.transfer_between_warehouses USING btree
    (id ASC NULLS LAST, enterprise ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE TABLE public.transfer_between_warehouses_detail
(
    id bigint NOT NULL,
    transfer_between_warehouses bigint NOT NULL,
    enterprise integer NOT NULL,
    product integer NOT NULL,
    quantity integer NOT NULL,
    quantity_transfered integer NOT NULL DEFAULT 0,
    finished boolean NOT NULL DEFAULT false,
    CONSTRAINT transfer_between_warehouses_detail_pkey PRIMARY KEY (id),
    CONSTRAINT transfer_between_warehouses_detail_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transfer_between_warehouses_detail_product FOREIGN KEY (enterprise, product)
        REFERENCES public.product (enterprise, id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transfer_between_warehouses_detail_transfer_between_warehouses FOREIGN KEY (enterprise, transfer_between_warehouses)
        REFERENCES public.transfer_between_warehouses (enterprise, id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE OR REPLACE FUNCTION set_transfer_between_warehouses_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(transfer_between_warehouses.id) END AS id FROM transfer_between_warehouses) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_transfer_between_warehouses_id
before insert on transfer_between_warehouses
for each row execute procedure set_transfer_between_warehouses_id();

CREATE OR REPLACE FUNCTION set_transfer_between_warehouses_detail_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(transfer_between_warehouses_detail.id) END AS id FROM transfer_between_warehouses_detail) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_transfer_between_warehouses_detail_id
before insert on transfer_between_warehouses_detail
for each row execute procedure set_transfer_between_warehouses_detail_id();

CREATE INDEX transfer_between_warehouses_detail_barcode
    ON public.transfer_between_warehouses_detail USING btree
    (enterprise ASC NULLS LAST, transfer_between_warehouses ASC NULLS LAST, product ASC NULLS LAST)

    WHERE quantity_transfered < quantity;

ALTER TABLE public.transfer_between_warehouses
    ALTER COLUMN lines_total SET DEFAULT 0;

ALTER TABLE public.transfer_between_warehouses
    ALTER COLUMN lines_transfered SET DEFAULT 0;

ALTER TABLE public.transfer_between_warehouses_detail
    ADD COLUMN warehouse_movement_out bigint;

ALTER TABLE public.transfer_between_warehouses_detail
    ADD COLUMN warehouse_movement_in bigint;
ALTER TABLE public.transfer_between_warehouses_detail
    ADD CONSTRAINT transfer_between_warehouses_detail_warehouse_movement_out FOREIGN KEY (warehouse_movement_out)
    REFERENCES public.warehouse_movement (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

ALTER TABLE public.transfer_between_warehouses_detail
    ADD CONSTRAINT transfer_between_warehouses_detail_warehouse_movement_in FOREIGN KEY (warehouse_movement_in)
    REFERENCES public.warehouse_movement (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;