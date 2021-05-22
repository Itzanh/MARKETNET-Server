--
-- PostgreSQL database dump
--

-- Dumped from database version 13.1
-- Dumped by pg_dump version 13.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: set_address_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_address_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(address.id) END AS id FROM address) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_city_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_city_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(city.id) END AS id FROM city) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_color_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_color_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(color.id) END AS id FROM color) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_country_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_country_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(country.id) END AS id FROM country) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_currency_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_currency_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(currency.id) END AS id FROM currency) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_customer_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_customer_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(customer.id) END AS id FROM customer) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_language_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_language_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(language.id) END AS id FROM language) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_manufacturing_order_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_manufacturing_order_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(manufacturing_order.id) END AS id FROM manufacturing_order) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_manufacturing_order_type_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_manufacturing_order_type_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(manufacturing_order_type.id) END AS id FROM manufacturing_order_type) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_payment_method_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_payment_method_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(payment_method.id) END AS id FROM payment_method) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_product_family_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_product_family_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(product_family.id) END AS id FROM product_family) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_product_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_product_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(product.id) END AS id FROM product) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_sales_invoice_detail_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_sales_invoice_detail_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_invoice_detail.id) END AS id FROM sales_invoice_detail) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_sales_invoice_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_sales_invoice_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_invoice.id) END AS id FROM sales_invoice) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_sales_order_detail_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_sales_order_detail_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_order_detail.id) END AS id FROM sales_order_detail) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_sales_order_discount_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_sales_order_discount_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_order_discount.id) END AS id FROM sales_order_discount) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_sales_order_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_sales_order_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_order.id) END AS id FROM sales_order) + 1;
    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: address; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.address (
    id integer NOT NULL,
    customer integer NOT NULL,
    address character varying(200) NOT NULL,
    address_2 character varying(200) NOT NULL,
    city integer NOT NULL,
    province character varying(100) NOT NULL,
    country smallint NOT NULL,
    private_business character(1) NOT NULL,
    notes text NOT NULL
);


--
-- Name: api_key; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.api_key (
    id smallint NOT NULL,
    name character(64) NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    user_created smallint NOT NULL,
    off boolean DEFAULT false NOT NULL,
    "user" smallint NOT NULL
);


--
-- Name: billing_series; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.billing_series (
    id character(3) NOT NULL,
    name character varying(50) NOT NULL,
    billing_type character(1) NOT NULL,
    year smallint NOT NULL
);


--
-- Name: carrier; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.carrier (
    id smallint NOT NULL,
    name character varying(50) NOT NULL,
    max_weight real NOT NULL,
    max_width real NOT NULL,
    max_height real NOT NULL,
    max_depth real NOT NULL,
    max_packages smallint NOT NULL,
    phone character varying(15) NOT NULL,
    email character varying(100) NOT NULL,
    web character varying(200) NOT NULL,
    off boolean NOT NULL
);


--
-- Name: city; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.city (
    id integer NOT NULL,
    country smallint NOT NULL,
    name character varying(100) NOT NULL,
    zip_code character varying(15) NOT NULL
);


--
-- Name: color; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.color (
    id smallint NOT NULL,
    name character varying(50) NOT NULL,
    hex_color character(6) NOT NULL
);


--
-- Name: config; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.config (
    id smallint NOT NULL,
    default_vat_percent real NOT NULL,
    default_warehouse character(2) NOT NULL,
    date_format character varying(25) NOT NULL
);


--
-- Name: country; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.country (
    id smallint NOT NULL,
    name character varying(50) NOT NULL,
    iso_2 character(2) NOT NULL,
    iso_3 character(3) NOT NULL,
    un_code smallint NOT NULL,
    zone character(1) NOT NULL,
    phone_prefix smallint NOT NULL,
    language smallint NOT NULL,
    currency smallint NOT NULL
);


--
-- Name: currency; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.currency (
    id smallint NOT NULL,
    name character varying(50) NOT NULL,
    sign character(3) NOT NULL,
    iso_code character(3) NOT NULL,
    iso_num smallint NOT NULL,
    change real NOT NULL
);


--
-- Name: customer; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.customer (
    id integer NOT NULL,
    name character varying(303) NOT NULL,
    tradename character varying(150) NOT NULL,
    fiscal_name character varying(150) NOT NULL,
    tax_id character varying(25) NOT NULL,
    vat_number character varying(25) NOT NULL,
    phone character varying(25) NOT NULL,
    email character varying(100) NOT NULL,
    main_address integer,
    country smallint,
    city integer,
    main_shipping_address integer,
    main_billing_address integer,
    language smallint,
    payment_method integer,
    billing_series character(3),
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL
);


--
-- Name: document; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.document (
    id integer NOT NULL,
    name character varying(250) NOT NULL,
    uuid uuid NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_updated timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    size integer NOT NULL,
    container smallint NOT NULL,
    dsc text NOT NULL,
    sales_order integer,
    sales_invoice integer,
    sales_delivery_note integer,
    shipping integer,
    purchase_order integer,
    purchase_invoice integer,
    purchase_delivery_note integer
);


--
-- Name: document_container; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.document_container (
    id smallint NOT NULL,
    name character varying(50) NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    path character varying(250) NOT NULL,
    max_file_size integer NOT NULL
);


--
-- Name: group; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."group" (
    id smallint NOT NULL,
    name character varying(50) NOT NULL,
    sales boolean NOT NULL,
    purchases boolean NOT NULL,
    masters boolean NOT NULL,
    warehouse boolean NOT NULL,
    manufacturing boolean NOT NULL,
    preparation boolean NOT NULL,
    admin boolean NOT NULL
);


