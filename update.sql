ALTER TABLE public.product
    ADD COLUMN digital_product boolean NOT NULL DEFAULT false;

CREATE TABLE public.sales_order_detail_digital_product_data
(
    id integer NOT NULL,
    detail bigint NOT NULL,
    key character varying(50) NOT NULL,
    value character varying(250) NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE public.sales_order_detail_digital_product_data
    OWNER to postgres;

CREATE OR REPLACE FUNCTION set_sales_order_detail_digital_product_data_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_order_detail_digital_product_data.id) END AS id FROM sales_order_detail_digital_product_data) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_sales_order_detail_digital_product_data_id
before insert on sales_order_detail_digital_product_data
for each row execute procedure set_sales_order_detail_digital_product_data_id();

ALTER TABLE public.sales_order_detail_digital_product_data
    ADD CONSTRAINT sales_order_detail_digital_product_data_sales_order_detail FOREIGN KEY (detail)
    REFERENCES public.sales_order_detail (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;