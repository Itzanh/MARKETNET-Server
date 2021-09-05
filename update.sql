ALTER TABLE public.config
    ADD COLUMN woocommerce_url character varying(100) NOT NULL DEFAULT '';

ALTER TABLE public.config
    ADD COLUMN woocommerce_consumer_key character varying(50) NOT NULL DEFAULT '';

ALTER TABLE public.config
    ADD COLUMN woocommerce_consumer_secret character varying(50) NOT NULL DEFAULT '';

ALTER TABLE public.config
    ADD COLUMN woocommerce_export_serie character(3);

ALTER TABLE public.config
    ADD COLUMN woocommerce_intracommunity_serie character(3);

ALTER TABLE public.config
    ADD COLUMN woocommerce_interior_serie character(3);

ALTER TABLE public.config
    ADD CONSTRAINT config_woocommerce_export_serie FOREIGN KEY (woocommerce_export_serie)
    REFERENCES public.billing_series (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.config
    ADD CONSTRAINT config_woocommerce_intracommunity_serie FOREIGN KEY (woocommerce_intracommunity_serie)
    REFERENCES public.billing_series (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.config
    ADD CONSTRAINT config_woocommerce_interior_serie FOREIGN KEY (woocommerce_interior_serie)
    REFERENCES public.billing_series (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.customer
    ADD COLUMN wc_id integer NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX customer_wc_id
    ON public.customer USING btree
    (wc_id ASC NULLS LAST)

    WHERE wc_id != 0;

ALTER TABLE public.product
    ADD COLUMN wc_id integer NOT NULL DEFAULT 0;

ALTER TABLE public.product
    ADD COLUMN wc_variation_id integer NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX products_wc_id
    ON public.product USING btree
    (wc_id ASC NULLS LAST, wc_variation_id ASC NULLS LAST)

    WHERE wc_id != 0;

ALTER TABLE public.payment_method
    ADD COLUMN woocommerce_module_name character varying(100) NOT NULL DEFAULT '';

ALTER TABLE public.sales_order
    ADD COLUMN wc_id integer NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX sales_order_wc_id
    ON public.sales_order USING btree
    (wc_id ASC NULLS LAST)

    WHERE wc_id != 0;

ALTER TABLE public.sales_order_detail
    ADD COLUMN wc_id integer NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX sales_order_detail_wc_id
    ON public.sales_order_detail USING btree
    (wc_id ASC NULLS LAST)

    WHERE wc_id != 0;

ALTER TABLE public.config
    ADD COLUMN woocommerce_default_payment_method smallint;
ALTER TABLE public.config
    ADD CONSTRAINT config_woocommerce_default_payment_method FOREIGN KEY (woocommerce_default_payment_method)
    REFERENCES public.payment_method (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

CREATE TABLE public.wc_customers
(
    id integer NOT NULL,
    date_created timestamp(0) without time zone NOT NULL,
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    first_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    last_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_address_1 character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_address_2 character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_city character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_postcode character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_country character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_state character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_phone character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_1 character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_2 character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_city character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_postcode character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_country character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_state character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_phone character varying(255) COLLATE pg_catalog."default" NOT NULL,
    wc_exists boolean NOT NULL DEFAULT true,
    billing_company character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT ws_customers_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.wc_customers
    OWNER to postgres;

CREATE TABLE public.wc_order_details
(
    id integer NOT NULL,
    "order" integer NOT NULL,
    product_id integer NOT NULL,
    variation_id integer NOT NULL,
    quantity integer NOT NULL,
    total_tax real NOT NULL,
    price real NOT NULL,
    CONSTRAINT wc_order_details_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.wc_order_details
    OWNER to postgres;

CREATE TABLE public.wc_orders
(
    id integer NOT NULL,
    status character varying(50) COLLATE pg_catalog."default" NOT NULL,
    currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
    date_created timestamp(0) without time zone NOT NULL,
    discount_tax real NOT NULL,
    shipping_total real NOT NULL,
    shipping_tax real NOT NULL,
    total_tax real NOT NULL,
    customer_id integer NOT NULL,
    order_key character varying(25) COLLATE pg_catalog."default" NOT NULL,
    billing_address_1 character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_address_2 character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_city character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_postcode character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_country character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_state character varying(255) COLLATE pg_catalog."default" NOT NULL,
    billing_phone character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_1 character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_2 character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_city character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_postcode character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_country character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_state character varying(255) COLLATE pg_catalog."default" NOT NULL,
    shipping_phone character varying(255) COLLATE pg_catalog."default" NOT NULL,
    payment_method character varying(50) COLLATE pg_catalog."default" NOT NULL,
    wc_exists boolean NOT NULL DEFAULT true,
    billing_company character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT wc_orders_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.wc_orders
    OWNER to postgres;

CREATE TABLE public.wc_product_variations
(
    id integer NOT NULL,
    sku character varying(25) COLLATE pg_catalog."default" NOT NULL,
    price real NOT NULL,
    weight character varying(10) COLLATE pg_catalog."default" NOT NULL,
    dimensions_length character varying(10) COLLATE pg_catalog."default" NOT NULL,
    dimensions_width character varying(10) COLLATE pg_catalog."default" NOT NULL,
    dimensions_height character varying(10) COLLATE pg_catalog."default" NOT NULL,
    attributes character varying(255)[] COLLATE pg_catalog."default" NOT NULL,
    wc_exists boolean NOT NULL DEFAULT true,
    CONSTRAINT wc_product_variations_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.wc_product_variations
    OWNER to postgres;

CREATE TABLE public.wc_products
(
    id integer NOT NULL,
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    date_created timestamp(0) without time zone NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    short_description character varying(255) COLLATE pg_catalog."default" NOT NULL,
    sku character varying(25) COLLATE pg_catalog."default" NOT NULL,
    price real NOT NULL,
    weight character varying(10) COLLATE pg_catalog."default" NOT NULL,
    dimensions_length character varying(10) COLLATE pg_catalog."default" NOT NULL,
    dimensions_width character varying(10) COLLATE pg_catalog."default" NOT NULL,
    dimensions_height character varying(10) COLLATE pg_catalog."default" NOT NULL,
    images character varying(255)[] COLLATE pg_catalog."default" NOT NULL,
    wc_exists boolean NOT NULL DEFAULT true,
    variations integer[] NOT NULL,
    CONSTRAINT wc_products_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.wc_products
    OWNER to postgres;