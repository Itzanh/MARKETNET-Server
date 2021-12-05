CREATE TABLE public.enterprise_logo
(
    enterprise integer NOT NULL,
    logo bytea NOT NULL,
    PRIMARY KEY (enterprise)
);

ALTER TABLE public.enterprise_logo
    OWNER to postgres;

ALTER TABLE public.enterprise_logo
    ADD CONSTRAINT enterprise_logo_enterprise FOREIGN KEY (enterprise)
    REFERENCES public.config (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.enterprise_logo
    VALIDATE CONSTRAINT enterprise_logo_enterprise;

ALTER TABLE public.manufacturing_order
    VALIDATE CONSTRAINT manufacturing_order_warehouse;

ALTER TABLE public.manufacturing_order
    VALIDATE CONSTRAINT manufacturing_order_warehouse_movement;

ALTER TABLE public.purchase_invoice
    VALIDATE CONSTRAINT purchase_invoice_amended_purchase_invoice;

ALTER TABLE public.sales_invoice
    VALIDATE CONSTRAINT sales_invoice_amended_sales_invoice;

ALTER TABLE public.sales_order_detail_digital_product_data
    VALIDATE CONSTRAINT sales_order_detail_digital_product_data_sales_order_detail;

ALTER TABLE public.shipping_status_history
    ADD CONSTRAINT shipping_status_history_shipping FOREIGN KEY (shipping)
    REFERENCES public.shipping (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

ALTER TABLE public.transactional_log
    VALIDATE CONSTRAINT transactional_log_user;

ALTER TABLE public.enterprise_logo
    ADD COLUMN mime_type character varying(150) NOT NULL;