--
-- Name: incoterm; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.incoterm (
    id smallint NOT NULL,
    key character(3) NOT NULL,
    name character varying(50) NOT NULL
);


--
-- Name: language; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language (
    id smallint NOT NULL,
    name character varying(50) NOT NULL,
    iso_2 character(2) NOT NULL,
    iso_3 character(3) NOT NULL
);


--
-- Name: login_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.login_tokens (
    id smallint NOT NULL,
    name character(128) NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    "user" smallint NOT NULL,
    ip_address character varying(15) NOT NULL
);


--
-- Name: manufacturing_order; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.manufacturing_order (
    id bigint NOT NULL,
    order_detail integer,
    product integer NOT NULL,
    type smallint NOT NULL,
    uuid uuid NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_last_update timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    manufactured boolean DEFAULT false NOT NULL,
    date_manufactured timestamp(3) without time zone,
    user_manufactured smallint,
    user_created smallint NOT NULL,
    tag_printed boolean DEFAULT false NOT NULL,
    date_tag_printed timestamp(3) without time zone,
    "order" integer,
    user_tag_printed smallint
);


--
-- Name: manufacturing_order_type; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.manufacturing_order_type (
    id smallint NOT NULL,
    name character varying(100) NOT NULL
);


--
-- Name: packages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.packages (
    id smallint NOT NULL,
    name character varying(50) NOT NULL,
    weight real NOT NULL,
    width real NOT NULL,
    height real NOT NULL,
    depth real NOT NULL
);


--
-- Name: packaging; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.packaging (
    id integer NOT NULL,
    packaging smallint NOT NULL,
    sales_order integer NOT NULL,
    weight real NOT NULL,
    shipping integer
);


--
-- Name: payment_method; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payment_method (
    id smallint NOT NULL,
    name character varying(100) NOT NULL,
    paid_in_advance boolean NOT NULL
);


--
-- Name: product; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product (
    id integer NOT NULL,
    name character varying(150) NOT NULL,
    reference character varying(40) NOT NULL,
    barcode character(13) NOT NULL,
    control_stock boolean NOT NULL,
    weight real NOT NULL,
    family smallint,
    width real NOT NULL,
    height real NOT NULL,
    depth real NOT NULL,
    off boolean NOT NULL,
    stock integer NOT NULL,
    vat_percent real NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    dsc text NOT NULL,
    color smallint,
    price real NOT NULL,
    manufacturing boolean NOT NULL,
    manufacturing_order_type smallint
);


--
-- Name: product_family; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_family (
    id smallint NOT NULL,
    name character varying(100) NOT NULL,
    reference character varying(40) NOT NULL
);


--
-- Name: product_image; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_image (
    product integer NOT NULL,
    url character varying(255) NOT NULL
);


--
-- Name: product_pack; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_pack (
    product_base integer NOT NULL,
    product_included integer NOT NULL,
    quantity smallint NOT NULL
);


--
-- Name: product_translation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_translation (
    product integer NOT NULL,
    language smallint NOT NULL,
    name character varying(150) NOT NULL
);


--
-- Name: purchase_delivery_note; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_delivery_note (
    id integer NOT NULL,
    warehouse character(2) NOT NULL,
    supplier integer NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    payment_method smallint NOT NULL,
    billing_series character(3) NOT NULL,
    shipping_address integer NOT NULL,
    total_products real NOT NULL,
    discount_percent real NOT NULL,
    fix_discount real NOT NULL,
    shipping_price real NOT NULL,
    shipping_discount real NOT NULL,
    total_with_discount real NOT NULL,
    total_vat real NOT NULL,
    total_amount real NOT NULL,
    lines_number smallint NOT NULL
);


--
-- Name: purchase_invoice; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_invoice (
    id integer NOT NULL,
    warehouse character(2) NOT NULL,
    supplier integer NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    payment_method smallint NOT NULL,
    billing_series character(3) NOT NULL,
    currency smallint NOT NULL,
    currency_change real NOT NULL,
    billing_address integer NOT NULL,
    total_products real NOT NULL,
    discount_percent real NOT NULL,
    fix_discount real NOT NULL,
    shipping_price real NOT NULL,
    shipping_discount real NOT NULL,
    total_with_discount real NOT NULL,
    total_vat real NOT NULL,
    total_amount real NOT NULL,
    lines_number smallint NOT NULL
);


--
-- Name: purchase_invoice_details; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_invoice_details (
    id integer NOT NULL,
    invoice integer NOT NULL,
    product integer NOT NULL,
    price real NOT NULL,
    quantity integer NOT NULL,
    vat_percent real NOT NULL,
    total_amount real NOT NULL
);


--
-- Name: purchase_order; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_order (
    id integer NOT NULL,
    warehouse character(2) NOT NULL,
    supplier_reference character varying(40) NOT NULL,
    supplier integer NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_paid timestamp(3) without time zone,
    payment_method smallint NOT NULL,
    billing_series character(3) NOT NULL,
    currency smallint NOT NULL,
    currency_change real NOT NULL,
    billing_address integer NOT NULL,
    shipping_address integer NOT NULL,
    lines_number smallint NOT NULL,
    invoiced_lines smallint NOT NULL,
    delivery_note_lines smallint NOT NULL,
    total_products real NOT NULL,
    discount_percent real NOT NULL,
    fix_discount real NOT NULL,
    shipping_price real NOT NULL,
    shipping_discount real NOT NULL,
    total_with_discount real NOT NULL,
    total_vat real NOT NULL,
    total_amount real NOT NULL,
    dsc text NOT NULL,
    notes character varying(250) NOT NULL,
    off boolean NOT NULL,
    cancelled boolean NOT NULL,
    status character(1) NOT NULL,
    order_number integer NOT NULL
);


