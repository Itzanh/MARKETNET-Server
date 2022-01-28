CREATE TABLE public.webhook_settings
(
    id integer NOT NULL,
    enterprise integer NOT NULL,
    url character varying(255) COLLATE pg_catalog."default" NOT NULL,
    auth_code uuid NOT NULL,
    auth_method character(1) COLLATE pg_catalog."default" NOT NULL,
    sale_orders boolean NOT NULL,
    sale_order_details boolean NOT NULL,
    sale_order_details_digital_product_data boolean NOT NULL,
    sale_invoices boolean NOT NULL,
    sale_invoice_details boolean NOT NULL,
    sale_delivery_notes boolean NOT NULL,
    purchase_orders boolean NOT NULL,
    purchase_order_details boolean NOT NULL,
    purchase_invoices boolean NOT NULL,
    purchase_invoice_details boolean NOT NULL,
    purchase_delivery_notes boolean NOT NULL,
    customers boolean NOT NULL,
    suppliers boolean NOT NULL,
    products boolean NOT NULL,
    CONSTRAINT webhook_settings_pkey PRIMARY KEY (id),
    CONSTRAINT webhook_settings_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE UNIQUE INDEX webhook_settings_id_enterprise
    ON public.webhook_settings USING btree
    (id ASC NULLS LAST, enterprise ASC NULLS LAST)
;

CREATE TABLE public.webhook_logs
(
    id bigint NOT NULL,
    webhook integer NOT NULL,
    enterprise integer NOT NULL,
    url character varying(255) COLLATE pg_catalog."default" NOT NULL,
    auth_code uuid NOT NULL,
    auth_method character(1) COLLATE pg_catalog."default" NOT NULL,
    date_created timestamp(3) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    sent text COLLATE pg_catalog."default" NOT NULL,
    received text COLLATE pg_catalog."default" NOT NULL,
    received_http_code smallint NOT NULL,
    method character varying(10) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT webhook_logs_pkey PRIMARY KEY (id),
    CONSTRAINT webhook_logs_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT webhook_logs_webhook FOREIGN KEY (enterprise, webhook)
        REFERENCES public.webhook_settings (enterprise, id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE TABLE public.webhook_queue
(
    id uuid NOT NULL,
    webhook integer NOT NULL,
    enterprise integer NOT NULL,
    url character varying(255) COLLATE pg_catalog."default" NOT NULL,
    auth_code uuid NOT NULL,
    auth_method character(1) COLLATE pg_catalog."default" NOT NULL,
    date_created timestamp(3) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    send text COLLATE pg_catalog."default" NOT NULL,
    method character varying(10) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT webhook_queue_pkey PRIMARY KEY (id),
    CONSTRAINT webhook_queue_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT webhook_queue_webhook FOREIGN KEY (enterprise, webhook)
        REFERENCES public.webhook_settings (enterprise, id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE OR REPLACE FUNCTION set_webhook_settings_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(webhook_settings.id) END AS id FROM webhook_settings) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_webhook_settings_id
before insert on webhook_settings
for each row execute procedure set_webhook_settings_id();

CREATE OR REPLACE FUNCTION set_webhook_logs_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(webhook_logs.id) END AS id FROM webhook_logs) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_webhook_logs_id
before insert on webhook_logs
for each row execute procedure set_webhook_logs_id();

ALTER TABLE public.inventory
    VALIDATE CONSTRAINT inventory_enterprise;

ALTER TABLE public.inventory
    VALIDATE CONSTRAINT inventory_warehouse;

ALTER TABLE public.report_template_translation
    VALIDATE CONSTRAINT report_template_translation_enterprise;

ALTER TABLE public.report_template_translation
    VALIDATE CONSTRAINT report_template_translation_language;

ALTER TABLE public.sales_order_discount
    VALIDATE CONSTRAINT sales_order_discount_sales_invoice_detail;