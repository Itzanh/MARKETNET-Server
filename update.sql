ALTER TABLE public.config
    ADD COLUMN shopify_url character varying(100) NOT NULL DEFAULT '';

ALTER TABLE public.config
    ADD COLUMN shopify_token character varying(50) NOT NULL DEFAULT '';

ALTER TABLE public.customer
    ADD COLUMN sy_id bigint NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX customer_sy_id
    ON public.customer USING btree
    (sy_id ASC NULLS LAST)

    WHERE sy_id != 0;

ALTER TABLE public.address
    ADD COLUMN sy_id bigint NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX address_sy_id
    ON public.address USING btree
    (sy_id ASC NULLS LAST)

    WHERE sy_id != 0;

ALTER TABLE public.product
    ADD COLUMN sy_id bigint NOT NULL DEFAULT 0;

ALTER TABLE public.product
    ADD COLUMN sy_variant_id bigint NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX product_sy_id
    ON public.product USING btree
    (sy_id ASC NULLS LAST, sy_variant_id ASC NULLS LAST)

    WHERE sy_id != 0;

ALTER TABLE public.sales_order
    ADD COLUMN sy_id bigint NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX sales_order_sy_id
    ON public.sales_order USING btree
    (sy_id ASC NULLS LAST)

    WHERE sy_id != 0;

ALTER TABLE public.sales_order_detail
    ADD COLUMN sy_id bigint NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX sales_order_detail_sy_id
    ON public.sales_order_detail USING btree
    (sy_id ASC NULLS LAST)

    WHERE sy_id != 0;

ALTER TABLE public.sales_order
    ADD COLUMN sy_draft_id bigint NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX sales_order_sy_draft_id
    ON public.sales_order USING btree
    (sy_draft_id ASC NULLS LAST)

    WHERE sy_draft_id != 0;

ALTER TABLE public.sales_order_detail
    ADD COLUMN sy_draft_id bigint NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX sales_order_detail_sy_draft_id
    ON public.sales_order_detail USING btree
    (sy_draft_id ASC NULLS LAST)
    TABLESPACE pg_default
    WHERE sy_draft_id <> 0;

ALTER TABLE public.config
    ADD COLUMN shopify_export_serie character(3);

ALTER TABLE public.config
    ADD COLUMN shopify_intracommunity_serie character(3);

ALTER TABLE public.config
    ADD COLUMN shopify_interior_serie character(3);

ALTER TABLE public.config
    ADD COLUMN shopify_default_payment_method smallint;