--
-- Name: purchase_order_detail; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_order_detail (
    id integer NOT NULL,
    "order" integer NOT NULL,
    product integer NOT NULL,
    price real NOT NULL,
    quantity integer NOT NULL,
    vat_percent real NOT NULL,
    total_amount real NOT NULL,
    quantity_invoiced integer NOT NULL,
    quantity_delivery_note integer NOT NULL,
    status character(1) NOT NULL,
    quantity_pending_packaging integer NOT NULL
);


--
-- Name: sales_delivery_note; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_delivery_note (
    id integer NOT NULL,
    warehouse character(2) NOT NULL,
    customer integer NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    payment_method smallint,
    billing_series character(3) NOT NULL,
    shipping_address integer NOT NULL,
    total_products real NOT NULL,
    discount_percent real NOT NULL,
    fix_discount real NOT NULL,
    shipping_price real NOT NULL,
    shipping_discount real NOT NULL,
    total_with_discount real NOT NULL,
    total_vat real NOT NULL,
    total_amount real NOT NULL,
    lines_number real NOT NULL
);


--
-- Name: sales_invoice; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_invoice (
    id integer NOT NULL,
    customer integer NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    payment_method smallint NOT NULL,
    billing_series character(3) NOT NULL,
    currency smallint NOT NULL,
    currency_change real NOT NULL,
    billing_address integer NOT NULL,
    total_products real DEFAULT 0 NOT NULL,
    discount_percent real DEFAULT 0 NOT NULL,
    fix_discount real DEFAULT 0 NOT NULL,
    shipping_price real DEFAULT 0 NOT NULL,
    shipping_discount real DEFAULT 0 NOT NULL,
    total_with_discount real DEFAULT 0 NOT NULL,
    vat_amount real DEFAULT 0 NOT NULL,
    total_amount real DEFAULT 0 NOT NULL,
    lines_number smallint DEFAULT 0 NOT NULL,
    invoice_number integer NOT NULL,
    invoice_name character(15) NOT NULL
);


--
-- Name: sales_invoice_detail; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_invoice_detail (
    id integer NOT NULL,
    invoice integer NOT NULL,
    product integer NOT NULL,
    price real NOT NULL,
    quantity integer NOT NULL,
    vat_percent real NOT NULL,
    total_amount real NOT NULL,
    order_detail integer
);


--
-- Name: sales_order; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_order (
    id integer NOT NULL,
    warehouse character(2) NOT NULL,
    reference character varying(9) NOT NULL,
    customer integer NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_payment_accepted timestamp(3) without time zone,
    payment_method smallint NOT NULL,
    billing_series character(3) NOT NULL,
    currency smallint NOT NULL,
    currency_change real NOT NULL,
    billing_address integer NOT NULL,
    shipping_address integer NOT NULL,
    lines_number smallint DEFAULT 0 NOT NULL,
    invoiced_lines smallint DEFAULT 0 NOT NULL,
    delivery_note_lines smallint DEFAULT 0 NOT NULL,
    total_products real DEFAULT 0 NOT NULL,
    discount_percent real NOT NULL,
    fix_discount real NOT NULL,
    shipping_price real NOT NULL,
    shipping_discount real NOT NULL,
    total_with_discount real DEFAULT 0 NOT NULL,
    vat_amount real DEFAULT 0 NOT NULL,
    total_amount real DEFAULT 0 NOT NULL,
    dsc text NOT NULL,
    notes character varying(250) NOT NULL,
    off boolean DEFAULT false NOT NULL,
    cancelled boolean DEFAULT false NOT NULL,
    status character(1) DEFAULT '_'::bpchar NOT NULL,
    order_number integer NOT NULL,
    billing_status character(1) DEFAULT 'P'::bpchar NOT NULL,
    order_name character(15) NOT NULL
);


--
-- Name: sales_order_detail; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_order_detail (
    id integer NOT NULL,
    "order" integer NOT NULL,
    product integer NOT NULL,
    price real NOT NULL,
    quantity integer NOT NULL,
    vat_percent real NOT NULL,
    total_amount real NOT NULL,
    quantity_invoiced integer DEFAULT 0 NOT NULL,
    quantity_delivery_note integer DEFAULT 0 NOT NULL,
    status character(1) NOT NULL,
    quantity_pending_packaging integer NOT NULL
);


--
-- Name: sales_order_detail_packaged; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_order_detail_packaged (
    order_detail integer NOT NULL,
    packaging integer NOT NULL,
    quantity integer NOT NULL
);


--
-- Name: sales_order_discount; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_order_discount (
    id integer NOT NULL,
    "order" integer NOT NULL,
    name character varying(100) NOT NULL,
    value_tax_included real NOT NULL,
    value_tax_excluded real NOT NULL
);


--
-- Name: shipping; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.shipping (
    id integer NOT NULL,
    "order" integer NOT NULL,
    delivery_note integer NOT NULL,
    delivery_address integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_sent timestamp(3) without time zone,
    sent boolean NOT NULL,
    collected boolean NOT NULL,
    "national" boolean NOT NULL,
    shipping_number character varying(50) NOT NULL,
    tracking_number character varying(50) NOT NULL,
    carrier smallint NOT NULL,
    weight real NOT NULL,
    packages_number smallint NOT NULL
);


