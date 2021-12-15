CREATE TABLE public.pos_terminals
(
    id bigint NOT NULL,
    uuid uuid NOT NULL,
    name character varying(150) NOT NULL,
    orders_customer integer,
    orders_invoice_address integer,
    orders_delivery_address integer,
    orders_payment_method integer,
    orders_billing_series character(3),
    orders_warehouse character(2),
    orders_currency integer,
    enterprise integer NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT pos_terminals_orders_customer FOREIGN KEY (orders_customer, enterprise)
        REFERENCES public.customer (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT pos_terminals_orders_invoice_address FOREIGN KEY (orders_invoice_address, enterprise)
        REFERENCES public.address (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT pos_terminals_orders_delivery_address FOREIGN KEY (orders_delivery_address, enterprise)
        REFERENCES public.address (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT pos_terminals_orders_payment_method FOREIGN KEY (orders_payment_method, enterprise)
        REFERENCES public.payment_method (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT pos_terminals_orders_billing_series FOREIGN KEY (orders_billing_series, enterprise)
        REFERENCES public.billing_series (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT pos_terminals_orders_warehouse FOREIGN KEY (orders_warehouse, enterprise)
        REFERENCES public.warehouse (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT pos_terminals_orders_currency FOREIGN KEY (orders_currency, enterprise)
        REFERENCES public.currency (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT pos_terminals_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE UNIQUE INDEX pos_terminals_uuid
    ON public.pos_terminals USING btree
    (uuid ASC NULLS LAST)
;

CREATE OR REPLACE FUNCTION set_pos_terminals_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(pos_terminals.id) END AS id FROM pos_terminals) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_pos_terminals_id
before insert on pos_terminals
for each row execute procedure set_pos_terminals_id();

ALTER TABLE public."group"
    ADD COLUMN point_of_sale boolean NOT NULL DEFAULT false;