ALTER TABLE public.config
    ADD CONSTRAINT config_shopify_export_serie FOREIGN KEY (shopify_export_serie)
    REFERENCES public.billing_series (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.config
    ADD CONSTRAINT config_shopify_interior_serie FOREIGN KEY (shopify_interior_serie)
    REFERENCES public.billing_series (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.config
    ADD CONSTRAINT config_shopify_intracommunity_serie FOREIGN KEY (shopify_intracommunity_serie)
    REFERENCES public.billing_series (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.config
    ADD CONSTRAINT config_shopify_default_payment_method FOREIGN KEY (shopify_default_payment_method)
    REFERENCES public.payment_method (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE public.config
    ADD COLUMN shopify_shop_location_id bigint NOT NULL DEFAULT 0;

ALTER TABLE public.payment_method
    ADD COLUMN shopify_module_name character varying(100) NOT NULL DEFAULT '';

CREATE TABLE public.sy_addresses
(
    id bigint NOT NULL,
    customer_id bigint NOT NULL,
    company character varying(100) COLLATE pg_catalog."default" NOT NULL,
    address1 character varying(100) COLLATE pg_catalog."default" NOT NULL,
    address2 character varying(100) COLLATE pg_catalog."default" NOT NULL,
    city character varying(50) COLLATE pg_catalog."default" NOT NULL,
    province character varying(50) COLLATE pg_catalog."default" NOT NULL,
    zip character varying(25) COLLATE pg_catalog."default" NOT NULL,
    country_code character varying(5) COLLATE pg_catalog."default" NOT NULL,
    sy_exists boolean NOT NULL DEFAULT true,
    CONSTRAINT sy_addresses_pkey PRIMARY KEY (id)
);

CREATE TABLE public.sy_customers
(
    id bigint NOT NULL,
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    first_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    last_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    tax_exempt boolean NOT NULL,
    phone character varying(25) COLLATE pg_catalog."default" NOT NULL,
    currency character varying(5) COLLATE pg_catalog."default" NOT NULL,
    sy_exists boolean NOT NULL DEFAULT true,
    default_address_id bigint NOT NULL,
    CONSTRAINT sy_customers_pkey PRIMARY KEY (id)
);

CREATE TABLE public.sy_draft_order_line_item
(
    id bigint NOT NULL,
    variant_id bigint NOT NULL,
    product_id bigint NOT NULL,
    quantity integer NOT NULL,
    taxable boolean NOT NULL,
    price real NOT NULL,
    draft_order_id bigint NOT NULL,
    sy_exists boolean NOT NULL DEFAULT true,
    CONSTRAINT sy_draft_order_line_item_pkey PRIMARY KEY (id)
);

CREATE TABLE public.sy_draft_orders
(
    id bigint NOT NULL,
    currency character varying(5) COLLATE pg_catalog."default" NOT NULL,
    tax_exempt boolean NOT NULL,
    name character varying(9) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_1 character varying(100) COLLATE pg_catalog."default" NOT NULL,
    shipping_address2 character varying(100) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_city character varying(50) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_zip character varying(25) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_country_code character varying(5) COLLATE pg_catalog."default" NOT NULL,
    billing_address_1 character varying(100) COLLATE pg_catalog."default" NOT NULL,
    billing_address2 character varying(100) COLLATE pg_catalog."default" NOT NULL,
    billing_address_city character varying(50) COLLATE pg_catalog."default" NOT NULL,
    billing_address_zip character varying(25) COLLATE pg_catalog."default" NOT NULL,
    billing_address_country_code character varying(5) COLLATE pg_catalog."default" NOT NULL,
    total_tax real NOT NULL,
    customer_id bigint NOT NULL,
    sy_exists boolean NOT NULL DEFAULT true,
    order_id bigint,
    CONSTRAINT sy_draft_orders_pkey PRIMARY KEY (id)
);

CREATE TABLE public.sy_order_line_item
(
    id bigint NOT NULL,
    variant_id bigint NOT NULL,
    product_id bigint NOT NULL,
    quantity integer NOT NULL,
    taxable boolean NOT NULL,
    price real NOT NULL,
    order_id bigint NOT NULL,
    sy_exists boolean NOT NULL DEFAULT true,
    CONSTRAINT sy_order_line_item_pkey PRIMARY KEY (id)
);

CREATE TABLE public.sy_orders
(
    id bigint NOT NULL,
    currency character varying(5) COLLATE pg_catalog."default" NOT NULL,
    current_total_discounts real NOT NULL,
    total_shipping_price_set_amount real NOT NULL,
    total_shipping_price_set_currency_code character varying(5) COLLATE pg_catalog."default" NOT NULL,
    tax_exempt boolean NOT NULL,
    name character varying(9) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_1 character varying(100) COLLATE pg_catalog."default" NOT NULL,
    shipping_address2 character varying(100) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_city character varying(50) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_zip character varying(25) COLLATE pg_catalog."default" NOT NULL,
    shipping_address_country_code character varying(5) COLLATE pg_catalog."default" NOT NULL,
    billing_address_1 character varying(100) COLLATE pg_catalog."default" NOT NULL,
    billing_address2 character varying(100) COLLATE pg_catalog."default" NOT NULL,
    billing_address_city character varying(50) COLLATE pg_catalog."default" NOT NULL,
    billing_address_zip character varying(25) COLLATE pg_catalog."default" NOT NULL,
    billing_address_country_code character varying(5) COLLATE pg_catalog."default" NOT NULL,
    total_tax real NOT NULL,
    customer_id bigint NOT NULL,
    sy_exists boolean NOT NULL DEFAULT true,
    gateway character varying(50) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT sy_orders_pkey PRIMARY KEY (id)
);

CREATE TABLE public.sy_products
(
    id bigint NOT NULL,
    title character varying(150) COLLATE pg_catalog."default" NOT NULL,
    body_html text COLLATE pg_catalog."default" NOT NULL,
    sy_exists boolean NOT NULL DEFAULT true,
    CONSTRAINT sy_products_pkey PRIMARY KEY (id)
);

CREATE TABLE public.sy_variants
(
    id bigint NOT NULL,
    product_id bigint NOT NULL,
    title character varying(150) COLLATE pg_catalog."default" NOT NULL,
    price real NOT NULL,
    sku character varying(25) COLLATE pg_catalog."default" NOT NULL,
    option1 character varying(150) COLLATE pg_catalog."default" NOT NULL,
    option2 character varying(150) COLLATE pg_catalog."default",
    option3 character varying(150) COLLATE pg_catalog."default",
    taxable boolean NOT NULL,
    barcode character varying(25) COLLATE pg_catalog."default" NOT NULL,
    grams integer NOT NULL,
    sy_exists boolean NOT NULL DEFAULT true,
    CONSTRAINT sy_variants_pkey PRIMARY KEY (id)
);

CREATE INDEX sy_draft_orders_order_id
    ON public.sy_draft_orders USING btree
    (order_id ASC NULLS LAST)
    TABLESPACE pg_default
    WHERE order_id IS NOT NULL;