--
-- Name: stock; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.stock (
    product integer NOT NULL,
    warehouse character(2) NOT NULL,
    quantity integer DEFAULT 0 NOT NULL,
    quantity_pending_received integer DEFAULT 0 NOT NULL,
    quantity_pending_served integer DEFAULT 0 NOT NULL,
    quantity_available integer DEFAULT 0 NOT NULL,
    quantity_pending_manufacture integer DEFAULT 0 NOT NULL
);


--
-- Name: suppliers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.suppliers (
    id integer NOT NULL,
    name character varying(300) NOT NULL,
    tradename character varying(150) NOT NULL,
    fiscal_name character varying(150) NOT NULL,
    tax_id character varying(25) NOT NULL,
    vat_number character varying(25) NOT NULL,
    phone character varying(15) NOT NULL,
    email character varying(100) NOT NULL,
    main_address integer,
    country smallint,
    city integer,
    main_shipping_address integer,
    main_billing_address integer,
    language smallint,
    payment_method integer,
    billing_series character(3)
);


--
-- Name: user; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."user" (
    id smallint NOT NULL,
    username character varying(40) NOT NULL,
    full_name character varying(150) NOT NULL,
    email character varying(100) NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_last_pwd timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    pwd_next_login boolean NOT NULL,
    off boolean NOT NULL,
    pwd bytea,
    salt character(30) NOT NULL,
    iterations integer NOT NULL,
    dsc text NOT NULL,
    date_last_login timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL
);


--
-- Name: user_group; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_group (
    "user" smallint NOT NULL,
    "group" smallint NOT NULL
);


--
-- Name: warehouse; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehouse (
    id character(2) NOT NULL,
    name character varying(50) NOT NULL
);


--
-- Name: warehouse_movement; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehouse_movement (
    id bigint NOT NULL,
    warehouse character(2) NOT NULL,
    product integer NOT NULL,
    quantity integer NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    type character(1) NOT NULL,
    current_stock integer NOT NULL,
    sales_order integer,
    sales_order_detail integer,
    sales_invoice integer,
    sales_invoice_detail integer,
    sales_delivery_note integer,
    dsc text NOT NULL,
    purchase_order integer,
    purchase_order_detail integer,
    purchase_invoice integer,
    purchase_invoice_details integer,
    purchase_delivery_note integer
);


--
-- Name: address address_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_pkey PRIMARY KEY (id);


--
-- Name: api_key api_key_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT api_key_pkey PRIMARY KEY (id);


--
-- Name: billing_series billing_series_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.billing_series
    ADD CONSTRAINT billing_series_pkey PRIMARY KEY (id);


--
-- Name: carrier carrier_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.carrier
    ADD CONSTRAINT carrier_pkey PRIMARY KEY (id);


--
-- Name: city city_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT city_pkey PRIMARY KEY (id);


--
-- Name: color color_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.color
    ADD CONSTRAINT color_pkey PRIMARY KEY (id);


--
-- Name: config config_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_pkey PRIMARY KEY (default_vat_percent);


--
-- Name: country country_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_pkey PRIMARY KEY (id);


--
-- Name: currency currency_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.currency
    ADD CONSTRAINT currency_pkey PRIMARY KEY (id);


--
-- Name: customer customer_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_pkey PRIMARY KEY (id);


--
-- Name: document_container document_container_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document_container
    ADD CONSTRAINT document_container_pkey PRIMARY KEY (id);


--
-- Name: document document_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_pkey PRIMARY KEY (id);


--
-- Name: group group_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."group"
    ADD CONSTRAINT group_pkey PRIMARY KEY (id);


--
-- Name: incoterm incoterm_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.incoterm
    ADD CONSTRAINT incoterm_pkey PRIMARY KEY (id);


--
-- Name: language language_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language
    ADD CONSTRAINT language_pkey PRIMARY KEY (id);


--
-- Name: login_tokens login_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.login_tokens
    ADD CONSTRAINT login_tokens_pkey PRIMARY KEY (id);


--
-- Name: manufacturing_order manufacturing_order_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_pkey PRIMARY KEY (id);


--
-- Name: manufacturing_order_type manufacturing_order_type_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order_type
    ADD CONSTRAINT manufacturing_order_type_pkey PRIMARY KEY (id);


--
-- Name: packages packages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packages
    ADD CONSTRAINT packages_pkey PRIMARY KEY (id);


--
-- Name: packaging packaing_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packaging
    ADD CONSTRAINT packaing_pkey PRIMARY KEY (id);


--
-- Name: payment_method payment_method_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_method
    ADD CONSTRAINT payment_method_pkey PRIMARY KEY (id);


--
-- Name: product_family product_family_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_family
    ADD CONSTRAINT product_family_pkey PRIMARY KEY (id);


--
-- Name: product_image product_image_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_image
    ADD CONSTRAINT product_image_pkey PRIMARY KEY (product);


--
-- Name: product_pack product_pack_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_pack
    ADD CONSTRAINT product_pack_pkey PRIMARY KEY (product_base, product_included);


--
-- Name: product product_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_pkey PRIMARY KEY (id);


--
-- Name: product_translation product_translation_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_translation
    ADD CONSTRAINT product_translation_pkey PRIMARY KEY (product, language);


