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
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


--
-- Name: set_account_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_account_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(account.id) END AS id FROM account) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_accounting_movement_detail_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_accounting_movement_detail_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(accounting_movement_detail.id) END AS id FROM accounting_movement_detail) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_accounting_movement_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_accounting_movement_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(accounting_movement.id) END AS id FROM accounting_movement) + 1;
    RETURN NEW;
END;
$$;


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
-- Name: set_api_key_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_api_key_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(api_key.id) END AS id FROM api_key) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_carrier_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_carrier_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(carrier.id) END AS id FROM carrier) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_charges_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_charges_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(charges.id) END AS id FROM charges) + 1;
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
-- Name: set_collection_operation_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_collection_operation_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(collection_operation.id) END AS id FROM collection_operation) + 1;
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
-- Name: set_complex_manufacturing_order_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_complex_manufacturing_order_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(complex_manufacturing_order.id) END AS id FROM complex_manufacturing_order) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_complex_manufacturing_order_manufacturing_order_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_complex_manufacturing_order_manufacturing_order_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(complex_manufacturing_order_manufacturing_order.id) END AS id FROM complex_manufacturing_order_manufacturing_order) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_config_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_config_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(config.id) END AS id FROM config) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_connection_filter_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_connection_filter_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(connection_filter.id) END AS id FROM connection_filter) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_connection_log_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_connection_log_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(connection_log.id) END AS id FROM connection_log) + 1;
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
-- Name: set_custom_fields_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_custom_fields_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$

BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(custom_fields.id) END AS id FROM custom_fields) + 1;
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
-- Name: set_document_container_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_document_container_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(document_container.id) END AS id FROM document_container) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_document_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_document_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(document.id) END AS id FROM document) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_email_log_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_email_log_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(email_log.id) END AS id FROM email_log) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_group_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_group_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX("group".id) END AS id FROM "group") + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_incoterm_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_incoterm_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(incoterm.id) END AS id FROM incoterm) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_inventory_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_inventory_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(inventory.id) END AS id FROM inventory) + 1;
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
-- Name: set_login_tokens_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_login_tokens_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(login_tokens.id) END AS id FROM login_tokens) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_logs_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_logs_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(logs.id) END AS id FROM logs) + 1;
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
-- Name: set_manufacturing_order_type_components_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_manufacturing_order_type_components_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(manufacturing_order_type_components.id) END AS id FROM manufacturing_order_type_components) + 1;
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
-- Name: set_packages_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_packages_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(packages.id) END AS id FROM packages) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_packaging_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_packaging_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(packaging.id) END AS id FROM packaging) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_pallets_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_pallets_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(pallets.id) END AS id FROM pallets) + 1;
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
-- Name: set_payment_transaction_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_payment_transaction_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(payment_transaction.id) END AS id FROM payment_transaction) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_payments_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_payments_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(payments.id) END AS id FROM payments) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_pos_terminals_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_pos_terminals_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(pos_terminals.id) END AS id FROM pos_terminals) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_product_account_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_product_account_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(product_account.id) END AS id FROM product_account) + 1;
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
-- Name: set_product_image_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_product_image_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(product_image.id) END AS id FROM product_image) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_purchase_delivery_note_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_purchase_delivery_note_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(purchase_delivery_note.id) END AS id FROM purchase_delivery_note) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_purchase_invoice_details_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_purchase_invoice_details_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(purchase_invoice_details.id) END AS id FROM purchase_invoice_details) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_purchase_invoice_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_purchase_invoice_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(purchase_invoice.id) END AS id FROM purchase_invoice) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_purchase_order_detail_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_purchase_order_detail_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(purchase_order_detail.id) END AS id FROM purchase_order_detail) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_purchase_order_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_purchase_order_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(purchase_order.id) END AS id FROM purchase_order) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_sales_delivery_note_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_sales_delivery_note_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_delivery_note.id) END AS id FROM sales_delivery_note) + 1;
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
-- Name: set_sales_order_detail_digital_product_data_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_sales_order_detail_digital_product_data_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_order_detail_digital_product_data.id) END AS id FROM sales_order_detail_digital_product_data) + 1;
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


--
-- Name: set_shipping_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_shipping_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(shipping.id) END AS id FROM shipping) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_shipping_status_history_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_shipping_status_history_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(shipping_status_history.id) END AS id FROM shipping_status_history) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_shipping_tag_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_shipping_tag_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(shipping_tag.id) END AS id FROM shipping_tag) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_state_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_state_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(state.id) END AS id FROM state) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_suppliers_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_suppliers_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(suppliers.id) END AS id FROM suppliers) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_transactional_log_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_transactional_log_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(transactional_log.id) END AS id FROM transactional_log) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_transfer_between_warehouses_detail_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_transfer_between_warehouses_detail_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(transfer_between_warehouses_detail.id) END AS id FROM transfer_between_warehouses_detail) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_transfer_between_warehouses_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_transfer_between_warehouses_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(transfer_between_warehouses.id) END AS id FROM transfer_between_warehouses) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_user_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_user_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX("user".id) END AS id FROM "user") + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_warehouse_movement_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_warehouse_movement_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(warehouse_movement.id) END AS id FROM warehouse_movement) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_webhook_logs_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_webhook_logs_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(webhook_logs.id) END AS id FROM webhook_logs) + 1;
    RETURN NEW;
END;
$$;


