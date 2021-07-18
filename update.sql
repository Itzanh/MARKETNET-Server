ALTER INDEX public.customer_name RENAME TO customer_name_trgm;
CREATE UNIQUE INDEX customer_name
    ON public.customer USING btree
    (name ASC NULLS LAST)
;
ALTER INDEX public.supplier_name RENAME TO supplier_name_trgm;
CREATE UNIQUE INDEX supplier_name
    ON public.suppliers USING btree
    (name COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
    