--
-- Name: purchase_delivery_note purchase_delivery_note_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_delivery_note
    ADD CONSTRAINT purchase_delivery_note_pkey PRIMARY KEY (id);


--
-- Name: purchase_invoice_details purchase_invoice_details_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice_details
    ADD CONSTRAINT purchase_invoice_details_pkey PRIMARY KEY (id);


--
-- Name: purchase_invoice purchase_invoice_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice
    ADD CONSTRAINT purchase_invoice_pkey PRIMARY KEY (id);


--
-- Name: purchase_order purchase_order_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_pkey PRIMARY KEY (id);


--
-- Name: purchase_order_detail purchase_orer_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order_detail
    ADD CONSTRAINT purchase_orer_detail_pkey PRIMARY KEY (id);


--
-- Name: sales_delivery_note sales_delivery_note_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_pkey PRIMARY KEY (id);


--
-- Name: sales_invoice_detail sales_invoice_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice_detail
    ADD CONSTRAINT sales_invoice_detail_pkey PRIMARY KEY (id);


--
-- Name: sales_invoice sales_invoice_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_pkey PRIMARY KEY (id);


--
-- Name: sales_order_detail_packaged sales_order_detail_packaged_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail_packaged
    ADD CONSTRAINT sales_order_detail_packaged_pkey PRIMARY KEY (order_detail, packaging);


--
-- Name: sales_order_detail sales_order_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail
    ADD CONSTRAINT sales_order_detail_pkey PRIMARY KEY (id);


--
-- Name: sales_order_discount sales_order_discount_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_discount
    ADD CONSTRAINT sales_order_discount_pkey PRIMARY KEY (id);


--
-- Name: sales_order sales_order_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_pkey PRIMARY KEY (id);


--
-- Name: shipping shipping_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_pkey PRIMARY KEY (id);


--
-- Name: stock stock_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_pkey PRIMARY KEY (product, warehouse);


--
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY (id);


--
-- Name: user_group user_group_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_group
    ADD CONSTRAINT user_group_pkey PRIMARY KEY ("user", "group");


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: warehouse_movement warehouse_movement_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_pkey PRIMARY KEY (id);


--
-- Name: warehouse warehouse_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse
    ADD CONSTRAINT warehouse_pkey PRIMARY KEY (id);


--
-- Name: city_zip_code; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX city_zip_code ON public.city USING btree (country, zip_code);


--
-- Name: country_iso_2; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX country_iso_2 ON public.country USING btree (iso_2);


--
-- Name: country_iso_3; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX country_iso_3 ON public.country USING btree (iso_3);


--
-- Name: currency_iso_code; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX currency_iso_code ON public.currency USING btree (iso_code);


--
-- Name: currency_iso_num; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX currency_iso_num ON public.currency USING btree (iso_num);


--
-- Name: currency_sign; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX currency_sign ON public.currency USING btree (sign);


--
-- Name: incoterm_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX incoterm_key ON public.incoterm USING btree (key);


--
-- Name: language_iso_2; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX language_iso_2 ON public.language USING btree (iso_2);


--
-- Name: language_iso_3; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX language_iso_3 ON public.language USING btree (iso_3);


--
-- Name: manufacturing_order_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX manufacturing_order_date_created ON public.manufacturing_order USING btree (date_created DESC NULLS LAST);


--
-- Name: manufacturing_order_uuid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX manufacturing_order_uuid ON public.manufacturing_order USING btree (uuid);


--
-- Name: product_barcode; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX product_barcode ON public.product USING btree (barcode);


--
-- Name: product_family_reference; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX product_family_reference ON public.product_family USING btree (reference);


--
-- Name: product_reference; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX product_reference ON public.product USING btree (reference);


--
-- Name: purchase_order_order_number; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX purchase_order_order_number ON public.purchase_order USING btree (billing_series, order_number DESC NULLS LAST);


--
-- Name: sales_invoice_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX sales_invoice_date_created ON public.sales_invoice USING btree (date_created DESC NULLS LAST);


--
-- Name: sales_invoice_invoice_number; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_invoice_invoice_number ON public.sales_invoice USING btree (billing_series, invoice_number DESC NULLS LAST);


--
-- Name: sales_order_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX sales_order_date_created ON public.sales_order USING btree (date_created DESC NULLS LAST);


--
-- Name: sales_order_order_number; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_order_number ON public.sales_order USING btree (billing_series, order_number DESC NULLS LAST);


--
-- Name: user_username; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX user_username ON public."user" USING btree (username);


--
-- Name: address set_address_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_address_id BEFORE INSERT ON public.address FOR EACH ROW EXECUTE FUNCTION public.set_address_id();


--
-- Name: city set_city_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_city_id BEFORE INSERT ON public.city FOR EACH ROW EXECUTE FUNCTION public.set_city_id();


--
-- Name: color set_color_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_color_id BEFORE INSERT ON public.color FOR EACH ROW EXECUTE FUNCTION public.set_color_id();


--
-- Name: country set_country_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_country_id BEFORE INSERT ON public.country FOR EACH ROW EXECUTE FUNCTION public.set_country_id();


--
-- Name: currency set_currency_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_currency_id BEFORE INSERT ON public.currency FOR EACH ROW EXECUTE FUNCTION public.set_currency_id();


--
-- Name: customer set_customer_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_customer_id BEFORE INSERT ON public.customer FOR EACH ROW EXECUTE FUNCTION public.set_customer_id();