--
-- Name: set_webhook_settings_id(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_webhook_settings_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(webhook_settings.id) END AS id FROM webhook_settings) + 1;
    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: account; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.account (
    id integer NOT NULL,
    journal integer NOT NULL,
    name character varying(150) NOT NULL,
    credit real DEFAULT 0 NOT NULL,
    debit real DEFAULT 0 NOT NULL,
    balance real DEFAULT 0 NOT NULL,
    account_number integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: accounting_movement; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.accounting_movement (
    id bigint NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    amount_debit real DEFAULT 0 NOT NULL,
    amount_credit real DEFAULT 0 NOT NULL,
    fiscal_year smallint NOT NULL,
    type character(1) NOT NULL,
    billing_serie character(3) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: accounting_movement_detail; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.accounting_movement_detail (
    id bigint NOT NULL,
    movement bigint NOT NULL,
    journal integer NOT NULL,
    account integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    credit real NOT NULL,
    debit real NOT NULL,
    type character(1) NOT NULL,
    note character varying(300) NOT NULL,
    document_name character(15) NOT NULL,
    payment_method integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: address; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.address (
    id integer NOT NULL,
    customer integer,
    address character varying(200) NOT NULL,
    address_2 character varying(200) NOT NULL,
    state integer,
    city character varying(100) NOT NULL,
    country integer NOT NULL,
    private_business character(1) NOT NULL,
    notes text NOT NULL,
    supplier integer,
    ps_id integer DEFAULT 0 NOT NULL,
    zip_code character varying(12) NOT NULL,
    sy_id bigint DEFAULT 0 NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: api_key; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.api_key (
    id integer NOT NULL,
    name character varying(64) NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    user_created integer NOT NULL,
    off boolean DEFAULT false NOT NULL,
    "user" integer NOT NULL,
    token uuid,
    enterprise integer NOT NULL,
    auth character(1) DEFAULT 'P'::bpchar NOT NULL,
    basic_auth_user character varying(20),
    basic_auth_password character varying(20),
    permissions json DEFAULT '{"saleOrders":{"get":true,"post":true,"put":true,"delete":true},"saleOrderDetails":{"get":true,"post":true,"put":true,"delete":true},"saleOrderDetailsDigitalProductData":{"get":true,"post":true,"put":true,"delete":true},"saleInvoices":{"get":true,"post":true,"put":true,"delete":true},"saleInvoiceDetails":{"get":true,"post":true,"put":true,"delete":true},"saleDeliveryNotes":{"get":true,"post":true,"put":true,"delete":true},"purchaseOrders":{"get":true,"post":true,"put":true,"delete":true},"purchaseOrderDetails":{"get":true,"post":true,"put":true,"delete":true},"purchaseInvoices":{"get":true,"post":true,"put":true,"delete":true},"purchaseInvoiceDetails":{"get":true,"post":true,"put":true,"delete":true},"purchaseDeliveryNotes":{"get":true,"post":true,"put":true,"delete":true},"customers":{"get":true,"post":true,"put":true,"delete":true},"suppliers":{"get":true,"post":true,"put":true,"delete":true},"products":{"get":true,"post":true,"put":true,"delete":true},"countries":{"get":true,"post":true,"put":true,"delete":true},"states":{"get":true,"post":true,"put":true,"delete":true},"colors":{"get":true,"post":true,"put":true,"delete":true},"productFamilies":{"get":true,"post":true,"put":true,"delete":true},"addresses":{"get":true,"post":true,"put":true,"delete":true},"carriers":{"get":true,"post":true,"put":true,"delete":true},"billingSeries":{"get":true,"post":true,"put":true,"delete":true},"currencies":{"get":true,"post":true,"put":true,"delete":true},"paymentMethods":{"get":true,"post":true,"put":true,"delete":true},"languages":{"get":true,"post":true,"put":true,"delete":true},"packages":{"get":true,"post":true,"put":true,"delete":true},"incoterms":{"get":true,"post":true,"put":true,"delete":true},"warehouses":{"get":true,"post":true,"put":true,"delete":true},"warehouseMovements":{"get":true,"post":true,"put":true,"delete":true},"manufacturingOrders":{"get":true,"post":true,"put":true,"delete":true},"manufacturingOrderTypes":{"get":true,"post":true,"put":true,"delete":true},"complexManufacturingOrders":{"get":true,"post":true,"put":true,"delete":true},"manufacturingOrderTypeComponents":{"get":true,"post":true,"put":true,"delete":true},"shippings":{"get":true,"post":true,"put":true,"delete":true},"shippingStatusHistory":{"get":true,"post":true,"put":true,"delete":true},"stock":{"get":true,"post":true,"put":true,"delete":true},"journal":{"get":true,"post":true,"put":true,"delete":true},"account":{"get":true,"post":true,"put":true,"delete":true},"accountingMovement":{"get":true,"post":true,"put":true,"delete":true},"accountingMovementDetail":{"get":true,"post":true,"put":true,"delete":true},"collectionOperation":{"get":true,"post":true,"put":true,"delete":true},"charges":{"get":true,"post":true,"put":true,"delete":true},"paymentTransaction":{"get":true,"post":true,"put":true,"delete":true},"payment":{"get":true,"post":true,"put":true,"delete":true},"postSaleInvoice":{"get":true,"post":true,"put":true,"delete":true},"postPurchaseInvoice":{"get":true,"post":true,"put":true,"delete":true}}'::json NOT NULL
);


--
-- Name: billing_series; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.billing_series (
    id character(3) NOT NULL,
    name character varying(50) NOT NULL,
    billing_type character(1) NOT NULL,
    year smallint NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: carrier; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.carrier (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    max_weight real NOT NULL,
    max_width real NOT NULL,
    max_height real NOT NULL,
    max_depth real NOT NULL,
    max_packages smallint NOT NULL,
    phone character varying(15) NOT NULL,
    email character varying(100) NOT NULL,
    web character varying(100) NOT NULL,
    off boolean NOT NULL,
    ps_id integer DEFAULT 0 NOT NULL,
    pallets boolean NOT NULL,
    webservice character(1) NOT NULL,
    sendcloud_url character varying(75) DEFAULT ''::character varying NOT NULL,
    sendcloud_key character varying(32) DEFAULT ''::character varying NOT NULL,
    sendcloud_secret character varying(32) DEFAULT ''::character varying NOT NULL,
    sendcloud_shipping_method integer DEFAULT 0 NOT NULL,
    sendcloud_sender_address bigint DEFAULT 0 NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: charges; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.charges (
    id integer NOT NULL,
    accounting_movement bigint NOT NULL,
    accounting_movement_detail_debit bigint NOT NULL,
    accounting_movement_detail_credit bigint NOT NULL,
    account integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    amount real NOT NULL,
    concept character varying(140) NOT NULL,
    collection_operation integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: collection_operation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.collection_operation (
    id integer NOT NULL,
    accounting_movement bigint NOT NULL,
    accounting_movement_detail bigint NOT NULL,
    account integer NOT NULL,
    bank integer,
    status character(1) DEFAULT 'P'::bpchar NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_expiration timestamp(3) with time zone NOT NULL,
    total real NOT NULL,
    paid real NOT NULL,
    pending real NOT NULL,
    document_name character(15) NOT NULL,
    payment_method integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: color; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.color (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    hex_color character(6) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: complex_manufacturing_order; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.complex_manufacturing_order (
    id bigint NOT NULL,
    type integer NOT NULL,
    manufactured boolean DEFAULT false NOT NULL,
    date_manufactured timestamp(3) with time zone,
    user_manufactured integer,
    enterprise integer NOT NULL,
    quantity_pending_manufacture integer DEFAULT 0 NOT NULL,
    quantity_manufactured integer DEFAULT 0 NOT NULL,
    warehouse character(2) NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    uuid uuid NOT NULL,
    user_created integer NOT NULL,
    tag_printed boolean DEFAULT false NOT NULL,
    date_tag_printed timestamp(3) with time zone,
    user_tag_printed integer
);


--
-- Name: complex_manufacturing_order_manufacturing_order; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.complex_manufacturing_order_manufacturing_order (
    id bigint NOT NULL,
    manufacturing_order bigint,
    type character(1) NOT NULL,
    complex_manufacturing_order bigint NOT NULL,
    enterprise integer NOT NULL,
    warehouse_movement bigint,
    manufactured boolean DEFAULT false NOT NULL,
    product integer NOT NULL,
    manufacturing_order_type_component integer NOT NULL,
    purchase_order_detail bigint,
    sale_order_detail bigint,
    complex_manufacturing_order_manufacturing_order_output bigint
);


--
-- Name: config; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.config (
    id integer NOT NULL,
    default_vat_percent real NOT NULL,
    default_warehouse character(2),
    date_format character varying(25) NOT NULL,
    enterprise_name character varying(50) NOT NULL,
    enterprise_description character varying(250) NOT NULL,
    ecommerce character(1) NOT NULL,
    email character(1) NOT NULL,
    currency character(1) NOT NULL,
    currency_ecb_url character varying(100) NOT NULL,
    barcode_prefix character varying(4) NOT NULL,
    prestashop_url character varying(100) NOT NULL,
    prestashop_api_key character varying(32) NOT NULL,
    prestashop_language_id integer NOT NULL,
    prestashop_export_serie character(3),
    prestashop_intracommunity_serie character(3),
    prestashop_interior_serie character(3),
    cron_currency character varying(25) NOT NULL,
    cron_prestashop character varying(25) NOT NULL,
    sendgrid_key character varying(75) NOT NULL,
    email_from character varying(50) NOT NULL,
    name_from character varying(50) NOT NULL,
    pallet_weight real NOT NULL,
    pallet_width real NOT NULL,
    pallet_height real NOT NULL,
    pallet_depth real NOT NULL,
    max_connections integer NOT NULL,
    prestashop_status_payment_accepted integer DEFAULT 0 NOT NULL,
    prestashop_status_shipped integer DEFAULT 0 NOT NULL,
    minimum_stock_sales_periods smallint DEFAULT 0 NOT NULL,
    minimum_stock_sales_days smallint DEFAULT 0 NOT NULL,
    customer_journal integer DEFAULT 430,
    sales_journal integer DEFAULT 700,
    sales_account integer,
    supplier_journal integer,
    purchase_journal integer,
    purchase_account integer,
    enable_api_key boolean DEFAULT false NOT NULL,
    cron_clear_labels character varying(25) DEFAULT '@midnight'::character varying NOT NULL,
    limit_accounting_date timestamp(0) with time zone,
    woocommerce_url character varying(100) DEFAULT ''::character varying NOT NULL,
    woocommerce_consumer_key character varying(50) DEFAULT ''::character varying NOT NULL,
    woocommerce_consumer_secret character varying(50) DEFAULT ''::character varying NOT NULL,
    woocommerce_export_serie character(3),
    woocommerce_intracommunity_serie character(3),
    woocommerce_interior_serie character(3),
    woocommerce_default_payment_method integer,
    connection_log boolean DEFAULT false NOT NULL,
    filter_connections boolean DEFAULT false NOT NULL,
    shopify_url character varying(100) DEFAULT ''::character varying NOT NULL,
    shopify_token character varying(50) DEFAULT ''::character varying NOT NULL,
    shopify_export_serie character(3),
    shopify_intracommunity_serie character(3),
    shopify_interior_serie character(3),
    shopify_default_payment_method integer,
    shopify_shop_location_id bigint DEFAULT 0 NOT NULL,
    enterprise_key character varying(25) DEFAULT ''::character varying NOT NULL,
    password_minimum_length smallint DEFAULT 8 NOT NULL,
    password_minumum_complexity character(1) DEFAULT 'B'::bpchar NOT NULL,
    invoice_delete_policy smallint DEFAULT 2 NOT NULL,
    transaction_log boolean DEFAULT true NOT NULL,
    undo_manufacturing_order_seconds smallint DEFAULT 120 NOT NULL,
    cron_sendcloud_tracking character varying(25) DEFAULT ''::character varying NOT NULL,
    smtp_identity character varying(50) DEFAULT ''::character varying NOT NULL,
    smtp_username character varying(50) DEFAULT ''::character varying NOT NULL,
    smtp_password character varying(50) DEFAULT ''::character varying NOT NULL,
    smtp_hostname character varying(50) DEFAULT ''::character varying NOT NULL,
    smtp_starttls boolean DEFAULT false NOT NULL,
    smtp_reply_to character varying(50) DEFAULT ''::character varying,
    email_send_error_ecommerce character varying(150) DEFAULT ''::character varying NOT NULL,
    email_send_error_sendcloud character varying(150) DEFAULT ''::character varying NOT NULL,
    product_barcode_label_width smallint DEFAULT 0 NOT NULL,
    product_barcode_label_height smallint DEFAULT 0 NOT NULL,
    product_barcode_label_size smallint DEFAULT 0 NOT NULL,
    product_barcode_label_margin_top smallint DEFAULT 0 NOT NULL,
    product_barcode_label_margin_bottom smallint DEFAULT 0 NOT NULL,
    product_barcode_label_margin_left smallint DEFAULT 0 NOT NULL,
    product_barcode_label_margin_right smallint DEFAULT 0 NOT NULL
);


--
-- Name: config_accounts_vat; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.config_accounts_vat (
    vat_percent real NOT NULL,
    account_sale integer NOT NULL,
    account_purchase integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: connection_filter; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.connection_filter (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    type character(1) NOT NULL,
    ip_address character varying(15),
    time_start time(0) with time zone,
    time_end time(0) with time zone,
    enterprise integer NOT NULL
);


--
-- Name: connection_filter_user; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.connection_filter_user (
    connection_filter integer NOT NULL,
    "user" integer NOT NULL
);


--
-- Name: connection_log; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.connection_log (
    id bigint NOT NULL,
    date_connected timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_disconnected timestamp(3) with time zone,
    "user" integer NOT NULL,
    ok boolean NOT NULL,
    ip_address character varying(15) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: country; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.country (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    iso_2 character(2) NOT NULL,
    iso_3 character(3) NOT NULL,
    un_code smallint NOT NULL,
    zone character(1) NOT NULL,
    phone_prefix smallint NOT NULL,
    language integer,
    currency integer,
    enterprise integer NOT NULL
);


--
-- Name: currency; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.currency (
    id integer NOT NULL,
    name character varying(150) NOT NULL,
    sign character(3) NOT NULL,
    iso_code character(3) NOT NULL,
    iso_num smallint NOT NULL,
    exchange real NOT NULL,
    exchange_date date DEFAULT CURRENT_DATE NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: custom_fields; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.custom_fields (
    id bigint NOT NULL,
    enterprise integer NOT NULL,
    product integer,
    customer integer,
    supplier integer,
    name character varying(255) NOT NULL,
    field_type smallint,
    value_string text,
    value_number real,
    value_boolean boolean,
    value_binary bytea,
    file_name character varying(255),
    file_size integer,
    image_mime_type character varying(255)
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
    country integer,
    state integer,
    main_shipping_address integer,
    main_billing_address integer,
    language integer,
    payment_method integer,
    billing_series character(3),
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    ps_id integer DEFAULT 0 NOT NULL,
    account integer,
    wc_id integer DEFAULT 0 NOT NULL,
    sy_id bigint DEFAULT 0 NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: document; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.document (
    id integer NOT NULL,
    name character varying(250) NOT NULL,
    uuid uuid NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_updated timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    size integer DEFAULT 0 NOT NULL,
    container integer NOT NULL,
    dsc text NOT NULL,
    sales_order bigint,
    sales_invoice bigint,
    sales_delivery_note bigint,
    shipping bigint,
    purchase_order bigint,
    purchase_invoice bigint,
    purchase_delivery_note bigint,
    mime_type character varying(100) DEFAULT ''::character varying NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: document_container; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.document_container (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    path character varying(250) NOT NULL,
    max_file_size integer NOT NULL,
    disallowed_mime_types character varying(250) NOT NULL,
    allowed_mime_types character varying(250) NOT NULL,
    enterprise integer NOT NULL,
    used_storage bigint DEFAULT 0 NOT NULL,
    max_storage bigint DEFAULT 0 NOT NULL
);


--
-- Name: email_log; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.email_log (
    id bigint NOT NULL,
    email_from character varying(100) NOT NULL,
    name_from character varying(100) NOT NULL,
    destination_email character varying(100) NOT NULL,
    destination_name character varying(100) NOT NULL,
    subject character varying(100) NOT NULL,
    content text NOT NULL,
    date_sent timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: enterprise_logo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.enterprise_logo (
    enterprise integer NOT NULL,
    logo bytea NOT NULL,
    mime_type character varying(150) NOT NULL
);


--
-- Name: group; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."group" (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    sales boolean NOT NULL,
    purchases boolean NOT NULL,
    masters boolean NOT NULL,
    warehouse boolean NOT NULL,
    manufacturing boolean NOT NULL,
    preparation boolean NOT NULL,
    admin boolean NOT NULL,
    prestashop boolean NOT NULL,
    accounting boolean DEFAULT false NOT NULL,
    enterprise integer NOT NULL,
    point_of_sale boolean DEFAULT false NOT NULL
);


--
-- Name: hs_codes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.hs_codes (
    id character varying(8) NOT NULL,
    name character varying(255) NOT NULL
);


--
-- Name: incoterm; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.incoterm (
    id integer NOT NULL,
    key character(3) NOT NULL,
    name character varying(50) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: inventory; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.inventory (
    id integer NOT NULL,
    enterprise integer NOT NULL,
    name character varying(50) NOT NULL,
    date_created timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    finished boolean DEFAULT false NOT NULL,
    date_finished timestamp(3) without time zone,
    warehouse character(2) NOT NULL
);


--
-- Name: inventory_products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.inventory_products (
    inventory integer NOT NULL,
    product integer NOT NULL,
    enterprise integer NOT NULL,
    quantity integer NOT NULL,
    warehouse_movement bigint
);


--
-- Name: journal; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.journal (
    id integer NOT NULL,
    name character varying(75) NOT NULL,
    type character(1) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: language; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    iso_2 character(2) NOT NULL,
    iso_3 character(3) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: login_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.login_tokens (
    id integer NOT NULL,
    name character(128) NOT NULL,
    date_last_used timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    "user" integer NOT NULL,
    ip_address character varying(15) NOT NULL
);


--
-- Name: logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.logs (
    id bigint NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    title character varying(255) NOT NULL,
    info text NOT NULL,
    stacktrace text NOT NULL
);


--
-- Name: manufacturing_order; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.manufacturing_order (
    id bigint NOT NULL,
    order_detail bigint,
    product integer NOT NULL,
    type integer NOT NULL,
    uuid uuid NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_last_update timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    manufactured boolean DEFAULT false NOT NULL,
    date_manufactured timestamp(3) with time zone,
    user_manufactured integer,
    user_created integer NOT NULL,
    tag_printed boolean DEFAULT false NOT NULL,
    date_tag_printed timestamp(3) with time zone,
    "order" bigint,
    user_tag_printed integer,
    enterprise integer NOT NULL,
    warehouse character(2) DEFAULT 'W1'::bpchar NOT NULL,
    warehouse_movement bigint,
    quantity_manufactured integer DEFAULT 1 NOT NULL,
    complex boolean DEFAULT false NOT NULL
);


--
-- Name: manufacturing_order_type; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.manufacturing_order_type (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    enterprise integer NOT NULL,
    quantity_manufactured integer DEFAULT 1 NOT NULL,
    complex boolean DEFAULT false NOT NULL
);


--
-- Name: manufacturing_order_type_components; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.manufacturing_order_type_components (
    id integer NOT NULL,
    manufacturing_order_type integer NOT NULL,
    type character(1) NOT NULL,
    product integer NOT NULL,
    quantity integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: packages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.packages (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    weight real NOT NULL,
    width real NOT NULL,
    height real NOT NULL,
    depth real NOT NULL,
    product integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: packaging; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.packaging (
    id bigint NOT NULL,
    package integer NOT NULL,
    sales_order bigint NOT NULL,
    weight real NOT NULL,
    shipping bigint,
    pallet integer,
    enterprise integer NOT NULL
);


--
-- Name: pallets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pallets (
    id integer NOT NULL,
    sales_order bigint NOT NULL,
    weight real NOT NULL,
    width real NOT NULL,
    height real NOT NULL,
    depth real NOT NULL,
    name character varying(40) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: payment_method; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payment_method (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    paid_in_advance boolean NOT NULL,
    prestashop_module_name character varying(100) NOT NULL,
    days_expiration smallint DEFAULT 0 NOT NULL,
    bank integer,
    woocommerce_module_name character varying(100) DEFAULT ''::character varying NOT NULL,
    shopify_module_name character varying(100) DEFAULT ''::character varying NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: payment_transaction; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payment_transaction (
    id integer NOT NULL,
    accounting_movement bigint NOT NULL,
    accounting_movement_detail bigint NOT NULL,
    account integer NOT NULL,
    bank integer,
    status character(1) DEFAULT 'P'::bpchar NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_expiration timestamp(3) with time zone NOT NULL,
    total real NOT NULL,
    paid real NOT NULL,
    pending real NOT NULL,
    document_name character(15) NOT NULL,
    payment_method integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: payments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payments (
    id integer NOT NULL,
    accounting_movement bigint NOT NULL,
    accounting_movement_detail_debit bigint NOT NULL,
    accounting_movement_detail_credit bigint NOT NULL,
    account integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    amount real NOT NULL,
    concept character varying(140) NOT NULL,
    payment_transaction integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: permission_dictionary; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.permission_dictionary (
    enterprise integer NOT NULL,
    key character varying(150) NOT NULL,
    description character varying(250) NOT NULL
);


--
-- Name: permission_dictionary_group; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.permission_dictionary_group (
    "group" integer NOT NULL,
    permission_key character varying(150) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: pos_terminals; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pos_terminals (
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
    enterprise integer NOT NULL
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
    family integer,
    width real NOT NULL,
    height real NOT NULL,
    depth real NOT NULL,
    off boolean NOT NULL,
    stock integer NOT NULL,
    vat_percent real NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    dsc text NOT NULL,
    color integer,
    price real NOT NULL,
    manufacturing boolean NOT NULL,
    manufacturing_order_type integer,
    supplier integer,
    ps_id integer DEFAULT 0 NOT NULL,
    ps_combination_id integer DEFAULT 0 NOT NULL,
    minimum_stock integer DEFAULT 0 NOT NULL,
    track_minimum_stock boolean DEFAULT false NOT NULL,
    wc_id integer DEFAULT 0 NOT NULL,
    wc_variation_id integer DEFAULT 0 NOT NULL,
    sy_id bigint DEFAULT 0 NOT NULL,
    sy_variant_id bigint DEFAULT 0 NOT NULL,
    enterprise integer NOT NULL,
    digital_product boolean DEFAULT false NOT NULL,
    purchase_price real DEFAULT 0 NOT NULL,
    minimum_purchase_quantity integer DEFAULT 0 NOT NULL,
    origin_country character varying(2) DEFAULT ''::character varying NOT NULL,
    hs_code character varying(8),
    cost_price real DEFAULT 0 NOT NULL
);


--
-- Name: product_account; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_account (
    id integer NOT NULL,
    product integer NOT NULL,
    account integer NOT NULL,
    jorunal integer NOT NULL,
    account_number integer NOT NULL,
    enterprise integer NOT NULL,
    type character(1) NOT NULL
);


--
-- Name: product_family; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_family (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    reference character varying(40) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: product_image; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_image (
    id integer NOT NULL,
    product integer NOT NULL,
    url character varying(255) NOT NULL
);


--
-- Name: product_translation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_translation (
    product integer NOT NULL,
    language integer NOT NULL,
    name character varying(150) NOT NULL
);


--
-- Name: ps_address; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_address (
    id integer NOT NULL,
    id_country integer NOT NULL,
    id_state integer NOT NULL,
    id_customer integer NOT NULL,
    alias character varying(32) NOT NULL,
    company character varying(255) NOT NULL,
    lastname character varying(255) NOT NULL,
    firstname character varying(255) NOT NULL,
    address1 character varying(128) NOT NULL,
    address2 character varying(128) NOT NULL,
    postcode character varying(12) NOT NULL,
    city character varying(64) NOT NULL,
    other text NOT NULL,
    phone character varying(32) NOT NULL,
    phone_mobile character varying(32) NOT NULL,
    vat_number character varying(32) NOT NULL,
    dni character varying(16) NOT NULL,
    date_add timestamp(0) with time zone NOT NULL,
    date_upd timestamp(0) with time zone NOT NULL,
    deleted boolean NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_carrier; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_carrier (
    id integer NOT NULL,
    deleted boolean NOT NULL,
    name character varying(64) NOT NULL,
    active boolean NOT NULL,
    url character varying(255) NOT NULL,
    max_width real NOT NULL,
    max_height real NOT NULL,
    max_depth real NOT NULL,
    max_weight real NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_country; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_country (
    id integer NOT NULL,
    id_zone integer NOT NULL,
    id_currency integer NOT NULL,
    iso_code character varying(3) NOT NULL,
    call_prefix integer NOT NULL,
    active boolean NOT NULL,
    name character varying(64) NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_currency; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_currency (
    id integer NOT NULL,
    name character varying(64) NOT NULL,
    iso_code character varying(3) NOT NULL,
    conversion_rate real NOT NULL,
    deleted boolean NOT NULL,
    active boolean NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    symbol character varying(3) NOT NULL,
    numeric_iso_code integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_customer; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_customer (
    id integer NOT NULL,
    id_lang integer NOT NULL,
    company character varying(255) NOT NULL,
    firstname character varying(255) NOT NULL,
    lastname character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    note text NOT NULL,
    active boolean NOT NULL,
    deleted boolean NOT NULL,
    date_add timestamp(0) with time zone NOT NULL,
    date_upd timestamp(0) with time zone NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_language; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_language (
    id integer NOT NULL,
    name character varying(32) NOT NULL,
    iso_code character varying(2) NOT NULL,
    active boolean NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_order; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_order (
    id integer NOT NULL,
    reference character varying(9) NOT NULL,
    id_carrier integer NOT NULL,
    id_lang integer NOT NULL,
    id_customer integer NOT NULL,
    id_currency integer NOT NULL,
    id_address_delivery integer NOT NULL,
    id_address_invoice integer NOT NULL,
    module character varying(255) NOT NULL,
    total_discounts_tax_excl real NOT NULL,
    total_shipping_tax_excl real NOT NULL,
    date_add timestamp(0) with time zone NOT NULL,
    date_upd timestamp(0) with time zone NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    tax_included boolean NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_order_detail; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_order_detail (
    id integer NOT NULL,
    id_order integer NOT NULL,
    product_id integer NOT NULL,
    product_attribute_id integer NOT NULL,
    product_quantity integer NOT NULL,
    product_price real NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_product; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_product (
    id integer NOT NULL,
    on_sale boolean NOT NULL,
    ean13 character varying(13) NOT NULL,
    price real NOT NULL,
    reference character varying(64) NOT NULL,
    active boolean NOT NULL,
    date_add timestamp(0) with time zone NOT NULL,
    date_upd timestamp(0) with time zone NOT NULL,
    name character varying(128) NOT NULL,
    description text NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_product_combination; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_product_combination (
    id integer NOT NULL,
    id_product integer NOT NULL,
    reference character varying(64) NOT NULL,
    ean13 character varying(13) NOT NULL,
    price real NOT NULL,
    product_option_values integer[] NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_product_option_values; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_product_option_values (
    id integer NOT NULL,
    name character varying(128) NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_state; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_state (
    id integer NOT NULL,
    id_country integer NOT NULL,
    id_zone integer NOT NULL,
    name character varying(64) NOT NULL,
    iso_code character varying(7) NOT NULL,
    active boolean NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: ps_zone; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ps_zone (
    id integer NOT NULL,
    name character varying(64) NOT NULL,
    active boolean NOT NULL,
    ps_exists boolean DEFAULT true NOT NULL,
    zone character(1) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: purchase_delivery_note; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_delivery_note (
    id bigint NOT NULL,
    warehouse character(2) NOT NULL,
    supplier integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    payment_method integer NOT NULL,
    billing_series character(3) NOT NULL,
    shipping_address integer NOT NULL,
    total_products real DEFAULT 0 NOT NULL,
    discount_percent real NOT NULL,
    fix_discount real NOT NULL,
    shipping_price real NOT NULL,
    shipping_discount real NOT NULL,
    total_with_discount real DEFAULT 0 NOT NULL,
    total_vat real DEFAULT 0 NOT NULL,
    total_amount real DEFAULT 0 NOT NULL,
    lines_number smallint DEFAULT 0 NOT NULL,
    delivery_note_name character(15) NOT NULL,
    delivery_note_number integer NOT NULL,
    currency integer NOT NULL,
    currency_change real NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: purchase_invoice; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_invoice (
    id bigint NOT NULL,
    supplier integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    payment_method integer NOT NULL,
    billing_series character(3) NOT NULL,
    currency integer NOT NULL,
    currency_change real NOT NULL,
    billing_address integer NOT NULL,
    total_products real DEFAULT 0 NOT NULL,
    discount_percent real NOT NULL,
    fix_discount real NOT NULL,
    shipping_price real NOT NULL,
    shipping_discount real NOT NULL,
    total_with_discount real DEFAULT 0 NOT NULL,
    vat_amount real DEFAULT 0 NOT NULL,
    total_amount real DEFAULT 0 NOT NULL,
    lines_number smallint DEFAULT 0 NOT NULL,
    invoice_number integer NOT NULL,
    invoice_name character(15) NOT NULL,
    accounting_movement bigint,
    enterprise integer NOT NULL,
    amending boolean DEFAULT false NOT NULL,
    amended_invoice bigint,
    income_tax boolean DEFAULT false NOT NULL,
    income_tax_base real DEFAULT 0 NOT NULL,
    income_tax_percentage real DEFAULT 0 NOT NULL,
    income_tax_value real DEFAULT 0 NOT NULL,
    rent boolean DEFAULT false NOT NULL,
    rent_base real DEFAULT 0 NOT NULL,
    rent_percentage real DEFAULT 0 NOT NULL,
    rent_value real DEFAULT 0 NOT NULL
);


--
-- Name: purchase_invoice_details; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_invoice_details (
    id bigint NOT NULL,
    invoice bigint NOT NULL,
    product integer,
    price real NOT NULL,
    quantity integer NOT NULL,
    vat_percent real NOT NULL,
    total_amount real NOT NULL,
    order_detail integer,
    enterprise integer NOT NULL,
    description character varying(150) DEFAULT ''::character varying NOT NULL,
    income_tax boolean DEFAULT false NOT NULL,
    rent boolean DEFAULT false NOT NULL
);


--
-- Name: purchase_order; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_order (
    id bigint NOT NULL,
    warehouse character(2) NOT NULL,
    supplier_reference character varying(40) NOT NULL,
    supplier integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_paid timestamp(3) with time zone,
    payment_method integer NOT NULL,
    billing_series character(3) NOT NULL,
    currency integer NOT NULL,
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
    total_vat real DEFAULT 0 NOT NULL,
    total_amount real DEFAULT 0 NOT NULL,
    dsc text NOT NULL,
    notes character varying(250) NOT NULL,
    off boolean DEFAULT false NOT NULL,
    cancelled boolean DEFAULT false NOT NULL,
    order_number integer NOT NULL,
    billing_status character(1) DEFAULT 'P'::bpchar NOT NULL,
    order_name character(15) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: purchase_order_detail; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_order_detail (
    id bigint NOT NULL,
    "order" bigint NOT NULL,
    product integer NOT NULL,
    price real NOT NULL,
    quantity integer NOT NULL,
    vat_percent real NOT NULL,
    total_amount real NOT NULL,
    quantity_invoiced integer DEFAULT 0 NOT NULL,
    quantity_delivery_note integer DEFAULT 0 NOT NULL,
    quantity_pending_packaging integer NOT NULL,
    quantity_assigned_sale integer DEFAULT 0 NOT NULL,
    enterprise integer NOT NULL,
    cancelled boolean DEFAULT false NOT NULL
);


--
-- Name: pwd_blacklist; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pwd_blacklist (
    pwd character varying(255) NOT NULL
);


--
-- Name: pwd_sha1_blacklist; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pwd_sha1_blacklist (
    hash bytea NOT NULL
);


--
-- Name: report_template; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.report_template (
    enterprise integer NOT NULL,
    key character varying(50) NOT NULL,
    html text NOT NULL
);


--
-- Name: report_template_translation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.report_template_translation (
    enterprise integer NOT NULL,
    key character varying(50) NOT NULL,
    language integer NOT NULL,
    translation character varying(255) NOT NULL
);


--
-- Name: sales_delivery_note; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_delivery_note (
    id bigint NOT NULL,
    warehouse character(2) NOT NULL,
    customer integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    payment_method integer,
    billing_series character(3) NOT NULL,
    shipping_address integer NOT NULL,
    total_products real DEFAULT 0 NOT NULL,
    discount_percent real DEFAULT 0 NOT NULL,
    fix_discount real DEFAULT 0 NOT NULL,
    shipping_price real DEFAULT 0 NOT NULL,
    shipping_discount real DEFAULT 0 NOT NULL,
    total_with_discount real DEFAULT 0 NOT NULL,
    vat_amount real DEFAULT 0 NOT NULL,
    total_amount real DEFAULT 0 NOT NULL,
    lines_number smallint DEFAULT 0 NOT NULL,
    delivery_note_name character(15) NOT NULL,
    delivery_note_number integer NOT NULL,
    currency integer NOT NULL,
    currency_change real NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: sales_invoice; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_invoice (
    id bigint NOT NULL,
    customer integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    payment_method integer NOT NULL,
    billing_series character(3) NOT NULL,
    currency integer NOT NULL,
    currency_change real NOT NULL,
    billing_address integer NOT NULL,
    total_products real DEFAULT 0 NOT NULL,
    discount_percent real NOT NULL,
    fix_discount real NOT NULL,
    shipping_price real NOT NULL,
    shipping_discount real NOT NULL,
    total_with_discount real DEFAULT 0 NOT NULL,
    vat_amount real DEFAULT 0 NOT NULL,
    total_amount real DEFAULT 0 NOT NULL,
    lines_number smallint DEFAULT 0 NOT NULL,
    invoice_number integer NOT NULL,
    invoice_name character(15) NOT NULL,
    accounting_movement bigint,
    enterprise integer NOT NULL,
    simplified_invoice boolean DEFAULT false NOT NULL,
    amending boolean DEFAULT false NOT NULL,
    amended_invoice bigint
);


--
-- Name: sales_invoice_detail; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_invoice_detail (
    id integer NOT NULL,
    invoice bigint NOT NULL,
    product integer,
    price real NOT NULL,
    quantity integer NOT NULL,
    vat_percent real NOT NULL,
    total_amount real NOT NULL,
    order_detail bigint,
    enterprise integer NOT NULL,
    description character varying(150) DEFAULT ''::character varying NOT NULL
);


--
-- Name: sales_order; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_order (
    id bigint NOT NULL,
    warehouse character(2) NOT NULL,
    reference character varying(15) NOT NULL,
    customer integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_payment_accepted timestamp(3) with time zone,
    payment_method integer NOT NULL,
    billing_series character(3) NOT NULL,
    currency integer NOT NULL,
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
    order_name character(15) NOT NULL,
    carrier integer,
    ps_id integer DEFAULT 0 NOT NULL,
    wc_id integer DEFAULT 0 NOT NULL,
    sy_id bigint DEFAULT 0 NOT NULL,
    sy_draft_id bigint DEFAULT 0 NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: sales_order_detail; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_order_detail (
    id bigint NOT NULL,
    "order" bigint NOT NULL,
    product integer NOT NULL,
    price real NOT NULL,
    quantity integer NOT NULL,
    vat_percent real NOT NULL,
    total_amount real NOT NULL,
    quantity_invoiced integer DEFAULT 0 NOT NULL,
    quantity_delivery_note integer DEFAULT 0 NOT NULL,
    status character(1) NOT NULL,
    quantity_pending_packaging integer NOT NULL,
    purchase_order_detail bigint,
    ps_id integer DEFAULT 0 NOT NULL,
    cancelled boolean DEFAULT false NOT NULL,
    wc_id integer DEFAULT 0 NOT NULL,
    sy_id bigint DEFAULT 0 NOT NULL,
    sy_draft_id bigint DEFAULT 0 NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: sales_order_detail_digital_product_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_order_detail_digital_product_data (
    id integer NOT NULL,
    detail bigint NOT NULL,
    key character varying(50) NOT NULL,
    value character varying(250) NOT NULL
);


--
-- Name: sales_order_detail_packaged; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_order_detail_packaged (
    order_detail bigint NOT NULL,
    packaging bigint NOT NULL,
    quantity integer NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: sales_order_discount; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales_order_discount (
    id integer NOT NULL,
    "order" bigint NOT NULL,
    name character varying(100) NOT NULL,
    value_tax_included real NOT NULL,
    value_tax_excluded real NOT NULL,
    enterprise integer NOT NULL,
    sales_invoice_detail integer
);


--
-- Name: shipping; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.shipping (
    id bigint NOT NULL,
    "order" bigint NOT NULL,
    delivery_note bigint NOT NULL,
    delivery_address integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_sent timestamp(3) with time zone,
    sent boolean DEFAULT false NOT NULL,
    collected boolean DEFAULT false NOT NULL,
    "national" boolean NOT NULL,
    shipping_number character varying(50) DEFAULT ''::character varying NOT NULL,
    tracking_number character varying(50) DEFAULT ''::character varying NOT NULL,
    carrier integer NOT NULL,
    weight real DEFAULT 0 NOT NULL,
    packages_number smallint DEFAULT 0 NOT NULL,
    incoterm integer,
    carrier_notes character varying(250) NOT NULL,
    description text NOT NULL,
    enterprise integer NOT NULL,
    delivered boolean DEFAULT false NOT NULL
);


--
-- Name: shipping_status_history; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.shipping_status_history (
    id bigint NOT NULL,
    shipping bigint NOT NULL,
    status_id smallint NOT NULL,
    message character varying(255) NOT NULL,
    delivered boolean NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL
);


--
-- Name: shipping_tag; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.shipping_tag (
    id bigint NOT NULL,
    shipping bigint NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    label bytea NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: state; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.state (
    id integer NOT NULL,
    country integer NOT NULL,
    name character varying(100) NOT NULL,
    iso_code character varying(7) NOT NULL,
    enterprise integer NOT NULL
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
    quantity_pending_manufacture integer DEFAULT 0 NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: suppliers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.suppliers (
    id integer NOT NULL,
    name character varying(303) NOT NULL,
    tradename character varying(150) NOT NULL,
    fiscal_name character varying(150) NOT NULL,
    tax_id character varying(25) NOT NULL,
    vat_number character varying(25) NOT NULL,
    phone character varying(15) NOT NULL,
    email character varying(100) NOT NULL,
    main_address integer,
    country integer,
    state integer,
    main_shipping_address integer,
    main_billing_address integer,
    language integer,
    payment_method integer,
    billing_series character(3),
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    account integer,
    enterprise integer NOT NULL
);


--
-- Name: sy_addresses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sy_addresses (
    id bigint NOT NULL,
    customer_id bigint NOT NULL,
    company character varying(100) NOT NULL,
    address1 character varying(100) NOT NULL,
    address2 character varying(100) NOT NULL,
    city character varying(50) NOT NULL,
    province character varying(50) NOT NULL,
    zip character varying(25) NOT NULL,
    country_code character varying(5) NOT NULL,
    sy_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: sy_customers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sy_customers (
    id bigint NOT NULL,
    email character varying(100) NOT NULL,
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    tax_exempt boolean NOT NULL,
    phone character varying(25) NOT NULL,
    currency character varying(5) NOT NULL,
    sy_exists boolean DEFAULT true NOT NULL,
    default_address_id bigint NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: sy_draft_order_line_item; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sy_draft_order_line_item (
    id bigint NOT NULL,
    variant_id bigint NOT NULL,
    product_id bigint NOT NULL,
    quantity integer NOT NULL,
    taxable boolean NOT NULL,
    price real NOT NULL,
    draft_order_id bigint NOT NULL,
    sy_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: sy_draft_orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sy_draft_orders (
    id bigint NOT NULL,
    currency character varying(5) NOT NULL,
    tax_exempt boolean NOT NULL,
    name character varying(9) NOT NULL,
    shipping_address_1 character varying(100) NOT NULL,
    shipping_address2 character varying(100) NOT NULL,
    shipping_address_city character varying(50) NOT NULL,
    shipping_address_zip character varying(25) NOT NULL,
    shipping_address_country_code character varying(5) NOT NULL,
    billing_address_1 character varying(100) NOT NULL,
    billing_address2 character varying(100) NOT NULL,
    billing_address_city character varying(50) NOT NULL,
    billing_address_zip character varying(25) NOT NULL,
    billing_address_country_code character varying(5) NOT NULL,
    total_tax real NOT NULL,
    customer_id bigint NOT NULL,
    sy_exists boolean DEFAULT true NOT NULL,
    order_id bigint,
    enterprise integer NOT NULL
);


--
-- Name: sy_order_line_item; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sy_order_line_item (
    id bigint NOT NULL,
    variant_id bigint NOT NULL,
    product_id bigint NOT NULL,
    quantity integer NOT NULL,
    taxable boolean NOT NULL,
    price real NOT NULL,
    order_id bigint NOT NULL,
    sy_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: sy_orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sy_orders (
    id bigint NOT NULL,
    currency character varying(5) NOT NULL,
    current_total_discounts real NOT NULL,
    total_shipping_price_set_amount real NOT NULL,
    total_shipping_price_set_currency_code character varying(5) NOT NULL,
    tax_exempt boolean NOT NULL,
    name character varying(9) NOT NULL,
    shipping_address_1 character varying(100) NOT NULL,
    shipping_address2 character varying(100) NOT NULL,
    shipping_address_city character varying(50) NOT NULL,
    shipping_address_zip character varying(25) NOT NULL,
    shipping_address_country_code character varying(5) NOT NULL,
    billing_address_1 character varying(100) NOT NULL,
    billing_address2 character varying(100) NOT NULL,
    billing_address_city character varying(50) NOT NULL,
    billing_address_zip character varying(25) NOT NULL,
    billing_address_country_code character varying(5) NOT NULL,
    total_tax real NOT NULL,
    customer_id bigint NOT NULL,
    sy_exists boolean DEFAULT true NOT NULL,
    gateway character varying(50) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: sy_products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sy_products (
    id bigint NOT NULL,
    title character varying(150) NOT NULL,
    body_html text NOT NULL,
    sy_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: sy_variants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sy_variants (
    id bigint NOT NULL,
    product_id bigint NOT NULL,
    title character varying(150) NOT NULL,
    price real NOT NULL,
    sku character varying(25) NOT NULL,
    option1 character varying(150) NOT NULL,
    option2 character varying(150),
    option3 character varying(150),
    taxable boolean NOT NULL,
    barcode character varying(25) NOT NULL,
    grams integer NOT NULL,
    sy_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: transactional_log; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.transactional_log (
    id bigint NOT NULL,
    enterprise integer NOT NULL,
    "table" character varying(150) NOT NULL,
    register jsonb NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    register_id bigint NOT NULL,
    "user" integer,
    mode character(1) NOT NULL
);


--
-- Name: transfer_between_warehouses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.transfer_between_warehouses (
    id bigint NOT NULL,
    warehouse_origin character(2) NOT NULL,
    warehouse_destination character(2) NOT NULL,
    enterprise integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_finished timestamp(3) without time zone,
    finished boolean DEFAULT false NOT NULL,
    lines_transfered integer DEFAULT 0 NOT NULL,
    lines_total integer DEFAULT 0 NOT NULL,
    name character varying(100) NOT NULL
);


--
-- Name: transfer_between_warehouses_detail; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.transfer_between_warehouses_detail (
    id bigint NOT NULL,
    transfer_between_warehouses bigint NOT NULL,
    enterprise integer NOT NULL,
    product integer NOT NULL,
    quantity integer NOT NULL,
    quantity_transferred integer DEFAULT 0 NOT NULL,
    finished boolean DEFAULT false NOT NULL,
    warehouse_movement_out bigint,
    warehouse_movement_in bigint
);


--
-- Name: user; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."user" (
    id integer NOT NULL,
    username character varying(40) NOT NULL,
    full_name character varying(150) NOT NULL,
    email character varying(100) DEFAULT ''::character varying NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    date_last_pwd timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    pwd_next_login boolean DEFAULT false NOT NULL,
    off boolean DEFAULT false NOT NULL,
    pwd bytea NOT NULL,
    salt character(30) NOT NULL,
    iterations integer NOT NULL,
    dsc text DEFAULT ''::text NOT NULL,
    date_last_login timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    failed_login_attemps smallint DEFAULT 0 NOT NULL,
    lang character(2) DEFAULT 'en'::bpchar NOT NULL,
    config integer NOT NULL,
    uses_google_authenticator boolean DEFAULT false NOT NULL,
    google_authenticator_secret character(8)
);


--
-- Name: user_group; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_group (
    "user" integer NOT NULL,
    "group" integer NOT NULL
);


--
-- Name: warehouse; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehouse (
    id character(2) NOT NULL,
    name character varying(50) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: warehouse_movement; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehouse_movement (
    id bigint NOT NULL,
    warehouse character(2) NOT NULL,
    product integer NOT NULL,
    quantity integer NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    type character(1) NOT NULL,
    sales_order bigint,
    sales_order_detail bigint,
    sales_delivery_note bigint,
    dsc text NOT NULL,
    purchase_order bigint,
    purchase_order_detail bigint,
    purchase_delivery_note bigint,
    dragged_stock integer DEFAULT 0 NOT NULL,
    price real DEFAULT 0 NOT NULL,
    vat_percent real DEFAULT 0 NOT NULL,
    total_amount real DEFAULT 0 NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: wc_customers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.wc_customers (
    id integer NOT NULL,
    date_created timestamp(0) with time zone NOT NULL,
    email character varying(100) NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    billing_address_1 character varying(255) NOT NULL,
    billing_address_2 character varying(255) NOT NULL,
    billing_city character varying(255) NOT NULL,
    billing_postcode character varying(255) NOT NULL,
    billing_country character varying(255) NOT NULL,
    billing_state character varying(255) NOT NULL,
    billing_phone character varying(255) NOT NULL,
    shipping_address_1 character varying(255) NOT NULL,
    shipping_address_2 character varying(255) NOT NULL,
    shipping_city character varying(255) NOT NULL,
    shipping_postcode character varying(255) NOT NULL,
    shipping_country character varying(255) NOT NULL,
    shipping_state character varying(255) NOT NULL,
    shipping_phone character varying(255) NOT NULL,
    wc_exists boolean DEFAULT true NOT NULL,
    billing_company character varying(255) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: wc_order_details; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.wc_order_details (
    id integer NOT NULL,
    "order" integer NOT NULL,
    product_id integer NOT NULL,
    variation_id integer NOT NULL,
    quantity integer NOT NULL,
    total_tax real NOT NULL,
    price real NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: wc_orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.wc_orders (
    id integer NOT NULL,
    status character varying(50) NOT NULL,
    currency character varying(3) NOT NULL,
    date_created timestamp(0) with time zone NOT NULL,
    discount_tax real NOT NULL,
    shipping_total real NOT NULL,
    shipping_tax real NOT NULL,
    total_tax real NOT NULL,
    customer_id integer NOT NULL,
    order_key character varying(25) NOT NULL,
    billing_address_1 character varying(255) NOT NULL,
    billing_address_2 character varying(255) NOT NULL,
    billing_city character varying(255) NOT NULL,
    billing_postcode character varying(255) NOT NULL,
    billing_country character varying(255) NOT NULL,
    billing_state character varying(255) NOT NULL,
    billing_phone character varying(255) NOT NULL,
    shipping_address_1 character varying(255) NOT NULL,
    shipping_address_2 character varying(255) NOT NULL,
    shipping_city character varying(255) NOT NULL,
    shipping_postcode character varying(255) NOT NULL,
    shipping_country character varying(255) NOT NULL,
    shipping_state character varying(255) NOT NULL,
    shipping_phone character varying(255) NOT NULL,
    payment_method character varying(50) NOT NULL,
    wc_exists boolean DEFAULT true NOT NULL,
    billing_company character varying(255) NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: wc_product_variations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.wc_product_variations (
    id integer NOT NULL,
    sku character varying(25) NOT NULL,
    price real NOT NULL,
    weight character varying(10) NOT NULL,
    dimensions_length character varying(10) NOT NULL,
    dimensions_width character varying(10) NOT NULL,
    dimensions_height character varying(10) NOT NULL,
    attributes character varying(255)[] NOT NULL,
    wc_exists boolean DEFAULT true NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: wc_products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.wc_products (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    date_created timestamp(0) with time zone NOT NULL,
    description text NOT NULL,
    short_description character varying(255) NOT NULL,
    sku character varying(25) NOT NULL,
    price real NOT NULL,
    weight character varying(10) NOT NULL,
    dimensions_length character varying(10) NOT NULL,
    dimensions_width character varying(10) NOT NULL,
    dimensions_height character varying(10) NOT NULL,
    images character varying(255)[] NOT NULL,
    wc_exists boolean DEFAULT true NOT NULL,
    variations integer[] NOT NULL,
    enterprise integer NOT NULL
);


--
-- Name: webhook_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.webhook_logs (
    id bigint NOT NULL,
    webhook integer NOT NULL,
    enterprise integer NOT NULL,
    url character varying(255) NOT NULL,
    auth_code uuid NOT NULL,
    auth_method character(1) NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    sent text NOT NULL,
    received text NOT NULL,
    received_http_code smallint NOT NULL,
    method character varying(10) NOT NULL
);


--
-- Name: webhook_queue; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.webhook_queue (
    id uuid NOT NULL,
    webhook integer NOT NULL,
    enterprise integer NOT NULL,
    url character varying(255) NOT NULL,
    auth_code uuid NOT NULL,
    auth_method character(1) NOT NULL,
    date_created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    send text NOT NULL,
    method character varying(10) NOT NULL
);


--
-- Name: webhook_settings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.webhook_settings (
    id integer NOT NULL,
    enterprise integer NOT NULL,
    url character varying(255) NOT NULL,
    auth_code uuid NOT NULL,
    auth_method character(1) NOT NULL,
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
    products boolean NOT NULL
);


--
-- Name: account account_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.account
    ADD CONSTRAINT account_pkey PRIMARY KEY (id, enterprise);


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
-- Name: accounting_movement_detail apunte_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounting_movement_detail
    ADD CONSTRAINT apunte_pkey PRIMARY KEY (id);


--
-- Name: accounting_movement asiento_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounting_movement
    ADD CONSTRAINT asiento_pkey PRIMARY KEY (id);


--
-- Name: billing_series billing_series_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.billing_series
    ADD CONSTRAINT billing_series_pkey PRIMARY KEY (id, enterprise);


--
-- Name: carrier carrier_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.carrier
    ADD CONSTRAINT carrier_pkey PRIMARY KEY (id);


--
-- Name: charges chages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.charges
    ADD CONSTRAINT chages_pkey PRIMARY KEY (id);


--
-- Name: state city_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.state
    ADD CONSTRAINT city_pkey PRIMARY KEY (id);


--
-- Name: collection_operation collection_operation_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.collection_operation
    ADD CONSTRAINT collection_operation_pkey PRIMARY KEY (id);


--
-- Name: color color_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.color
    ADD CONSTRAINT color_pkey PRIMARY KEY (id);


--
-- Name: complex_manufacturing_order_manufacturing_order complex_manufacturing_order_manufacturing_order_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_pkey PRIMARY KEY (id);


--
-- Name: complex_manufacturing_order complex_manufacturing_order_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_pkey PRIMARY KEY (id);


--
-- Name: config_accounts_vat config_accounts_vat_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config_accounts_vat
    ADD CONSTRAINT config_accounts_vat_pkey PRIMARY KEY (vat_percent, enterprise);


--
-- Name: config config_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_pkey PRIMARY KEY (id);


--
-- Name: connection_filter connection_filter_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connection_filter
    ADD CONSTRAINT connection_filter_pkey PRIMARY KEY (id);


--
-- Name: connection_filter_user connection_filter_user_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connection_filter_user
    ADD CONSTRAINT connection_filter_user_pkey PRIMARY KEY (connection_filter, "user");


--
-- Name: connection_log connection_log_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connection_log
    ADD CONSTRAINT connection_log_pkey PRIMARY KEY (id);


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
-- Name: custom_fields custom_fields_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_fields
    ADD CONSTRAINT custom_fields_pkey PRIMARY KEY (id);


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
-- Name: email_log email_log_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_log
    ADD CONSTRAINT email_log_pkey PRIMARY KEY (id);


--
-- Name: enterprise_logo enterprise_logo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.enterprise_logo
    ADD CONSTRAINT enterprise_logo_pkey PRIMARY KEY (enterprise);


--
-- Name: group group_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."group"
    ADD CONSTRAINT group_pkey PRIMARY KEY (id);


--
-- Name: hs_codes hs_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hs_codes
    ADD CONSTRAINT hs_codes_pkey PRIMARY KEY (id);


--
-- Name: incoterm incoterm_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.incoterm
    ADD CONSTRAINT incoterm_pkey PRIMARY KEY (id);


--
-- Name: inventory inventory_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory
    ADD CONSTRAINT inventory_pkey PRIMARY KEY (id);


--
-- Name: inventory_products inventory_products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_products
    ADD CONSTRAINT inventory_products_pkey PRIMARY KEY (inventory, product);


--
-- Name: journal journal_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.journal
    ADD CONSTRAINT journal_pkey PRIMARY KEY (id, enterprise);


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
-- Name: logs logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (id);


--
-- Name: manufacturing_order manufacturing_order_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_pkey PRIMARY KEY (id);


--
-- Name: manufacturing_order_type_components manufacturing_order_type_components_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order_type_components
    ADD CONSTRAINT manufacturing_order_type_components_pkey PRIMARY KEY (id);


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
-- Name: pallets pallets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pallets
    ADD CONSTRAINT pallets_pkey PRIMARY KEY (id);


--
-- Name: payment_method payment_method_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_method
    ADD CONSTRAINT payment_method_pkey PRIMARY KEY (id);


--
-- Name: payment_transaction payment_transaction_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_transaction
    ADD CONSTRAINT payment_transaction_pkey PRIMARY KEY (id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: permission_dictionary_group permission_dictionary_group_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permission_dictionary_group
    ADD CONSTRAINT permission_dictionary_group_pkey PRIMARY KEY (enterprise, permission_key, "group");


--
-- Name: permission_dictionary permission_dictionary_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permission_dictionary
    ADD CONSTRAINT permission_dictionary_pkey PRIMARY KEY (enterprise, key);


--
-- Name: product_account poduct_account_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_account
    ADD CONSTRAINT poduct_account_pkey PRIMARY KEY (id);


--
-- Name: pos_terminals pos_terminals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pos_terminals
    ADD CONSTRAINT pos_terminals_pkey PRIMARY KEY (id);


--
-- Name: product_family product_family_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_family
    ADD CONSTRAINT product_family_pkey PRIMARY KEY (id);


--
-- Name: product_image product_image_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_image
    ADD CONSTRAINT product_image_pkey PRIMARY KEY (id);


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
-- Name: ps_address ps_address_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_address
    ADD CONSTRAINT ps_address_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_carrier ps_carrier_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_carrier
    ADD CONSTRAINT ps_carrier_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_country ps_country_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_country
    ADD CONSTRAINT ps_country_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_currency ps_currency_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_currency
    ADD CONSTRAINT ps_currency_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_customer ps_customer_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_customer
    ADD CONSTRAINT ps_customer_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_language ps_language_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_language
    ADD CONSTRAINT ps_language_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_order_detail ps_order_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_order_detail
    ADD CONSTRAINT ps_order_detail_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_order ps_order_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_order
    ADD CONSTRAINT ps_order_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_product_combination ps_product_combination_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_product_combination
    ADD CONSTRAINT ps_product_combination_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_product_option_values ps_product_option_values_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_product_option_values
    ADD CONSTRAINT ps_product_option_values_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_product ps_product_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_product
    ADD CONSTRAINT ps_product_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_state ps_state_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_state
    ADD CONSTRAINT ps_state_pkey PRIMARY KEY (id, enterprise);


--
-- Name: ps_zone ps_zone_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_zone
    ADD CONSTRAINT ps_zone_pkey PRIMARY KEY (id, enterprise);


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
-- Name: pwd_blacklist pwd_blacklist_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pwd_blacklist
    ADD CONSTRAINT pwd_blacklist_pkey PRIMARY KEY (pwd);


--
-- Name: pwd_sha1_blacklist pwd_sha1_blacklist_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pwd_sha1_blacklist
    ADD CONSTRAINT pwd_sha1_blacklist_pkey PRIMARY KEY (hash);


--
-- Name: report_template report_template_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.report_template
    ADD CONSTRAINT report_template_pkey PRIMARY KEY (enterprise, key);


--
-- Name: report_template_translation report_template_translation_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.report_template_translation
    ADD CONSTRAINT report_template_translation_pkey PRIMARY KEY (enterprise, key, language);


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
-- Name: sales_order_detail_digital_product_data sales_order_detail_digital_product_data_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail_digital_product_data
    ADD CONSTRAINT sales_order_detail_digital_product_data_pkey PRIMARY KEY (id);


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
-- Name: shipping_status_history shipping_status_history_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping_status_history
    ADD CONSTRAINT shipping_status_history_pkey PRIMARY KEY (shipping);


--
-- Name: shipping_tag shipping_tag_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping_tag
    ADD CONSTRAINT shipping_tag_pkey PRIMARY KEY (id);


--
-- Name: stock stock_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_pkey PRIMARY KEY (product, warehouse, enterprise);


--
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY (id);


--
-- Name: sy_addresses sy_addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_addresses
    ADD CONSTRAINT sy_addresses_pkey PRIMARY KEY (id, enterprise);


--
-- Name: sy_customers sy_customers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_customers
    ADD CONSTRAINT sy_customers_pkey PRIMARY KEY (id, enterprise);


--
-- Name: sy_draft_order_line_item sy_draft_order_line_item_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_draft_order_line_item
    ADD CONSTRAINT sy_draft_order_line_item_pkey PRIMARY KEY (id, enterprise);


--
-- Name: sy_draft_orders sy_draft_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_draft_orders
    ADD CONSTRAINT sy_draft_orders_pkey PRIMARY KEY (id, enterprise);


--
-- Name: sy_order_line_item sy_order_line_item_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_order_line_item
    ADD CONSTRAINT sy_order_line_item_pkey PRIMARY KEY (id, enterprise);


--
-- Name: sy_orders sy_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_orders
    ADD CONSTRAINT sy_orders_pkey PRIMARY KEY (id, enterprise);


--
-- Name: sy_products sy_products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_products
    ADD CONSTRAINT sy_products_pkey PRIMARY KEY (id, enterprise);


--
-- Name: sy_variants sy_variants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_variants
    ADD CONSTRAINT sy_variants_pkey PRIMARY KEY (id, enterprise);


--
-- Name: transactional_log transactional_log_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactional_log
    ADD CONSTRAINT transactional_log_pkey PRIMARY KEY (id);


--
-- Name: transfer_between_warehouses_detail transfer_between_warehouses_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfer_between_warehouses_detail
    ADD CONSTRAINT transfer_between_warehouses_detail_pkey PRIMARY KEY (id);


--
-- Name: transfer_between_warehouses transfer_between_warehouses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfer_between_warehouses
    ADD CONSTRAINT transfer_between_warehouses_pkey PRIMARY KEY (id);


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
    ADD CONSTRAINT warehouse_pkey PRIMARY KEY (id, enterprise);


--
-- Name: wc_customers wc_customers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wc_customers
    ADD CONSTRAINT wc_customers_pkey PRIMARY KEY (id, enterprise);


--
-- Name: wc_order_details wc_order_details_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wc_order_details
    ADD CONSTRAINT wc_order_details_pkey PRIMARY KEY (id, enterprise);


--
-- Name: wc_orders wc_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wc_orders
    ADD CONSTRAINT wc_orders_pkey PRIMARY KEY (id, enterprise);


--
-- Name: wc_product_variations wc_product_variations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wc_product_variations
    ADD CONSTRAINT wc_product_variations_pkey PRIMARY KEY (id, enterprise);


--
-- Name: wc_products wc_products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wc_products
    ADD CONSTRAINT wc_products_pkey PRIMARY KEY (id, enterprise);


--
-- Name: webhook_logs webhook_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_logs
    ADD CONSTRAINT webhook_logs_pkey PRIMARY KEY (id);


--
-- Name: webhook_queue webhook_queue_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_queue
    ADD CONSTRAINT webhook_queue_pkey PRIMARY KEY (id);


--
-- Name: webhook_settings webhook_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_settings
    ADD CONSTRAINT webhook_settings_pkey PRIMARY KEY (id);


--
-- Name: account_account_number_journal; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX account_account_number_journal ON public.account USING btree (enterprise, account_number, journal);


--
-- Name: account_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX account_id_enterprise ON public.account USING btree (id, enterprise);


--
-- Name: account_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX account_name ON public.account USING gin (name public.gin_trgm_ops);


--
-- Name: accounting_movement_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX accounting_movement_date_created ON public.accounting_movement USING btree (date_created DESC NULLS LAST);


--
-- Name: accounting_movement_detail_document_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX accounting_movement_detail_document_name ON public.accounting_movement_detail USING hash (document_name);


--
-- Name: accounting_movement_detail_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX accounting_movement_detail_id_enterprise ON public.accounting_movement_detail USING btree (id, enterprise);


--
-- Name: accounting_movement_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX accounting_movement_id_enterprise ON public.accounting_movement USING btree (id, enterprise);


--
-- Name: address_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX address_id_enterprise ON public.address USING btree (id, enterprise);


--
-- Name: address_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX address_name ON public.address USING gin (address public.gin_trgm_ops);


--
-- Name: address_ps_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX address_ps_id ON public.address USING btree (enterprise, ps_id) WHERE (ps_id <> 0);


--
-- Name: address_sy_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX address_sy_id ON public.address USING btree (enterprise, sy_id) WHERE (sy_id <> 0);


--
-- Name: api_key_basic_auth; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX api_key_basic_auth ON public.api_key USING btree (basic_auth_user, basic_auth_password) WHERE (auth = 'B'::bpchar);


--
-- Name: api_key_token; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX api_key_token ON public.api_key USING btree (token) WHERE (token IS NOT NULL);


--
-- Name: billing_series_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX billing_series_id_enterprise ON public.billing_series USING btree (id, enterprise);


--
-- Name: carrier_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX carrier_id_enterprise ON public.carrier USING btree (id, enterprise);


--
-- Name: carrier_ps_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX carrier_ps_id ON public.carrier USING btree (ps_id) WHERE (ps_id <> 0);


--
-- Name: collection_operation_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX collection_operation_date_created ON public.collection_operation USING btree (date_created);


--
-- Name: collection_operation_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX collection_operation_id_enterprise ON public.collection_operation USING btree (id, enterprise);


--
-- Name: collection_operation_status_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX collection_operation_status_enterprise ON public.collection_operation USING btree (status, enterprise);


--
-- Name: color_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX color_id_enterprise ON public.color USING btree (id, enterprise);


--
-- Name: complex_manufacturing_order_complex_manufacturing_order_manufac; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX complex_manufacturing_order_complex_manufacturing_order_manufac ON public.complex_manufacturing_order_manufacturing_order USING btree (complex_manufacturing_order, manufacturing_order);


--
-- Name: complex_manufacturing_order_complex_manufacturing_order_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX complex_manufacturing_order_complex_manufacturing_order_type ON public.complex_manufacturing_order_manufacturing_order USING btree (complex_manufacturing_order, type);


--
-- Name: complex_manufacturing_order_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX complex_manufacturing_order_id_enterprise ON public.complex_manufacturing_order USING btree (id, enterprise);


--
-- Name: complex_manufacturing_order_manufacturing_order_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX complex_manufacturing_order_manufacturing_order_id_enterprise ON public.complex_manufacturing_order_manufacturing_order USING btree (id, enterprise);


--
-- Name: complex_manufacturing_order_manufacturing_order_pending_product; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX complex_manufacturing_order_manufacturing_order_pending_product ON public.complex_manufacturing_order_manufacturing_order USING btree (product, manufactured, type, sale_order_detail) WHERE ((NOT manufactured) AND (type = 'O'::bpchar) AND (sale_order_detail IS NULL));


--
-- Name: config_enterprise_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX config_enterprise_key ON public.config USING btree (enterprise_key);


--
-- Name: connection_log_date_connected; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX connection_log_date_connected ON public.connection_log USING btree (date_connected DESC NULLS LAST);


--
-- Name: connection_log_user_date_connected; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX connection_log_user_date_connected ON public.connection_log USING btree ("user", date_connected DESC NULLS LAST) WHERE (date_disconnected IS NULL);


--
-- Name: country_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX country_id_enterprise ON public.country USING btree (id, enterprise);


--
-- Name: country_iso_2; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX country_iso_2 ON public.country USING btree (enterprise, iso_2);


--
-- Name: country_iso_3; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX country_iso_3 ON public.country USING btree (enterprise, iso_3) WHERE (iso_3 <> ''::bpchar);


--
-- Name: country_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX country_name ON public.country USING gin (name public.gin_trgm_ops);


--
-- Name: currency_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX currency_id_enterprise ON public.currency USING btree (id, enterprise);


--
-- Name: currency_iso_code; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX currency_iso_code ON public.currency USING btree (enterprise, iso_code);


--
-- Name: currency_iso_num; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX currency_iso_num ON public.currency USING btree (enterprise, iso_num);


--
-- Name: customer_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX customer_email ON public.customer USING gin (email public.gin_trgm_ops);


--
-- Name: customer_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX customer_id_enterprise ON public.customer USING btree (id, enterprise);


--
-- Name: customer_name_trgm; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX customer_name_trgm ON public.customer USING gin (name public.gin_trgm_ops);


--
-- Name: customer_ps_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX customer_ps_id ON public.customer USING btree (enterprise, ps_id) WHERE (ps_id <> 0);


--
-- Name: customer_sy_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX customer_sy_id ON public.customer USING btree (enterprise, sy_id) WHERE (sy_id <> 0);


--
-- Name: customer_tax_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX customer_tax_id ON public.customer USING gin (tax_id public.gin_trgm_ops);


--
-- Name: customer_wc_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX customer_wc_id ON public.customer USING btree (enterprise, wc_id) WHERE (wc_id <> 0);


--
-- Name: document_container_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX document_container_id_enterprise ON public.document_container USING btree (id, enterprise);


--
-- Name: document_uuid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX document_uuid ON public.document USING btree (uuid);


--
-- Name: email_log_date_sent; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX email_log_date_sent ON public.email_log USING btree (date_sent DESC NULLS LAST);


--
-- Name: email_log_trgm; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX email_log_trgm ON public.email_log USING gin (email_from public.gin_trgm_ops, name_from public.gin_trgm_ops, destination_email public.gin_trgm_ops, destination_name public.gin_trgm_ops, subject public.gin_trgm_ops);


--
-- Name: group_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX group_id_enterprise ON public."group" USING btree (id, enterprise);


--
-- Name: incoterm_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX incoterm_id_enterprise ON public.incoterm USING btree (id, enterprise);


--
-- Name: incoterm_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX incoterm_key ON public.incoterm USING btree (enterprise, key);


--
-- Name: journal_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX journal_id_enterprise ON public.journal USING btree (id, enterprise);


--
-- Name: language_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX language_id_enterprise ON public.language USING btree (id, enterprise);


--
-- Name: language_iso_2; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX language_iso_2 ON public.language USING btree (enterprise, iso_2);


--
-- Name: language_iso_3; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX language_iso_3 ON public.language USING btree (enterprise, iso_3) WHERE (iso_3 <> ''::bpchar);


--
-- Name: language_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX language_name ON public.language USING gin (name public.gin_trgm_ops);


--
-- Name: login_tokens_name_ip_address; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX login_tokens_name_ip_address ON public.login_tokens USING btree (name, ip_address);


--
-- Name: manufacturing_order_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX manufacturing_order_date_created ON public.manufacturing_order USING btree (date_created DESC NULLS LAST);


--
-- Name: manufacturing_order_for_stock_pending; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX manufacturing_order_for_stock_pending ON public.manufacturing_order USING btree (enterprise, product, manufactured, order_detail, complex) WHERE ((NOT manufactured) AND (order_detail IS NULL) AND (NOT complex));


--
-- Name: manufacturing_order_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX manufacturing_order_id_enterprise ON public.manufacturing_order USING btree (id, enterprise);


--
-- Name: manufacturing_order_type_components_component; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX manufacturing_order_type_components_component ON public.manufacturing_order_type_components USING btree (manufacturing_order_type, product);


--
-- Name: manufacturing_order_type_components_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX manufacturing_order_type_components_id_enterprise ON public.manufacturing_order_type_components USING btree (id, enterprise);


--
-- Name: manufacturing_order_type_components_manufacturing_order_type_ty; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX manufacturing_order_type_components_manufacturing_order_type_ty ON public.manufacturing_order_type_components USING btree (manufacturing_order_type, type, product);


--
-- Name: manufacturing_order_type_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX manufacturing_order_type_id_enterprise ON public.manufacturing_order_type USING btree (id, enterprise);


--
-- Name: manufacturing_order_uuid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX manufacturing_order_uuid ON public.manufacturing_order USING btree (uuid);


--
-- Name: packaging_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX packaging_id_enterprise ON public.packaging USING btree (id, enterprise);


--
-- Name: payment_method_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX payment_method_id_enterprise ON public.payment_method USING btree (id, enterprise);


--
-- Name: payment_transaction_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX payment_transaction_id_enterprise ON public.payment_transaction USING btree (id, enterprise);


--
-- Name: pos_terminals_uuid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX pos_terminals_uuid ON public.pos_terminals USING btree (uuid);


--
-- Name: product_account_product_type; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX product_account_product_type ON public.product_account USING btree (product, type);


--
-- Name: product_barcode; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX product_barcode ON public.product USING btree (enterprise, barcode) WHERE (barcode <> ''::bpchar);


--
-- Name: product_family_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX product_family_id_enterprise ON public.product_family USING btree (id, enterprise);


--
-- Name: product_family_reference; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX product_family_reference ON public.product_family USING btree (enterprise, reference);


--
-- Name: product_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX product_id_enterprise ON public.product USING btree (id, enterprise);


--
-- Name: product_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX product_name ON public.product USING gin (name public.gin_trgm_ops);


--
-- Name: product_ps_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX product_ps_id ON public.product USING btree (enterprise, ps_id, ps_combination_id) WHERE (ps_id <> 0);


--
-- Name: product_reference; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX product_reference ON public.product USING gin (reference public.gin_trgm_ops) WHERE ((reference)::text <> ''::text);


--
-- Name: product_sy_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX product_sy_id ON public.product USING btree (enterprise, sy_id, sy_variant_id) WHERE (sy_id <> 0);


--
-- Name: product_track_minimum_stock; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX product_track_minimum_stock ON public.product USING btree (track_minimum_stock) WHERE (track_minimum_stock = true);


--
-- Name: products_wc_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX products_wc_id ON public.product USING btree (enterprise, wc_id, wc_variation_id) WHERE (wc_id <> 0);


--
-- Name: ps_address_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_address_ps_exists ON public.ps_address USING btree (ps_exists);


--
-- Name: ps_carrier_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_carrier_ps_exists ON public.ps_carrier USING btree (ps_exists);


--
-- Name: ps_country_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_country_ps_exists ON public.ps_country USING btree (ps_exists);


--
-- Name: ps_currency_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_currency_ps_exists ON public.ps_currency USING btree (ps_exists);


--
-- Name: ps_language_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_language_ps_exists ON public.ps_language USING btree (ps_exists);


--
-- Name: ps_order_detail_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_order_detail_ps_exists ON public.ps_order_detail USING btree (ps_exists);


--
-- Name: ps_order_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_order_ps_exists ON public.ps_order USING btree (ps_exists);


--
-- Name: ps_product_combination_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_product_combination_ps_exists ON public.ps_product_combination USING btree (ps_exists);


--
-- Name: ps_product_option_values_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_product_option_values_ps_exists ON public.ps_product_option_values USING btree (ps_exists);


--
-- Name: ps_product_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_product_ps_exists ON public.ps_product USING btree (ps_exists);


--
-- Name: ps_state_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_state_ps_exists ON public.ps_state USING btree (ps_exists);


--
-- Name: ps_zone_ps_exists; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ps_zone_ps_exists ON public.ps_zone USING btree (ps_exists);


--
-- Name: purcahse_invoice_invoice_number; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX purcahse_invoice_invoice_number ON public.purchase_invoice USING btree (enterprise, billing_series, invoice_number DESC NULLS LAST);


--
-- Name: purchase_delivery_note_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX purchase_delivery_note_date_created ON public.purchase_delivery_note USING btree (date_created DESC NULLS LAST);


--
-- Name: purchase_delivery_note_delivery_note_number; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX purchase_delivery_note_delivery_note_number ON public.purchase_delivery_note USING btree (enterprise, billing_series, delivery_note_number DESC NULLS LAST);


--
-- Name: purchase_delivery_note_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX purchase_delivery_note_id_enterprise ON public.purchase_delivery_note USING btree (id, enterprise);


--
-- Name: purchase_invoice_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX purchase_invoice_date_created ON public.purchase_invoice USING btree (date_created DESC NULLS LAST);


--
-- Name: purchase_invoice_details_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX purchase_invoice_details_id_enterprise ON public.purchase_invoice_details USING btree (id, enterprise);


--
-- Name: purchase_invoice_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX purchase_invoice_id_enterprise ON public.purchase_invoice USING btree (id, enterprise);


--
-- Name: purchase_order_detail_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX purchase_order_detail_id_enterprise ON public.purchase_order_detail USING btree (id, enterprise);


--
-- Name: purchase_order_detail_purchase_order_product; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX purchase_order_detail_purchase_order_product ON public.purchase_order_detail USING btree ("order", product);


--
-- Name: purchase_order_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX purchase_order_id_enterprise ON public.purchase_order USING btree (id, enterprise);


--
-- Name: purchase_order_order_number; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX purchase_order_order_number ON public.purchase_order USING btree (enterprise, billing_series, order_number DESC NULLS LAST);


--
-- Name: purchase_order_supplier_reference; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX purchase_order_supplier_reference ON public.purchase_order USING gin (supplier_reference public.gin_trgm_ops);


--
-- Name: sales_delivery_note_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX sales_delivery_note_date_created ON public.sales_delivery_note USING btree (date_created DESC NULLS LAST);


--
-- Name: sales_delivery_note_delivery_note_number; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_delivery_note_delivery_note_number ON public.sales_delivery_note USING btree (enterprise, billing_series, delivery_note_number DESC NULLS LAST);


--
-- Name: sales_delivery_note_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_delivery_note_id_enterprise ON public.sales_delivery_note USING btree (id, enterprise);


--
-- Name: sales_invoice_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX sales_invoice_date_created ON public.sales_invoice USING btree (date_created DESC NULLS LAST);


--
-- Name: sales_invoice_detail_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_invoice_detail_id_enterprise ON public.sales_invoice_detail USING btree (id, enterprise);


--
-- Name: sales_invoice_detail_invoice_product; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_invoice_detail_invoice_product ON public.sales_invoice_detail USING btree (invoice, product);


--
-- Name: sales_invoice_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_invoice_id_enterprise ON public.sales_invoice USING btree (id, enterprise);


--
-- Name: sales_invoice_invoice_number; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_invoice_invoice_number ON public.sales_invoice USING btree (enterprise, billing_series, invoice_number DESC NULLS LAST);


--
-- Name: sales_order_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX sales_order_date_created ON public.sales_order USING btree (date_created DESC NULLS LAST);


--
-- Name: sales_order_detail_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_detail_id_enterprise ON public.sales_order_detail USING btree (id, enterprise);


--
-- Name: sales_order_detail_ps_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_detail_ps_id ON public.sales_order_detail USING btree (enterprise, ps_id) WHERE (ps_id <> 0);


--
-- Name: sales_order_detail_sales_order_product; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_detail_sales_order_product ON public.sales_order_detail USING btree ("order", product);


--
-- Name: sales_order_detail_sy_draft_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_detail_sy_draft_id ON public.sales_order_detail USING btree (enterprise, sy_draft_id) WHERE (sy_draft_id <> 0);


--
-- Name: sales_order_detail_sy_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_detail_sy_id ON public.sales_order_detail USING btree (enterprise, sy_id) WHERE (sy_id <> 0);


--
-- Name: sales_order_detail_wc_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_detail_wc_id ON public.sales_order_detail USING btree (enterprise, wc_id) WHERE (wc_id <> 0);


--
-- Name: sales_order_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_id_enterprise ON public.sales_order USING btree (id, enterprise);


--
-- Name: sales_order_order_number; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_order_number ON public.sales_order USING btree (enterprise, billing_series, order_number DESC NULLS LAST);


--
-- Name: sales_order_ps_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_ps_id ON public.sales_order USING btree (enterprise, ps_id) WHERE (ps_id <> 0);


--
-- Name: sales_order_reference; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX sales_order_reference ON public.sales_order USING gin (reference public.gin_trgm_ops);


--
-- Name: sales_order_sy_draft_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_sy_draft_id ON public.sales_order USING btree (enterprise, sy_draft_id) WHERE (sy_draft_id <> 0);


--
-- Name: sales_order_sy_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_sy_id ON public.sales_order USING btree (enterprise, sy_id) WHERE (sy_id <> 0);


--
-- Name: sales_order_wc_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX sales_order_wc_id ON public.sales_order USING btree (enterprise, wc_id) WHERE (wc_id <> 0);


--
-- Name: shipping_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX shipping_id_enterprise ON public.shipping USING btree (id, enterprise);


--
-- Name: shipping_sent_collected; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX shipping_sent_collected ON public.shipping USING btree (sent, collected);


--
-- Name: shipping_sent_collected_delivered; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX shipping_sent_collected_delivered ON public.shipping USING btree (sent, collected, delivered);


--
-- Name: state_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX state_id_enterprise ON public.state USING btree (id, enterprise);


--
-- Name: state_iso_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX state_iso_code ON public.state USING btree (iso_code) WHERE ((iso_code)::text <> ''::text);


--
-- Name: state_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX state_name ON public.state USING gin (name public.gin_trgm_ops);


--
-- Name: supplier_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX supplier_email ON public.suppliers USING gin (email public.gin_trgm_ops);


--
-- Name: supplier_name_trgm; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX supplier_name_trgm ON public.suppliers USING gin (name public.gin_trgm_ops);


--
-- Name: supplier_tax_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX supplier_tax_id ON public.suppliers USING gin (tax_id public.gin_trgm_ops);


--
-- Name: suppliers_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX suppliers_id_enterprise ON public.suppliers USING btree (id, enterprise);


--
-- Name: sy_draft_orders_order_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX sy_draft_orders_order_id ON public.sy_draft_orders USING btree (order_id) WHERE (order_id IS NOT NULL);


--
-- Name: transactional_log_enterprise_table_register_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX transactional_log_enterprise_table_register_id ON public.transactional_log USING btree (enterprise, "table", register_id);


--
-- Name: transfer_between_warehouses_detail_barcode; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX transfer_between_warehouses_detail_barcode ON public.transfer_between_warehouses_detail USING btree (enterprise, transfer_between_warehouses, product) WHERE (quantity_transferred < quantity);


--
-- Name: transfer_between_warehouses_enterprise_finished_date_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX transfer_between_warehouses_enterprise_finished_date_created ON public.transfer_between_warehouses USING btree (enterprise, finished, date_created);


--
-- Name: transfer_between_warehouses_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX transfer_between_warehouses_id_enterprise ON public.transfer_between_warehouses USING btree (id, enterprise);


--
-- Name: user_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX user_id_enterprise ON public."user" USING btree (id, config);


--
-- Name: user_username; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX user_username ON public."user" USING btree (config, username);


--
-- Name: webhook_settings_id_enterprise; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX webhook_settings_id_enterprise ON public.webhook_settings USING btree (id, enterprise);


--
-- Name: account set_account_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_account_id BEFORE INSERT ON public.account FOR EACH ROW EXECUTE FUNCTION public.set_account_id();


--
-- Name: accounting_movement_detail set_accounting_movement_detail_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_accounting_movement_detail_id BEFORE INSERT ON public.accounting_movement_detail FOR EACH ROW EXECUTE FUNCTION public.set_accounting_movement_detail_id();


--
-- Name: accounting_movement set_accounting_movement_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_accounting_movement_id BEFORE INSERT ON public.accounting_movement FOR EACH ROW EXECUTE FUNCTION public.set_accounting_movement_id();


--
-- Name: address set_address_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_address_id BEFORE INSERT ON public.address FOR EACH ROW EXECUTE FUNCTION public.set_address_id();


--
-- Name: api_key set_api_key_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_api_key_id BEFORE INSERT ON public.api_key FOR EACH ROW EXECUTE FUNCTION public.set_api_key_id();


--
-- Name: carrier set_carrier_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_carrier_id BEFORE INSERT ON public.carrier FOR EACH ROW EXECUTE FUNCTION public.set_carrier_id();


--
-- Name: charges set_charges_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_charges_id BEFORE INSERT ON public.charges FOR EACH ROW EXECUTE FUNCTION public.set_charges_id();


--
-- Name: collection_operation set_collection_operation_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_collection_operation_id BEFORE INSERT ON public.collection_operation FOR EACH ROW EXECUTE FUNCTION public.set_collection_operation_id();


--
-- Name: color set_color_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_color_id BEFORE INSERT ON public.color FOR EACH ROW EXECUTE FUNCTION public.set_color_id();


--
-- Name: complex_manufacturing_order set_complex_manufacturing_order_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_complex_manufacturing_order_id BEFORE INSERT ON public.complex_manufacturing_order FOR EACH ROW EXECUTE FUNCTION public.set_complex_manufacturing_order_id();


--
-- Name: complex_manufacturing_order_manufacturing_order set_complex_manufacturing_order_manufacturing_order_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_complex_manufacturing_order_manufacturing_order_id BEFORE INSERT ON public.complex_manufacturing_order_manufacturing_order FOR EACH ROW EXECUTE FUNCTION public.set_complex_manufacturing_order_manufacturing_order_id();


--
-- Name: config set_config_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_config_id BEFORE INSERT ON public.config FOR EACH ROW EXECUTE FUNCTION public.set_config_id();


--
-- Name: connection_filter set_connection_filter_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_connection_filter_id BEFORE INSERT ON public.connection_filter FOR EACH ROW EXECUTE FUNCTION public.set_connection_filter_id();


--
-- Name: connection_log set_connection_log_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_connection_log_id BEFORE INSERT ON public.connection_log FOR EACH ROW EXECUTE FUNCTION public.set_connection_log_id();


--
-- Name: country set_country_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_country_id BEFORE INSERT ON public.country FOR EACH ROW EXECUTE FUNCTION public.set_country_id();


--
-- Name: currency set_currency_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_currency_id BEFORE INSERT ON public.currency FOR EACH ROW EXECUTE FUNCTION public.set_currency_id();


--
-- Name: custom_fields set_custom_fields_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_custom_fields_id BEFORE INSERT ON public.custom_fields FOR EACH ROW EXECUTE FUNCTION public.set_custom_fields_id();


--
-- Name: customer set_customer_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_customer_id BEFORE INSERT ON public.customer FOR EACH ROW EXECUTE FUNCTION public.set_customer_id();


--
-- Name: document_container set_document_container_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_document_container_id BEFORE INSERT ON public.document_container FOR EACH ROW EXECUTE FUNCTION public.set_document_container_id();


--
-- Name: document set_document_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_document_id BEFORE INSERT ON public.document FOR EACH ROW EXECUTE FUNCTION public.set_document_id();


--
-- Name: email_log set_email_log_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_email_log_id BEFORE INSERT ON public.email_log FOR EACH ROW EXECUTE FUNCTION public.set_email_log_id();


--
-- Name: group set_group_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_group_id BEFORE INSERT ON public."group" FOR EACH ROW EXECUTE FUNCTION public.set_group_id();


--
-- Name: incoterm set_incoterm_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_incoterm_id BEFORE INSERT ON public.incoterm FOR EACH ROW EXECUTE FUNCTION public.set_incoterm_id();


--
-- Name: inventory set_inventory_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_inventory_id BEFORE INSERT ON public.inventory FOR EACH ROW EXECUTE FUNCTION public.set_inventory_id();


--
-- Name: language set_language_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_language_id BEFORE INSERT ON public.language FOR EACH ROW EXECUTE FUNCTION public.set_language_id();


--
-- Name: login_tokens set_login_tokens_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_login_tokens_id BEFORE INSERT ON public.login_tokens FOR EACH ROW EXECUTE FUNCTION public.set_login_tokens_id();


--
-- Name: logs set_logs_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_logs_id BEFORE INSERT ON public.logs FOR EACH ROW EXECUTE FUNCTION public.set_logs_id();


--
-- Name: manufacturing_order set_manufacturing_order_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_manufacturing_order_id BEFORE INSERT ON public.manufacturing_order FOR EACH ROW EXECUTE FUNCTION public.set_manufacturing_order_id();


--
-- Name: manufacturing_order_type_components set_manufacturing_order_type_components_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_manufacturing_order_type_components_id BEFORE INSERT ON public.manufacturing_order_type_components FOR EACH ROW EXECUTE FUNCTION public.set_manufacturing_order_type_components_id();


--
-- Name: manufacturing_order_type set_manufacturing_order_type_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_manufacturing_order_type_id BEFORE INSERT ON public.manufacturing_order_type FOR EACH ROW EXECUTE FUNCTION public.set_manufacturing_order_type_id();


--
-- Name: packages set_packages_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_packages_id BEFORE INSERT ON public.packages FOR EACH ROW EXECUTE FUNCTION public.set_packages_id();


--
-- Name: packaging set_packaging_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_packaging_id BEFORE INSERT ON public.packaging FOR EACH ROW EXECUTE FUNCTION public.set_packaging_id();


--
-- Name: pallets set_pallets_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_pallets_id BEFORE INSERT ON public.pallets FOR EACH ROW EXECUTE FUNCTION public.set_pallets_id();


--
-- Name: payment_method set_payment_method_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_payment_method_id BEFORE INSERT ON public.payment_method FOR EACH ROW EXECUTE FUNCTION public.set_payment_method_id();


--
-- Name: payment_transaction set_payment_transaction_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_payment_transaction_id BEFORE INSERT ON public.payment_transaction FOR EACH ROW EXECUTE FUNCTION public.set_payment_transaction_id();


--
-- Name: payments set_payments_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_payments_id BEFORE INSERT ON public.payments FOR EACH ROW EXECUTE FUNCTION public.set_payments_id();


--
-- Name: pos_terminals set_pos_terminals_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_pos_terminals_id BEFORE INSERT ON public.pos_terminals FOR EACH ROW EXECUTE FUNCTION public.set_pos_terminals_id();


--
-- Name: product_account set_product_account_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_product_account_id BEFORE INSERT ON public.product_account FOR EACH ROW EXECUTE FUNCTION public.set_product_account_id();


--
-- Name: product_family set_product_family_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_product_family_id BEFORE INSERT ON public.product_family FOR EACH ROW EXECUTE FUNCTION public.set_product_family_id();


--
-- Name: product set_product_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_product_id BEFORE INSERT ON public.product FOR EACH ROW EXECUTE FUNCTION public.set_product_id();


--
-- Name: product_image set_product_image_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_product_image_id BEFORE INSERT ON public.product_image FOR EACH ROW EXECUTE FUNCTION public.set_product_image_id();


--
-- Name: purchase_delivery_note set_purchase_delivery_note_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_purchase_delivery_note_id BEFORE INSERT ON public.purchase_delivery_note FOR EACH ROW EXECUTE FUNCTION public.set_purchase_delivery_note_id();


--
-- Name: purchase_invoice_details set_purchase_invoice_details_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_purchase_invoice_details_id BEFORE INSERT ON public.purchase_invoice_details FOR EACH ROW EXECUTE FUNCTION public.set_purchase_invoice_details_id();


--
-- Name: purchase_invoice set_purchase_invoice_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_purchase_invoice_id BEFORE INSERT ON public.purchase_invoice FOR EACH ROW EXECUTE FUNCTION public.set_purchase_invoice_id();


--
-- Name: purchase_order_detail set_purchase_order_detail_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_purchase_order_detail_id BEFORE INSERT ON public.purchase_order_detail FOR EACH ROW EXECUTE FUNCTION public.set_purchase_order_detail_id();


--
-- Name: purchase_order set_purchase_order_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_purchase_order_id BEFORE INSERT ON public.purchase_order FOR EACH ROW EXECUTE FUNCTION public.set_purchase_order_id();


--
-- Name: sales_delivery_note set_sales_delivery_note_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_sales_delivery_note_id BEFORE INSERT ON public.sales_delivery_note FOR EACH ROW EXECUTE FUNCTION public.set_sales_delivery_note_id();


--
-- Name: sales_invoice_detail set_sales_invoice_detail_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_sales_invoice_detail_id BEFORE INSERT ON public.sales_invoice_detail FOR EACH ROW EXECUTE FUNCTION public.set_sales_invoice_detail_id();


--
-- Name: sales_invoice set_sales_invoice_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_sales_invoice_id BEFORE INSERT ON public.sales_invoice FOR EACH ROW EXECUTE FUNCTION public.set_sales_invoice_id();


--
-- Name: sales_order_detail_digital_product_data set_sales_order_detail_digital_product_data_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_sales_order_detail_digital_product_data_id BEFORE INSERT ON public.sales_order_detail_digital_product_data FOR EACH ROW EXECUTE FUNCTION public.set_sales_order_detail_digital_product_data_id();


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
-- Name: shipping set_shipping_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_shipping_id BEFORE INSERT ON public.shipping FOR EACH ROW EXECUTE FUNCTION public.set_shipping_id();


--
-- Name: shipping_status_history set_shipping_status_history_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_shipping_status_history_id BEFORE INSERT ON public.shipping_status_history FOR EACH ROW EXECUTE FUNCTION public.set_shipping_status_history_id();


--
-- Name: shipping_tag set_shipping_tag_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_shipping_tag_id BEFORE INSERT ON public.shipping_tag FOR EACH ROW EXECUTE FUNCTION public.set_shipping_tag_id();


--
-- Name: state set_state_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_state_id BEFORE INSERT ON public.state FOR EACH ROW EXECUTE FUNCTION public.set_state_id();


--
-- Name: suppliers set_suppliers_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_suppliers_id BEFORE INSERT ON public.suppliers FOR EACH ROW EXECUTE FUNCTION public.set_suppliers_id();


--
-- Name: transactional_log set_transactional_log_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_transactional_log_id BEFORE INSERT ON public.transactional_log FOR EACH ROW EXECUTE FUNCTION public.set_transactional_log_id();


--
-- Name: transfer_between_warehouses_detail set_transfer_between_warehouses_detail_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_transfer_between_warehouses_detail_id BEFORE INSERT ON public.transfer_between_warehouses_detail FOR EACH ROW EXECUTE FUNCTION public.set_transfer_between_warehouses_detail_id();


--
-- Name: transfer_between_warehouses set_transfer_between_warehouses_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_transfer_between_warehouses_id BEFORE INSERT ON public.transfer_between_warehouses FOR EACH ROW EXECUTE FUNCTION public.set_transfer_between_warehouses_id();


--
-- Name: user set_user_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_user_id BEFORE INSERT ON public."user" FOR EACH ROW EXECUTE FUNCTION public.set_user_id();


--
-- Name: warehouse_movement set_warehouse_movement_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_warehouse_movement_id BEFORE INSERT ON public.warehouse_movement FOR EACH ROW EXECUTE FUNCTION public.set_warehouse_movement_id();


--
-- Name: webhook_logs set_webhook_logs_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_webhook_logs_id BEFORE INSERT ON public.webhook_logs FOR EACH ROW EXECUTE FUNCTION public.set_webhook_logs_id();


--
-- Name: webhook_settings set_webhook_settings_id; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_webhook_settings_id BEFORE INSERT ON public.webhook_settings FOR EACH ROW EXECUTE FUNCTION public.set_webhook_settings_id();


--
-- Name: account account_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.account
    ADD CONSTRAINT account_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: account account_journal; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.account
    ADD CONSTRAINT account_journal FOREIGN KEY (journal, enterprise) REFERENCES public.journal(id, enterprise);


--
-- Name: accounting_movement accounting_movement_billing_serie; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounting_movement
    ADD CONSTRAINT accounting_movement_billing_serie FOREIGN KEY (billing_serie, enterprise) REFERENCES public.billing_series(id, enterprise);


--
-- Name: accounting_movement_detail accounting_movement_detail_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounting_movement_detail
    ADD CONSTRAINT accounting_movement_detail_account FOREIGN KEY (account, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: accounting_movement_detail accounting_movement_detail_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounting_movement_detail
    ADD CONSTRAINT accounting_movement_detail_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: accounting_movement_detail accounting_movement_detail_journal; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounting_movement_detail
    ADD CONSTRAINT accounting_movement_detail_journal FOREIGN KEY (journal, enterprise) REFERENCES public.journal(id, enterprise);


--
-- Name: accounting_movement_detail accounting_movement_detail_movement; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounting_movement_detail
    ADD CONSTRAINT accounting_movement_detail_movement FOREIGN KEY (movement, enterprise) REFERENCES public.accounting_movement(id, enterprise);


--
-- Name: accounting_movement_detail accounting_movement_detail_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounting_movement_detail
    ADD CONSTRAINT accounting_movement_detail_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: accounting_movement accounting_movement_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounting_movement
    ADD CONSTRAINT accounting_movement_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: address address_country; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_country FOREIGN KEY (country, enterprise) REFERENCES public.country(id, enterprise);


--
-- Name: address address_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_customer FOREIGN KEY (customer, enterprise) REFERENCES public.customer(id, enterprise);


--
-- Name: address address_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: address address_state; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_state FOREIGN KEY (state, enterprise) REFERENCES public.state(id, enterprise);


--
-- Name: api_key api_key_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT api_key_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: api_key api_key_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT api_key_user FOREIGN KEY ("user", enterprise) REFERENCES public."user"(id, config);


--
-- Name: api_key api_key_user_created; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT api_key_user_created FOREIGN KEY (user_created, enterprise) REFERENCES public."user"(id, config);


--
-- Name: billing_series billing_series_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.billing_series
    ADD CONSTRAINT billing_series_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: carrier carrier_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.carrier
    ADD CONSTRAINT carrier_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: charges charges_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.charges
    ADD CONSTRAINT charges_account FOREIGN KEY (account, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: charges charges_accounting_movement; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.charges
    ADD CONSTRAINT charges_accounting_movement FOREIGN KEY (accounting_movement, enterprise) REFERENCES public.accounting_movement(id, enterprise);


--
-- Name: charges charges_accounting_movement_detail_credit; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.charges
    ADD CONSTRAINT charges_accounting_movement_detail_credit FOREIGN KEY (accounting_movement_detail_credit, enterprise) REFERENCES public.accounting_movement_detail(id, enterprise);


--
-- Name: charges charges_accounting_movement_detail_debit; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.charges
    ADD CONSTRAINT charges_accounting_movement_detail_debit FOREIGN KEY (accounting_movement_detail_debit, enterprise) REFERENCES public.accounting_movement_detail(id, enterprise);


--
-- Name: charges charges_collection_operation; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.charges
    ADD CONSTRAINT charges_collection_operation FOREIGN KEY (collection_operation, enterprise) REFERENCES public.collection_operation(id, enterprise);


--
-- Name: charges charges_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.charges
    ADD CONSTRAINT charges_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: collection_operation collection_operation_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.collection_operation
    ADD CONSTRAINT collection_operation_account FOREIGN KEY (account, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: collection_operation collection_operation_accounting_movement; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.collection_operation
    ADD CONSTRAINT collection_operation_accounting_movement FOREIGN KEY (accounting_movement, enterprise) REFERENCES public.accounting_movement(id, enterprise);


--
-- Name: collection_operation collection_operation_accounting_movement_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.collection_operation
    ADD CONSTRAINT collection_operation_accounting_movement_detail FOREIGN KEY (accounting_movement_detail, enterprise) REFERENCES public.accounting_movement_detail(id, enterprise);


--
-- Name: collection_operation collection_operation_bank; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.collection_operation
    ADD CONSTRAINT collection_operation_bank FOREIGN KEY (bank, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: collection_operation collection_operation_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.collection_operation
    ADD CONSTRAINT collection_operation_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: collection_operation collection_operation_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.collection_operation
    ADD CONSTRAINT collection_operation_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: complex_manufacturing_order_manufacturing_order complex_manufacturing_order_complex_manufacturing_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_complex_manufacturing_order FOREIGN KEY (complex_manufacturing_order, enterprise) REFERENCES public.complex_manufacturing_order(id, enterprise);


--
-- Name: complex_manufacturing_order complex_manufacturing_order_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: complex_manufacturing_order_manufacturing_order complex_manufacturing_order_manufacturing_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order FOREIGN KEY (manufacturing_order, enterprise) REFERENCES public.manufacturing_order(id, enterprise);


--
-- Name: complex_manufacturing_order_manufacturing_order complex_manufacturing_order_manufacturing_order_complex_manufac; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_complex_manufac FOREIGN KEY (complex_manufacturing_order_manufacturing_order_output, enterprise) REFERENCES public.complex_manufacturing_order_manufacturing_order(id, enterprise);


--
-- Name: complex_manufacturing_order_manufacturing_order complex_manufacturing_order_manufacturing_order_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: complex_manufacturing_order_manufacturing_order complex_manufacturing_order_manufacturing_order_manufacturing_o; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_manufacturing_o FOREIGN KEY (manufacturing_order_type_component, enterprise) REFERENCES public.manufacturing_order_type_components(id, enterprise);


--
-- Name: complex_manufacturing_order_manufacturing_order complex_manufacturing_order_manufacturing_order_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: complex_manufacturing_order_manufacturing_order complex_manufacturing_order_manufacturing_order_purchase_order_; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_purchase_order_ FOREIGN KEY (purchase_order_detail, enterprise) REFERENCES public.purchase_order_detail(id, enterprise);


--
-- Name: complex_manufacturing_order_manufacturing_order complex_manufacturing_order_manufacturing_order_sale_order_deta; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_sale_order_deta FOREIGN KEY (sale_order_detail, enterprise) REFERENCES public.sales_order_detail(id, enterprise);


--
-- Name: complex_manufacturing_order_manufacturing_order complex_manufacturing_order_manufacturing_order_warehouse_movem; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_warehouse_movem FOREIGN KEY (warehouse_movement) REFERENCES public.warehouse_movement(id);


--
-- Name: complex_manufacturing_order complex_manufacturing_order_type; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_type FOREIGN KEY (type, enterprise) REFERENCES public.manufacturing_order_type(id, enterprise);


--
-- Name: complex_manufacturing_order complex_manufacturing_order_user_manufactured; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_user_manufactured FOREIGN KEY (user_manufactured, enterprise) REFERENCES public."user"(id, config);


--
-- Name: complex_manufacturing_order complex_manufacturing_order_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.complex_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_warehouse FOREIGN KEY (warehouse, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: config_accounts_vat config_accounts_vat_account_purchase; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config_accounts_vat
    ADD CONSTRAINT config_accounts_vat_account_purchase FOREIGN KEY (account_purchase, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: config_accounts_vat config_accounts_vat_account_sales; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config_accounts_vat
    ADD CONSTRAINT config_accounts_vat_account_sales FOREIGN KEY (account_sale, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: config_accounts_vat config_accounts_vat_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config_accounts_vat
    ADD CONSTRAINT config_accounts_vat_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: config config_customer_journal; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_customer_journal FOREIGN KEY (customer_journal, id) REFERENCES public.journal(id, enterprise);


--
-- Name: config config_prestashop_export_serie; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_prestashop_export_serie FOREIGN KEY (prestashop_export_serie, id) REFERENCES public.billing_series(id, enterprise);


--
-- Name: config config_prestashop_interior_serie; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_prestashop_interior_serie FOREIGN KEY (prestashop_interior_serie, id) REFERENCES public.billing_series(id, enterprise);


--
-- Name: config config_prestashop_intracommunity_serie; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_prestashop_intracommunity_serie FOREIGN KEY (prestashop_intracommunity_serie, id) REFERENCES public.billing_series(id, enterprise);


--
-- Name: config config_purchase_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_purchase_account FOREIGN KEY (purchase_account, id) REFERENCES public.account(id, enterprise);


--
-- Name: config config_purchase_journal; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_purchase_journal FOREIGN KEY (purchase_journal, id) REFERENCES public.journal(id, enterprise);


--
-- Name: config config_sales_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_sales_account FOREIGN KEY (sales_account, id) REFERENCES public.account(id, enterprise);


--
-- Name: config config_sales_journal; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_sales_journal FOREIGN KEY (sales_journal, id) REFERENCES public.journal(id, enterprise);


--
-- Name: config config_shopify_default_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_shopify_default_payment_method FOREIGN KEY (shopify_default_payment_method, id) REFERENCES public.payment_method(id, enterprise);


--
-- Name: config config_shopify_export_serie; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_shopify_export_serie FOREIGN KEY (shopify_export_serie, id) REFERENCES public.billing_series(id, enterprise);


--
-- Name: config config_shopify_interior_serie; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_shopify_interior_serie FOREIGN KEY (shopify_interior_serie, id) REFERENCES public.billing_series(id, enterprise);


--
-- Name: config config_shopify_intracommunity_serie; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_shopify_intracommunity_serie FOREIGN KEY (shopify_intracommunity_serie, id) REFERENCES public.billing_series(id, enterprise);


--
-- Name: config config_supplier_journal; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_supplier_journal FOREIGN KEY (supplier_journal, id) REFERENCES public.journal(id, enterprise);


--
-- Name: config config_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_warehouse FOREIGN KEY (default_warehouse, id) REFERENCES public.warehouse(id, enterprise);


--
-- Name: config config_woocommerce_default_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_woocommerce_default_payment_method FOREIGN KEY (woocommerce_default_payment_method, id) REFERENCES public.payment_method(id, enterprise);


--
-- Name: config config_woocommerce_export_serie; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_woocommerce_export_serie FOREIGN KEY (woocommerce_export_serie, id) REFERENCES public.billing_series(id, enterprise);


--
-- Name: config config_woocommerce_interior_serie; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_woocommerce_interior_serie FOREIGN KEY (woocommerce_interior_serie, id) REFERENCES public.billing_series(id, enterprise);


--
-- Name: config config_woocommerce_intracommunity_serie; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_woocommerce_intracommunity_serie FOREIGN KEY (woocommerce_intracommunity_serie, id) REFERENCES public.billing_series(id, enterprise);


--
-- Name: connection_filter connection_filter_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connection_filter
    ADD CONSTRAINT connection_filter_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: connection_filter_user connection_filter_user_connection_filter; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connection_filter_user
    ADD CONSTRAINT connection_filter_user_connection_filter FOREIGN KEY (connection_filter) REFERENCES public.connection_filter(id);


--
-- Name: connection_filter_user connection_filter_user_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connection_filter_user
    ADD CONSTRAINT connection_filter_user_user FOREIGN KEY ("user") REFERENCES public."user"(id);


--
-- Name: connection_log connection_log_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connection_log
    ADD CONSTRAINT connection_log_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: connection_log connection_log_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connection_log
    ADD CONSTRAINT connection_log_user FOREIGN KEY ("user", enterprise) REFERENCES public."user"(id, config);


--
-- Name: country country_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_currency FOREIGN KEY (currency, enterprise) REFERENCES public.currency(id, enterprise);


--
-- Name: country country_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: country country_language; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_language FOREIGN KEY (language, enterprise) REFERENCES public.language(id, enterprise);


--
-- Name: currency currency_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.currency
    ADD CONSTRAINT currency_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: custom_fields custom_fields_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_fields
    ADD CONSTRAINT custom_fields_customer FOREIGN KEY (customer, enterprise) REFERENCES public.customer(id, enterprise);


--
-- Name: custom_fields custom_fields_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_fields
    ADD CONSTRAINT custom_fields_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: custom_fields custom_fields_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_fields
    ADD CONSTRAINT custom_fields_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: custom_fields custom_fields_supplier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_fields
    ADD CONSTRAINT custom_fields_supplier FOREIGN KEY (supplier, enterprise) REFERENCES public.suppliers(id, enterprise);


--
-- Name: customer customer_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_account FOREIGN KEY (account, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: customer customer_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_billing_series FOREIGN KEY (billing_series, enterprise) REFERENCES public.billing_series(id, enterprise);


--
-- Name: customer customer_country; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_country FOREIGN KEY (country, enterprise) REFERENCES public.country(id, enterprise);


--
-- Name: customer customer_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: customer customer_language; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_language FOREIGN KEY (language, enterprise) REFERENCES public.language(id, enterprise);


--
-- Name: customer customer_main_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_main_address FOREIGN KEY (main_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: customer customer_main_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_main_billing_address FOREIGN KEY (main_billing_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: customer customer_main_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_main_shipping_address FOREIGN KEY (main_shipping_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: customer customer_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: customer customer_state; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_state FOREIGN KEY (state, enterprise) REFERENCES public.state(id, enterprise);


--
-- Name: document document_container; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_container FOREIGN KEY (container, enterprise) REFERENCES public.document_container(id, enterprise);


--
-- Name: document_container document_container_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document_container
    ADD CONSTRAINT document_container_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: document document_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: document document_purchase_delivery_note; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_purchase_delivery_note FOREIGN KEY (purchase_delivery_note, enterprise) REFERENCES public.purchase_delivery_note(id, enterprise);


--
-- Name: document document_purchase_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_purchase_invoice FOREIGN KEY (purchase_invoice, enterprise) REFERENCES public.purchase_invoice(id, enterprise);


--
-- Name: document document_purchase_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_purchase_order FOREIGN KEY (purchase_order, enterprise) REFERENCES public.purchase_order(id, enterprise);


--
-- Name: document document_sales_delivery_note; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_sales_delivery_note FOREIGN KEY (sales_delivery_note, enterprise) REFERENCES public.sales_delivery_note(id, enterprise);


--
-- Name: document document_sales_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_sales_invoice FOREIGN KEY (sales_invoice, enterprise) REFERENCES public.sales_invoice(id, enterprise);


--
-- Name: document document_sales_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_sales_order FOREIGN KEY (sales_order, enterprise) REFERENCES public.sales_order(id, enterprise);


--
-- Name: document document_shipping; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_shipping FOREIGN KEY (shipping, enterprise) REFERENCES public.shipping(id, enterprise);


--
-- Name: email_log email_log_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_log
    ADD CONSTRAINT email_log_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: enterprise_logo enterprise_logo_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.enterprise_logo
    ADD CONSTRAINT enterprise_logo_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: group group_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."group"
    ADD CONSTRAINT group_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: incoterm incoterm_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.incoterm
    ADD CONSTRAINT incoterm_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: inventory inventory_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory
    ADD CONSTRAINT inventory_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: inventory_products inventory_products_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_products
    ADD CONSTRAINT inventory_products_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: inventory_products inventory_products_inventory; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_products
    ADD CONSTRAINT inventory_products_inventory FOREIGN KEY (inventory) REFERENCES public.inventory(id);


--
-- Name: inventory_products inventory_products_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_products
    ADD CONSTRAINT inventory_products_product FOREIGN KEY (product) REFERENCES public.product(id);


--
-- Name: inventory_products inventory_products_warehouse_movement; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_products
    ADD CONSTRAINT inventory_products_warehouse_movement FOREIGN KEY (warehouse_movement) REFERENCES public.warehouse_movement(id);


--
-- Name: inventory inventory_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory
    ADD CONSTRAINT inventory_warehouse FOREIGN KEY (warehouse, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: journal journal_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.journal
    ADD CONSTRAINT journal_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: language language_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language
    ADD CONSTRAINT language_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: login_tokens login_tokens_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.login_tokens
    ADD CONSTRAINT login_tokens_user FOREIGN KEY ("user") REFERENCES public."user"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: manufacturing_order manufacturing_order_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: manufacturing_order manufacturing_order_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_order FOREIGN KEY ("order", enterprise) REFERENCES public.sales_order(id, enterprise);


--
-- Name: manufacturing_order manufacturing_order_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_order_detail FOREIGN KEY (order_detail, enterprise) REFERENCES public.sales_order_detail(id, enterprise);


--
-- Name: manufacturing_order manufacturing_order_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: manufacturing_order manufacturing_order_type; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_type FOREIGN KEY (type, enterprise) REFERENCES public.manufacturing_order_type(id, enterprise);


--
-- Name: manufacturing_order_type_components manufacturing_order_type_components_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order_type_components
    ADD CONSTRAINT manufacturing_order_type_components_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: manufacturing_order_type_components manufacturing_order_type_components_manufacturing_order_type; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order_type_components
    ADD CONSTRAINT manufacturing_order_type_components_manufacturing_order_type FOREIGN KEY (manufacturing_order_type, enterprise) REFERENCES public.manufacturing_order_type(id, enterprise);


--
-- Name: manufacturing_order_type_components manufacturing_order_type_components_manufacturing_order_type2; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order_type_components
    ADD CONSTRAINT manufacturing_order_type_components_manufacturing_order_type2 FOREIGN KEY (manufacturing_order_type, enterprise) REFERENCES public.manufacturing_order_type(id, enterprise);


--
-- Name: manufacturing_order_type_components manufacturing_order_type_components_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order_type_components
    ADD CONSTRAINT manufacturing_order_type_components_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: manufacturing_order_type manufacturing_order_type_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order_type
    ADD CONSTRAINT manufacturing_order_type_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: manufacturing_order manufacturing_order_user_created; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_user_created FOREIGN KEY (user_created, enterprise) REFERENCES public."user"(id, config);


--
-- Name: manufacturing_order manufacturing_order_user_manufactured; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_user_manufactured FOREIGN KEY (user_manufactured, enterprise) REFERENCES public."user"(id, config);


--
-- Name: manufacturing_order manufacturing_order_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_warehouse FOREIGN KEY (warehouse, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: manufacturing_order manufacturing_order_warehouse_movement; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.manufacturing_order
    ADD CONSTRAINT manufacturing_order_warehouse_movement FOREIGN KEY (warehouse_movement) REFERENCES public.warehouse_movement(id);


--
-- Name: packages packages_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packages
    ADD CONSTRAINT packages_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: packages packages_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packages
    ADD CONSTRAINT packages_product FOREIGN KEY (product) REFERENCES public.product(id);


--
-- Name: packaging packaging_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packaging
    ADD CONSTRAINT packaging_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: packaging packaging_package; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packaging
    ADD CONSTRAINT packaging_package FOREIGN KEY (package) REFERENCES public.packages(id);


--
-- Name: packaging packaging_pallet; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packaging
    ADD CONSTRAINT packaging_pallet FOREIGN KEY (pallet) REFERENCES public.pallets(id);


--
-- Name: packaging packaging_sales_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packaging
    ADD CONSTRAINT packaging_sales_order FOREIGN KEY (sales_order) REFERENCES public.sales_order(id);


--
-- Name: packaging packaging_shipping; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.packaging
    ADD CONSTRAINT packaging_shipping FOREIGN KEY (shipping) REFERENCES public.shipping(id);


--
-- Name: pallets pallets_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pallets
    ADD CONSTRAINT pallets_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: pallets pallets_sales_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pallets
    ADD CONSTRAINT pallets_sales_order FOREIGN KEY (sales_order) REFERENCES public.sales_order(id);


--
-- Name: payment_method payment_method_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_method
    ADD CONSTRAINT payment_method_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: payment_transaction payment_transaction_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_transaction
    ADD CONSTRAINT payment_transaction_account FOREIGN KEY (account, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: payment_transaction payment_transaction_accounting_movement; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_transaction
    ADD CONSTRAINT payment_transaction_accounting_movement FOREIGN KEY (accounting_movement, enterprise) REFERENCES public.accounting_movement(id, enterprise);


--
-- Name: payment_transaction payment_transaction_accounting_movement_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_transaction
    ADD CONSTRAINT payment_transaction_accounting_movement_detail FOREIGN KEY (accounting_movement_detail, enterprise) REFERENCES public.accounting_movement_detail(id, enterprise);


--
-- Name: payment_transaction payment_transaction_bank; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_transaction
    ADD CONSTRAINT payment_transaction_bank FOREIGN KEY (bank, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: payment_transaction payment_transaction_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_transaction
    ADD CONSTRAINT payment_transaction_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: payment_transaction payment_transaction_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_transaction
    ADD CONSTRAINT payment_transaction_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: payments payments_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_account FOREIGN KEY (account, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: payments payments_accounting_movement; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_accounting_movement FOREIGN KEY (accounting_movement, enterprise) REFERENCES public.accounting_movement(id, enterprise);


--
-- Name: payments payments_accounting_movement_detail_credit; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_accounting_movement_detail_credit FOREIGN KEY (accounting_movement_detail_credit, enterprise) REFERENCES public.accounting_movement_detail(id, enterprise);


--
-- Name: payments payments_accounting_movement_detail_debit; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_accounting_movement_detail_debit FOREIGN KEY (accounting_movement_detail_debit, enterprise) REFERENCES public.accounting_movement_detail(id, enterprise);


--
-- Name: payments payments_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: payments payments_payment_transaction; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_payment_transaction FOREIGN KEY (payment_transaction, enterprise) REFERENCES public.payment_transaction(id, enterprise);


--
-- Name: permission_dictionary permission_dictionary_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permission_dictionary
    ADD CONSTRAINT permission_dictionary_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: permission_dictionary_group permission_dictionary_group_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permission_dictionary_group
    ADD CONSTRAINT permission_dictionary_group_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: permission_dictionary_group permission_dictionary_group_group; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permission_dictionary_group
    ADD CONSTRAINT permission_dictionary_group_group FOREIGN KEY ("group", enterprise) REFERENCES public."group"(id, enterprise);


--
-- Name: permission_dictionary_group permission_dictionary_group_permission; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permission_dictionary_group
    ADD CONSTRAINT permission_dictionary_group_permission FOREIGN KEY (permission_key, enterprise) REFERENCES public.permission_dictionary(key, enterprise);


--
-- Name: product_account poduct_account_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_account
    ADD CONSTRAINT poduct_account_account FOREIGN KEY (account, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: product_account poduct_account_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_account
    ADD CONSTRAINT poduct_account_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: product_account poduct_account_jorunal_account_number; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_account
    ADD CONSTRAINT poduct_account_jorunal_account_number FOREIGN KEY (jorunal, account_number, enterprise) REFERENCES public.account(journal, account_number, enterprise);


--
-- Name: product_account poduct_account_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_account
    ADD CONSTRAINT poduct_account_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: pos_terminals pos_terminals_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pos_terminals
    ADD CONSTRAINT pos_terminals_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: pos_terminals pos_terminals_orders_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pos_terminals
    ADD CONSTRAINT pos_terminals_orders_billing_series FOREIGN KEY (orders_billing_series, enterprise) REFERENCES public.billing_series(id, enterprise);


--
-- Name: pos_terminals pos_terminals_orders_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pos_terminals
    ADD CONSTRAINT pos_terminals_orders_currency FOREIGN KEY (orders_currency, enterprise) REFERENCES public.currency(id, enterprise);


--
-- Name: pos_terminals pos_terminals_orders_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pos_terminals
    ADD CONSTRAINT pos_terminals_orders_customer FOREIGN KEY (orders_customer, enterprise) REFERENCES public.customer(id, enterprise);


--
-- Name: pos_terminals pos_terminals_orders_delivery_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pos_terminals
    ADD CONSTRAINT pos_terminals_orders_delivery_address FOREIGN KEY (orders_delivery_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: pos_terminals pos_terminals_orders_invoice_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pos_terminals
    ADD CONSTRAINT pos_terminals_orders_invoice_address FOREIGN KEY (orders_invoice_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: pos_terminals pos_terminals_orders_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pos_terminals
    ADD CONSTRAINT pos_terminals_orders_payment_method FOREIGN KEY (orders_payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: pos_terminals pos_terminals_orders_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pos_terminals
    ADD CONSTRAINT pos_terminals_orders_warehouse FOREIGN KEY (orders_warehouse, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: product product_color; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_color FOREIGN KEY (color, enterprise) REFERENCES public.color(id, enterprise);


--
-- Name: product product_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: product_family product_family_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_family
    ADD CONSTRAINT product_family_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: product_image product_image_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_image
    ADD CONSTRAINT product_image_product FOREIGN KEY (product) REFERENCES public.product(id);


--
-- Name: product product_manufacturing_order_type; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_manufacturing_order_type FOREIGN KEY (manufacturing_order_type, enterprise) REFERENCES public.manufacturing_order_type(id, enterprise);


--
-- Name: product product_product_family; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_product_family FOREIGN KEY (family, enterprise) REFERENCES public.product_family(id, enterprise);


--
-- Name: product product_supplier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_supplier FOREIGN KEY (supplier, enterprise) REFERENCES public.suppliers(id, enterprise);


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
-- Name: ps_address ps_address_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_address
    ADD CONSTRAINT ps_address_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_carrier ps_carrier_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_carrier
    ADD CONSTRAINT ps_carrier_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_country ps_country_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_country
    ADD CONSTRAINT ps_country_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_currency ps_currency_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_currency
    ADD CONSTRAINT ps_currency_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_customer ps_customer_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_customer
    ADD CONSTRAINT ps_customer_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_language ps_language_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_language
    ADD CONSTRAINT ps_language_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_order_detail ps_order_detail_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_order_detail
    ADD CONSTRAINT ps_order_detail_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_order ps_order_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_order
    ADD CONSTRAINT ps_order_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_product_combination ps_product_combination; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_product_combination
    ADD CONSTRAINT ps_product_combination FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_product ps_product_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_product
    ADD CONSTRAINT ps_product_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_product_option_values ps_product_option_values; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_product_option_values
    ADD CONSTRAINT ps_product_option_values FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_state ps_state_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_state
    ADD CONSTRAINT ps_state_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: ps_zone ps_zone_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ps_zone
    ADD CONSTRAINT ps_zone_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: purchase_delivery_note purchase_delivery_note_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_delivery_note
    ADD CONSTRAINT purchase_delivery_note_billing_series FOREIGN KEY (billing_series, enterprise) REFERENCES public.billing_series(id, enterprise);


--
-- Name: purchase_delivery_note purchase_delivery_note_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_delivery_note
    ADD CONSTRAINT purchase_delivery_note_currency FOREIGN KEY (currency, enterprise) REFERENCES public.currency(id, enterprise);


--
-- Name: purchase_delivery_note purchase_delivery_note_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_delivery_note
    ADD CONSTRAINT purchase_delivery_note_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: purchase_delivery_note purchase_delivery_note_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_delivery_note
    ADD CONSTRAINT purchase_delivery_note_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: purchase_delivery_note purchase_delivery_note_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_delivery_note
    ADD CONSTRAINT purchase_delivery_note_shipping_address FOREIGN KEY (shipping_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: purchase_delivery_note purchase_delivery_note_supplier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_delivery_note
    ADD CONSTRAINT purchase_delivery_note_supplier FOREIGN KEY (supplier, enterprise) REFERENCES public.suppliers(id, enterprise);


--
-- Name: purchase_delivery_note purchase_delivery_note_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_delivery_note
    ADD CONSTRAINT purchase_delivery_note_warehouse FOREIGN KEY (warehouse, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: purchase_invoice purchase_invoice_accounting_movement; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice
    ADD CONSTRAINT purchase_invoice_accounting_movement FOREIGN KEY (accounting_movement, enterprise) REFERENCES public.accounting_movement(id, enterprise);


--
-- Name: purchase_invoice purchase_invoice_amended_purchase_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice
    ADD CONSTRAINT purchase_invoice_amended_purchase_invoice FOREIGN KEY (amended_invoice) REFERENCES public.purchase_invoice(id);


--
-- Name: purchase_invoice purchase_invoice_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice
    ADD CONSTRAINT purchase_invoice_billing_address FOREIGN KEY (billing_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: purchase_invoice purchase_invoice_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice
    ADD CONSTRAINT purchase_invoice_billing_series FOREIGN KEY (billing_series, enterprise) REFERENCES public.billing_series(id, enterprise);


--
-- Name: purchase_invoice purchase_invoice_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice
    ADD CONSTRAINT purchase_invoice_currency FOREIGN KEY (currency, enterprise) REFERENCES public.currency(id, enterprise);


--
-- Name: purchase_invoice_details purchase_invoice_details_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice_details
    ADD CONSTRAINT purchase_invoice_details_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: purchase_invoice_details purchase_invoice_details_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice_details
    ADD CONSTRAINT purchase_invoice_details_invoice FOREIGN KEY (invoice, enterprise) REFERENCES public.purchase_invoice(id, enterprise);


--
-- Name: purchase_invoice_details purchase_invoice_details_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice_details
    ADD CONSTRAINT purchase_invoice_details_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: purchase_invoice purchase_invoice_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice
    ADD CONSTRAINT purchase_invoice_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: purchase_invoice purchase_invoice_supplier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_invoice
    ADD CONSTRAINT purchase_invoice_supplier FOREIGN KEY (supplier, enterprise) REFERENCES public.suppliers(id, enterprise);


--
-- Name: purchase_order purchase_order_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_billing_address FOREIGN KEY (billing_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: purchase_order purchase_order_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_billing_series FOREIGN KEY (billing_series, enterprise) REFERENCES public.billing_series(id, enterprise);


--
-- Name: purchase_order purchase_order_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_currency FOREIGN KEY (currency, enterprise) REFERENCES public.currency(id, enterprise);


--
-- Name: purchase_order_detail purchase_order_detail_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order_detail
    ADD CONSTRAINT purchase_order_detail_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: purchase_order_detail purchase_order_detail_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order_detail
    ADD CONSTRAINT purchase_order_detail_order FOREIGN KEY ("order", enterprise) REFERENCES public.purchase_order(id, enterprise);


--
-- Name: purchase_order_detail purchase_order_detail_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order_detail
    ADD CONSTRAINT purchase_order_detail_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: purchase_order purchase_order_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: purchase_order purchase_order_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: purchase_order purchase_order_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_shipping_address FOREIGN KEY (shipping_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: purchase_order purchase_order_supplier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_supplier FOREIGN KEY (supplier, enterprise) REFERENCES public.suppliers(id, enterprise);


--
-- Name: purchase_order purchase_order_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_order
    ADD CONSTRAINT purchase_order_warehouse FOREIGN KEY (warehouse, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: report_template report_template_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.report_template
    ADD CONSTRAINT report_template_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: report_template_translation report_template_translation_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.report_template_translation
    ADD CONSTRAINT report_template_translation_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: report_template_translation report_template_translation_language; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.report_template_translation
    ADD CONSTRAINT report_template_translation_language FOREIGN KEY (language, enterprise) REFERENCES public.language(id, enterprise);


--
-- Name: sales_delivery_note sales_delivery_note_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_billing_series FOREIGN KEY (billing_series, enterprise) REFERENCES public.billing_series(id, enterprise);


--
-- Name: sales_delivery_note sales_delivery_note_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_currency FOREIGN KEY (currency, enterprise) REFERENCES public.currency(id, enterprise);


--
-- Name: sales_delivery_note sales_delivery_note_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_customer FOREIGN KEY (customer, enterprise) REFERENCES public.customer(id, enterprise);


--
-- Name: sales_delivery_note sales_delivery_note_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sales_delivery_note sales_delivery_note_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: sales_delivery_note sales_delivery_note_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_shipping_address FOREIGN KEY (shipping_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: sales_delivery_note sales_delivery_note_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_delivery_note
    ADD CONSTRAINT sales_delivery_note_warehouse FOREIGN KEY (warehouse, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: sales_invoice sales_invoice_accounting_movement; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_accounting_movement FOREIGN KEY (accounting_movement, enterprise) REFERENCES public.accounting_movement(id, enterprise);


--
-- Name: sales_invoice sales_invoice_amended_sales_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_amended_sales_invoice FOREIGN KEY (amended_invoice) REFERENCES public.sales_invoice(id);


--
-- Name: sales_invoice sales_invoice_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_billing_address FOREIGN KEY (billing_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: sales_invoice sales_invoice_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_billing_series FOREIGN KEY (billing_series, enterprise) REFERENCES public.billing_series(id, enterprise);


--
-- Name: sales_invoice sales_invoice_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_currency FOREIGN KEY (currency, enterprise) REFERENCES public.currency(id, enterprise);


--
-- Name: sales_invoice sales_invoice_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_customer FOREIGN KEY (customer, enterprise) REFERENCES public.customer(id, enterprise);


--
-- Name: sales_invoice_detail sales_invoice_detail_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice_detail
    ADD CONSTRAINT sales_invoice_detail_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sales_invoice_detail sales_invoice_detail_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice_detail
    ADD CONSTRAINT sales_invoice_detail_invoice FOREIGN KEY (invoice, enterprise) REFERENCES public.sales_invoice(id, enterprise);


--
-- Name: sales_invoice_detail sales_invoice_detail_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice_detail
    ADD CONSTRAINT sales_invoice_detail_order_detail FOREIGN KEY (order_detail, enterprise) REFERENCES public.sales_order_detail(id, enterprise);


--
-- Name: sales_invoice_detail sales_invoice_detail_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice_detail
    ADD CONSTRAINT sales_invoice_detail_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: sales_invoice sales_invoice_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sales_invoice sales_invoice_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_invoice
    ADD CONSTRAINT sales_invoice_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: sales_order sales_order_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_billing_address FOREIGN KEY (billing_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: sales_order sales_order_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_billing_series FOREIGN KEY (billing_series, enterprise) REFERENCES public.billing_series(id, enterprise);


--
-- Name: sales_order sales_order_carrier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_carrier FOREIGN KEY (carrier, enterprise) REFERENCES public.carrier(id, enterprise);


--
-- Name: sales_order sales_order_currency; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_currency FOREIGN KEY (currency, enterprise) REFERENCES public.currency(id, enterprise);


--
-- Name: sales_order sales_order_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_customer FOREIGN KEY (customer, enterprise) REFERENCES public.customer(id, enterprise);


--
-- Name: sales_order_detail_digital_product_data sales_order_detail_digital_product_data_sales_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail_digital_product_data
    ADD CONSTRAINT sales_order_detail_digital_product_data_sales_order_detail FOREIGN KEY (detail) REFERENCES public.sales_order_detail(id);


--
-- Name: sales_order_detail sales_order_detail_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail
    ADD CONSTRAINT sales_order_detail_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sales_order_detail_packaged sales_order_detail_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail_packaged
    ADD CONSTRAINT sales_order_detail_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sales_order_detail sales_order_detail_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail
    ADD CONSTRAINT sales_order_detail_order FOREIGN KEY ("order", enterprise) REFERENCES public.sales_order(id, enterprise);


--
-- Name: sales_order_detail_packaged sales_order_detail_packaged_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail_packaged
    ADD CONSTRAINT sales_order_detail_packaged_order_detail FOREIGN KEY (order_detail, enterprise) REFERENCES public.sales_order_detail(id, enterprise);


--
-- Name: sales_order_detail_packaged sales_order_detail_packaged_packaging; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail_packaged
    ADD CONSTRAINT sales_order_detail_packaged_packaging FOREIGN KEY (packaging, enterprise) REFERENCES public.packaging(id, enterprise);


--
-- Name: sales_order_detail sales_order_detail_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail
    ADD CONSTRAINT sales_order_detail_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: sales_order_detail sales_order_detail_purchase_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_detail
    ADD CONSTRAINT sales_order_detail_purchase_order_detail FOREIGN KEY (purchase_order_detail, enterprise) REFERENCES public.purchase_order_detail(id, enterprise);


--
-- Name: sales_order_discount sales_order_discount_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_discount
    ADD CONSTRAINT sales_order_discount_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sales_order_discount sales_order_discount_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_discount
    ADD CONSTRAINT sales_order_discount_order FOREIGN KEY ("order", enterprise) REFERENCES public.sales_order(id, enterprise);


--
-- Name: sales_order_discount sales_order_discount_sales_invoice_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order_discount
    ADD CONSTRAINT sales_order_discount_sales_invoice_detail FOREIGN KEY (sales_invoice_detail, enterprise) REFERENCES public.sales_invoice_detail(id, enterprise);


--
-- Name: sales_order sales_order_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sales_order sales_order_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: sales_order sales_order_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_shipping_address FOREIGN KEY (shipping_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: sales_order sales_order_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales_order
    ADD CONSTRAINT sales_order_warehouse FOREIGN KEY (warehouse, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: shipping shipping_carrier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_carrier FOREIGN KEY (carrier, enterprise) REFERENCES public.carrier(id, enterprise);


--
-- Name: shipping shipping_delivery_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_delivery_address FOREIGN KEY (delivery_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: shipping shipping_delivery_note; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_delivery_note FOREIGN KEY (delivery_note, enterprise) REFERENCES public.sales_delivery_note(id, enterprise);


--
-- Name: shipping shipping_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: shipping shipping_incoterm; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_incoterm FOREIGN KEY (incoterm, enterprise) REFERENCES public.incoterm(id, enterprise);


--
-- Name: shipping shipping_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping
    ADD CONSTRAINT shipping_order FOREIGN KEY ("order", enterprise) REFERENCES public.sales_order(id, enterprise);


--
-- Name: shipping_status_history shipping_status_history_shipping; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping_status_history
    ADD CONSTRAINT shipping_status_history_shipping FOREIGN KEY (shipping) REFERENCES public.shipping(id);


--
-- Name: shipping_tag shipping_tag_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping_tag
    ADD CONSTRAINT shipping_tag_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: shipping_tag shipping_tag_shipping; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipping_tag
    ADD CONSTRAINT shipping_tag_shipping FOREIGN KEY (shipping, enterprise) REFERENCES public.shipping(id, enterprise);


--
-- Name: state state_country; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.state
    ADD CONSTRAINT state_country FOREIGN KEY (country, enterprise) REFERENCES public.country(id, enterprise);


--
-- Name: state state_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.state
    ADD CONSTRAINT state_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: stock stock_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: stock stock_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise) ON DELETE CASCADE;


--
-- Name: stock stock_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_warehouse FOREIGN KEY (warehouse, enterprise) REFERENCES public.warehouse(id, enterprise) ON DELETE CASCADE;


--
-- Name: suppliers suppliers_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_account FOREIGN KEY (account, enterprise) REFERENCES public.account(id, enterprise);


--
-- Name: suppliers suppliers_billing_series; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_billing_series FOREIGN KEY (billing_series, enterprise) REFERENCES public.billing_series(id, enterprise);


--
-- Name: suppliers suppliers_country; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_country FOREIGN KEY (country, enterprise) REFERENCES public.country(id, enterprise);


--
-- Name: suppliers suppliers_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: suppliers suppliers_language; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_language FOREIGN KEY (language, enterprise) REFERENCES public.language(id, enterprise);


--
-- Name: suppliers suppliers_main_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_main_address FOREIGN KEY (main_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: suppliers suppliers_main_billing_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_main_billing_address FOREIGN KEY (main_billing_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: suppliers suppliers_main_shipping_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_main_shipping_address FOREIGN KEY (main_shipping_address, enterprise) REFERENCES public.address(id, enterprise);


--
-- Name: suppliers suppliers_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_payment_method FOREIGN KEY (payment_method, enterprise) REFERENCES public.payment_method(id, enterprise);


--
-- Name: suppliers suppliers_state; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_state FOREIGN KEY (state, enterprise) REFERENCES public.state(id, enterprise);


--
-- Name: sy_addresses sy_addresses_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_addresses
    ADD CONSTRAINT sy_addresses_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sy_customers sy_customers_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_customers
    ADD CONSTRAINT sy_customers_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sy_draft_order_line_item sy_draft_order_line_item_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_draft_order_line_item
    ADD CONSTRAINT sy_draft_order_line_item_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sy_draft_orders sy_draft_orders_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_draft_orders
    ADD CONSTRAINT sy_draft_orders_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sy_order_line_item sy_order_line_item_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_order_line_item
    ADD CONSTRAINT sy_order_line_item_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sy_orders sy_orders_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_orders
    ADD CONSTRAINT sy_orders_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sy_products sy_products_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_products
    ADD CONSTRAINT sy_products_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: sy_variants sy_variants_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sy_variants
    ADD CONSTRAINT sy_variants_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: transactional_log transactional_log_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactional_log
    ADD CONSTRAINT transactional_log_user FOREIGN KEY ("user", enterprise) REFERENCES public."user"(id, config);


--
-- Name: transfer_between_warehouses_detail transfer_between_warehouses_detail_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfer_between_warehouses_detail
    ADD CONSTRAINT transfer_between_warehouses_detail_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: transfer_between_warehouses_detail transfer_between_warehouses_detail_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfer_between_warehouses_detail
    ADD CONSTRAINT transfer_between_warehouses_detail_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: transfer_between_warehouses_detail transfer_between_warehouses_detail_transfer_between_warehouses; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfer_between_warehouses_detail
    ADD CONSTRAINT transfer_between_warehouses_detail_transfer_between_warehouses FOREIGN KEY (transfer_between_warehouses, enterprise) REFERENCES public.transfer_between_warehouses(id, enterprise);


--
-- Name: transfer_between_warehouses_detail transfer_between_warehouses_detail_warehouse_movement_in; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfer_between_warehouses_detail
    ADD CONSTRAINT transfer_between_warehouses_detail_warehouse_movement_in FOREIGN KEY (warehouse_movement_in) REFERENCES public.warehouse_movement(id);


--
-- Name: transfer_between_warehouses_detail transfer_between_warehouses_detail_warehouse_movement_out; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfer_between_warehouses_detail
    ADD CONSTRAINT transfer_between_warehouses_detail_warehouse_movement_out FOREIGN KEY (warehouse_movement_out) REFERENCES public.warehouse_movement(id);


--
-- Name: transfer_between_warehouses transfer_between_warehouses_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfer_between_warehouses
    ADD CONSTRAINT transfer_between_warehouses_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: transfer_between_warehouses transfer_between_warehouses_warehouse_destination; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfer_between_warehouses
    ADD CONSTRAINT transfer_between_warehouses_warehouse_destination FOREIGN KEY (warehouse_destination, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: transfer_between_warehouses transfer_between_warehouses_warehouse_origin; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfer_between_warehouses
    ADD CONSTRAINT transfer_between_warehouses_warehouse_origin FOREIGN KEY (warehouse_origin, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: user user_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_enterprise FOREIGN KEY (config) REFERENCES public.config(id);


--
-- Name: user_group user_group_group; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_group
    ADD CONSTRAINT user_group_group FOREIGN KEY ("group") REFERENCES public."group"(id);


--
-- Name: user_group user_group_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_group
    ADD CONSTRAINT user_group_user FOREIGN KEY ("user") REFERENCES public."user"(id);


--
-- Name: warehouse warehouse_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse
    ADD CONSTRAINT warehouse_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: warehouse_movement warehouse_movement_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: warehouse_movement warehouse_movement_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_product FOREIGN KEY (product, enterprise) REFERENCES public.product(id, enterprise);


--
-- Name: warehouse_movement warehouse_movement_purchase_delivery_note; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_purchase_delivery_note FOREIGN KEY (purchase_delivery_note, enterprise) REFERENCES public.purchase_delivery_note(id, enterprise);


--
-- Name: warehouse_movement warehouse_movement_purchase_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_purchase_order FOREIGN KEY (purchase_order, enterprise) REFERENCES public.purchase_order(id, enterprise);


--
-- Name: warehouse_movement warehouse_movement_purchase_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_purchase_order_detail FOREIGN KEY (purchase_order_detail, enterprise) REFERENCES public.purchase_order_detail(id, enterprise);


--
-- Name: warehouse_movement warehouse_movement_sales_delivery_note; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_sales_delivery_note FOREIGN KEY (sales_delivery_note, enterprise) REFERENCES public.sales_delivery_note(id, enterprise);


--
-- Name: warehouse_movement warehouse_movement_sales_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_sales_order FOREIGN KEY (sales_order, enterprise) REFERENCES public.sales_order(id, enterprise);


--
-- Name: warehouse_movement warehouse_movement_sales_order_detail; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_sales_order_detail FOREIGN KEY (sales_order_detail, enterprise) REFERENCES public.sales_order_detail(id, enterprise);


--
-- Name: warehouse_movement warehouse_movement_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouse_movement
    ADD CONSTRAINT warehouse_movement_warehouse FOREIGN KEY (warehouse, enterprise) REFERENCES public.warehouse(id, enterprise);


--
-- Name: wc_customers wc_customers_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wc_customers
    ADD CONSTRAINT wc_customers_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: wc_order_details wc_order_details_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wc_order_details
    ADD CONSTRAINT wc_order_details_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: wc_orders wc_orders_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wc_orders
    ADD CONSTRAINT wc_orders_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: wc_product_variations wc_product_variations_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wc_product_variations
    ADD CONSTRAINT wc_product_variations_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: wc_products wc_products_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wc_products
    ADD CONSTRAINT wc_products_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: webhook_logs webhook_logs_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_logs
    ADD CONSTRAINT webhook_logs_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: webhook_logs webhook_logs_webhook; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_logs
    ADD CONSTRAINT webhook_logs_webhook FOREIGN KEY (webhook, enterprise) REFERENCES public.webhook_settings(id, enterprise);


--
-- Name: webhook_queue webhook_queue_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_queue
    ADD CONSTRAINT webhook_queue_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- Name: webhook_queue webhook_queue_webhook; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_queue
    ADD CONSTRAINT webhook_queue_webhook FOREIGN KEY (webhook, enterprise) REFERENCES public.webhook_settings(id, enterprise);


--
-- Name: webhook_settings webhook_settings_enterprise; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_settings
    ADD CONSTRAINT webhook_settings_enterprise FOREIGN KEY (enterprise) REFERENCES public.config(id);


--
-- PostgreSQL database dump complete
--

