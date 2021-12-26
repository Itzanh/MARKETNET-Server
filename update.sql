CREATE TABLE public.poduct_account
(
    id integer NOT NULL,
    product integer NOT NULL,
    account integer NOT NULL,
    jorunal integer NOT NULL,
    account_number integer NOT NULL,
    enterprise integer NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT poduct_account_product FOREIGN KEY (product, enterprise)
        REFERENCES public.product (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT poduct_account_account FOREIGN KEY (account, enterprise)
        REFERENCES public.account (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT poduct_account_jorunal_account_number FOREIGN KEY (jorunal, account_number, enterprise)
        REFERENCES public.account (journal, account_number, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT poduct_account_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE public.poduct_account
    RENAME TO product_account;

CREATE OR REPLACE FUNCTION set_product_account_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(product_account.id) END AS id FROM product_account) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_product_account_id
before insert on product_account
for each row execute procedure set_product_account_id();

ALTER TABLE public.product_account
    ADD COLUMN type character(1) NOT NULL;

CREATE UNIQUE INDEX product_account_product_type
    ON public.product_account USING btree
    (product ASC NULLS LAST, type ASC NULLS LAST)
;