--
-- Name: language set_language_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_language_id BEFORE INSERT ON public.language FOR EACH ROW EXECUTE FUNCTION public.set_language_id();


--
-- Name: manufacturing_order set_manufacturing_order_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_manufacturing_order_id BEFORE INSERT ON public.manufacturing_order FOR EACH ROW EXECUTE FUNCTION public.set_manufacturing_order_id();


--
-- Name: manufacturing_order_type set_manufacturing_order_type_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_manufacturing_order_type_id BEFORE INSERT ON public.manufacturing_order_type FOR EACH ROW EXECUTE FUNCTION public.set_manufacturing_order_type_id();


--
-- Name: payment_method set_payment_method_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_payment_method_id BEFORE INSERT ON public.payment_method FOR EACH ROW EXECUTE FUNCTION public.set_payment_method_id();


--
-- Name: product_family set_product_family_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_product_family_id BEFORE INSERT ON public.product_family FOR EACH ROW EXECUTE FUNCTION public.set_product_family_id();


--
-- Name: product set_product_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_product_id BEFORE INSERT ON public.product FOR EACH ROW EXECUTE FUNCTION public.set_product_id();


--
-- Name: sales_invoice_detail set_sales_invoice_detail_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_sales_invoice_detail_id BEFORE INSERT ON public.sales_invoice_detail FOR EACH ROW EXECUTE FUNCTION public.set_sales_invoice_detail_id();


--
-- Name: sales_invoice set_sales_invoice_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_sales_invoice_id BEFORE INSERT ON public.sales_invoice FOR EACH ROW EXECUTE FUNCTION public.set_sales_invoice_id();


--
-- Name: sales_order_detail set_sales_order_detail_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_sales_order_detail_id BEFORE INSERT ON public.sales_order_detail FOR EACH ROW EXECUTE FUNCTION public.set_sales_order_detail_id();


--
-- Name: sales_order_discount set_sales_order_discount_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_sales_order_discount_id BEFORE INSERT ON public.sales_order_discount FOR EACH ROW EXECUTE FUNCTION public.set_sales_order_discount_id();


--
-- Name: sales_order set_sales_order_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_sales_order_id BEFORE INSERT ON public.sales_order FOR EACH ROW EXECUTE FUNCTION public.set_sales_order_id();


--
-- Name: address address_city; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_city FOREIGN KEY (city) REFERENCES public.city(id) NOT VALID;


--
-- Name: address address_country; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_country FOREIGN KEY (country) REFERENCES public.country(id) NOT VALID;


--
-- Name: address address_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_customer FOREIGN KEY (customer) REFERENCES public.customer(id) NOT VALID;


--
-- Name: api_key api_key_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT api_key_user FOREIGN KEY ("user") REFERENCES public."user"(id);


--
-- Name: api_key api_key_user_created; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT api_key_user_created FOREIGN KEY (user_created) REFERENCES public."user"(id);


--
-- Name: city city_country; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT city_country FOREIGN KEY (country) REFERENCES public.country(id) NOT VALID;


--
-- Name: country country_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_currency FOREIGN KEY (currency) REFERENCES public.currency(id) NOT VALID;


--
-- Name: country country_language; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_language FOREIGN KEY (language) REFERENCES public.language(id) NOT VALID;


--
-- Name: customer customer_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_billing_series FOREIGN KEY (billing_series) REFERENCES public.billing_series(id) NOT VALID;


--
-- Name: customer customer_city; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_city FOREIGN KEY (city) REFERENCES public.city(id) NOT VALID;


--
-- Name: customer customer_country; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_country FOREIGN KEY (country) REFERENCES public.country(id) NOT VALID;


--
-- Name: customer customer_language; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_language FOREIGN KEY (language) REFERENCES public.language(id) NOT VALID;


--
-- Name: customer customer_main_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_main_address FOREIGN KEY (main_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: customer customer_main_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_main_billing_address FOREIGN KEY (main_billing_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: customer customer_main_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_main_shipping_address FOREIGN KEY (main_shipping_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: customer customer_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_payment_method FOREIGN KEY (payment_method) REFERENCES public.payment_method(id) NOT VALID;


--
-- Name: manufacturing_order manufacturing_order_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_order FOREIGN KEY ("order") REFERENCES public.sales_order(id) NOT VALID;


--
-- Name: manufacturing_order manufacturing_order_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_order_detail FOREIGN KEY (order_detail) REFERENCES public.sales_order_detail(id) NOT VALID;


--
-- Name: manufacturing_order manufacturing_order_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_product FOREIGN KEY (product) REFERENCES public.product(id) NOT VALID;


--
-- Name: manufacturing_order manufacturing_order_type; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_type FOREIGN KEY (type) REFERENCES public.manufacturing_order_type(id) NOT VALID;


--
-- Name: manufacturing_order manufacturing_order_user_created; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_user_created FOREIGN KEY (user_created) REFERENCES public."user"(id) NOT VALID;


--
-- Name: manufacturing_order manufacturing_order_user_manufactured; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_user_manufactured FOREIGN KEY (user_manufactured) REFERENCES public."user"(id) NOT VALID;


--
-- Name: packaging packaging_packaging; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packaging
    ADD CONSTRAINT packaging_packaging FOREIGN KEY (packaging) REFERENCES public.packages(id) NOT VALID;


--
-- Name: packaging packaging_sales_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packaging
    ADD CONSTRAINT packaging_sales_order FOREIGN KEY (sales_order) REFERENCES public.sales_order(id) NOT VALID;


--
-- Name: packaging packaging_shipping; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packaging
    ADD CONSTRAINT packaging_shipping FOREIGN KEY (shipping) REFERENCES public.shipping(id) NOT VALID;


--
-- Name: product product_color; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_color FOREIGN KEY (color) REFERENCES public.color(id) NOT VALID;


--
-- Name: product_image product_image_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_image
    ADD CONSTRAINT product_image_product FOREIGN KEY (product) REFERENCES public.product(id);


--
-- Name: product_pack product_pack_product_base; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_pack
    ADD CONSTRAINT product_pack_product_base FOREIGN KEY (product_base) REFERENCES public.product(id) NOT VALID;


--
-- Name: product_pack product_pack_product_included; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_pack
    ADD CONSTRAINT product_pack_product_included FOREIGN KEY (product_included) REFERENCES public.product(id) NOT VALID;


--
-- Name: product product_product_family; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_product_family FOREIGN KEY (family) REFERENCES public.product_family(id) NOT VALID;


--
-- Name: product_translation product_translation_language; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_translation
    ADD CONSTRAINT product_translation_language FOREIGN KEY (language) REFERENCES public.language(id);


--
-- Name: product_translation product_translation_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_translation
    ADD CONSTRAINT product_translation_product FOREIGN KEY (product) REFERENCES public.product(id);


--
-- Name: purchase_invoice_details purchase_invoice_details_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice_details
    ADD CONSTRAINT purchase_invoice_details_invoice FOREIGN KEY (invoice) REFERENCES public.purchase_invoice(id) NOT VALID;


--
-- Name: purchase_invoice_details purchase_invoice_details_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice_details
    ADD CONSTRAINT purchase_invoice_details_product FOREIGN KEY (product) REFERENCES public.product(id) NOT VALID;


--
-- Name: purchase_order purchase_order_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_billing_address FOREIGN KEY (billing_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: purchase_order purchase_order_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_billing_series FOREIGN KEY (billing_series) REFERENCES public.billing_series(id) NOT VALID;


--
-- Name: purchase_order purchase_order_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_currency FOREIGN KEY (currency) REFERENCES public.currency(id) NOT VALID;


--
-- Name: purchase_order_detail purchase_order_detail_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order_detail
    ADD CONSTRAINT purchase_order_detail_order FOREIGN KEY ("order") REFERENCES public.purchase_order(id) NOT VALID;


--
-- Name: purchase_order_detail purchase_order_detail_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order_detail
    ADD CONSTRAINT purchase_order_detail_product FOREIGN KEY (product) REFERENCES public.product(id) NOT VALID;


--
-- Name: purchase_order purchase_order_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_payment_method FOREIGN KEY (payment_method) REFERENCES public.payment_method(id) NOT VALID;


--
-- Name: purchase_order purchase_order_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_shipping_address FOREIGN KEY (shipping_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: purchase_order purchase_order_supplier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_supplier FOREIGN KEY (supplier) REFERENCES public.suppliers(id) NOT VALID;


--
-- Name: purchase_order purchase_order_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_warehouse FOREIGN KEY (warehouse) REFERENCES public.warehouse(id) NOT VALID;


--
-- Name: sales_delivery_note sales_delivery_note_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_billing_series FOREIGN KEY (billing_series) REFERENCES public.billing_series(id) NOT VALID;


--
-- Name: sales_delivery_note sales_delivery_note_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_customer FOREIGN KEY (customer) REFERENCES public.customer(id) NOT VALID;


--
-- Name: sales_delivery_note sales_delivery_note_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_payment_method FOREIGN KEY (payment_method) REFERENCES public.payment_method(id) NOT VALID;


--
-- Name: sales_delivery_note sales_delivery_note_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_shipping_address FOREIGN KEY (shipping_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: sales_delivery_note sales_delivery_note_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_warehouse FOREIGN KEY (warehouse) REFERENCES public.warehouse(id) NOT VALID;


--
-- Name: sales_invoice sales_invoice_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_billing_address FOREIGN KEY (billing_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: sales_invoice sales_invoice_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_billing_series FOREIGN KEY (billing_series) REFERENCES public.billing_series(id) NOT VALID;


--
-- Name: sales_invoice sales_invoice_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_currency FOREIGN KEY (currency) REFERENCES public.currency(id) NOT VALID;


--
-- Name: sales_invoice sales_invoice_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_customer FOREIGN KEY (customer) REFERENCES public.customer(id) NOT VALID;


--
-- Name: sales_invoice_detail sales_invoice_detail_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice_detail
    ADD CONSTRAINT sales_invoice_detail_invoice FOREIGN KEY (invoice) REFERENCES public.sales_invoice(id) NOT VALID;


--
-- Name: sales_invoice_detail sales_invoice_detail_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice_detail
    ADD CONSTRAINT sales_invoice_detail_order_detail FOREIGN KEY (order_detail) REFERENCES public.sales_order_detail(id) NOT VALID;


--
-- Name: sales_invoice_detail sales_invoice_detail_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice_detail
    ADD CONSTRAINT sales_invoice_detail_product FOREIGN KEY (product) REFERENCES public.product(id) NOT VALID;


--
-- Name: sales_invoice sales_invoice_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_payment_method FOREIGN KEY (payment_method) REFERENCES public.payment_method(id) NOT VALID;


--
-- Name: sales_order sales_order_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_billing_address FOREIGN KEY (billing_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: sales_order sales_order_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_billing_series FOREIGN KEY (billing_series) REFERENCES public.billing_series(id) NOT VALID;


--
-- Name: sales_order sales_order_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_currency FOREIGN KEY (currency) REFERENCES public.currency(id) NOT VALID;


--
-- Name: sales_order sales_order_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_customer FOREIGN KEY (customer) REFERENCES public.customer(id) NOT VALID;


--
-- Name: sales_order_detail sales_order_detail_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail
    ADD CONSTRAINT sales_order_detail_order FOREIGN KEY ("order") REFERENCES public.sales_order(id) NOT VALID;


--
-- Name: sales_order_detail_packaged sales_order_detail_packaged_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail_packaged
    ADD CONSTRAINT sales_order_detail_packaged_order_detail FOREIGN KEY (order_detail) REFERENCES public.sales_order(id) NOT VALID;


--
-- Name: sales_order_detail_packaged sales_order_detail_packaged_packaging; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail_packaged
    ADD CONSTRAINT sales_order_detail_packaged_packaging FOREIGN KEY (packaging) REFERENCES public.packaging(id) NOT VALID;


--
-- Name: sales_order_detail sales_order_detail_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail
    ADD CONSTRAINT sales_order_detail_product FOREIGN KEY (product) REFERENCES public.product(id) NOT VALID;


--
-- Name: sales_order_discount sales_order_discount_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_discount
    ADD CONSTRAINT sales_order_discount_order FOREIGN KEY ("order") REFERENCES public.sales_order(id) NOT VALID;


--
-- Name: sales_order sales_order_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_payment_method FOREIGN KEY (payment_method) REFERENCES public.payment_method(id) NOT VALID;


--
-- Name: sales_order sales_order_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_shipping_address FOREIGN KEY (shipping_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: sales_order sales_order_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_warehouse FOREIGN KEY (warehouse) REFERENCES public.warehouse(id) NOT VALID;


--
-- Name: shipping shipping_carrier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_carrier FOREIGN KEY (carrier) REFERENCES public.carrier(id) NOT VALID;


--
-- Name: shipping shipping_delivery_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_delivery_address FOREIGN KEY (delivery_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: shipping shipping_delivery_note; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_delivery_note FOREIGN KEY (delivery_note) REFERENCES public.sales_delivery_note(id) NOT VALID;


--
-- Name: shipping shipping_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_order FOREIGN KEY ("order") REFERENCES public.sales_order(id) NOT VALID;


--
-- Name: stock stock_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_product FOREIGN KEY (product) REFERENCES public.product(id) NOT VALID;


--
-- Name: stock stock_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_warehouse FOREIGN KEY (warehouse) REFERENCES public.warehouse(id) NOT VALID;


--
-- Name: suppliers suppliers_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_billing_series FOREIGN KEY (billing_series) REFERENCES public.billing_series(id) NOT VALID;


--
-- Name: suppliers suppliers_city; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_city FOREIGN KEY (city) REFERENCES public.city(id) NOT VALID;


--
-- Name: suppliers suppliers_country; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_country FOREIGN KEY (country) REFERENCES public.country(id) NOT VALID;


--
-- Name: suppliers suppliers_language; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_language FOREIGN KEY (language) REFERENCES public.language(id) NOT VALID;


--
-- Name: suppliers suppliers_main_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_main_address FOREIGN KEY (main_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: suppliers suppliers_main_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_main_billing_address FOREIGN KEY (main_billing_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: suppliers suppliers_main_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_main_shipping_address FOREIGN KEY (main_shipping_address) REFERENCES public.address(id) NOT VALID;


--
-- Name: suppliers suppliers_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_payment_method FOREIGN KEY (payment_method) REFERENCES public.payment_method(id) NOT VALID;


--
-- Name: user_group user_group_group; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_group
    ADD CONSTRAINT user_group_group FOREIGN KEY ("group") REFERENCES public."group"(id) NOT VALID;


--
-- Name: user_group user_group_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_group
    ADD CONSTRAINT user_group_user FOREIGN KEY ("user") REFERENCES public."user"(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_delivery_note; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_delivery_note FOREIGN KEY (sales_delivery_note) REFERENCES public.sales_delivery_note(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_invoice FOREIGN KEY (sales_invoice) REFERENCES public.sales_invoice(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_invoice_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_invoice_detail FOREIGN KEY (sales_invoice_detail) REFERENCES public.sales_invoice_detail(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_order FOREIGN KEY (sales_order) REFERENCES public.sales_order(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_order_detail FOREIGN KEY (sales_order_detail) REFERENCES public.sales_order_detail(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_product FOREIGN KEY (product) REFERENCES public.product(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_purchase_delivery_note; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_purchase_delivery_note FOREIGN KEY (purchase_delivery_note) REFERENCES public.purchase_delivery_note(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_purchase_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_purchase_invoice FOREIGN KEY (purchase_invoice) REFERENCES public.purchase_invoice(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_purchase_invoice_details; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_purchase_invoice_details FOREIGN KEY (purchase_invoice_details) REFERENCES public.purchase_invoice_details(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_purchase_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_purchase_order FOREIGN KEY (purchase_order) REFERENCES public.purchase_order(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_purchase_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_purchase_order_detail FOREIGN KEY (purchase_order_detail) REFERENCES public.purchase_order_detail(id) NOT VALID;


--
-- Name: warehouse_movement warehouse_movement_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_warehouse FOREIGN KEY (warehouse) REFERENCES public.warehouse(id) NOT VALID;


--
-- PostgreSQL database dump complete
--

