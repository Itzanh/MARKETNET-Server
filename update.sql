CREATE UNIQUE INDEX complex_manufacturing_order_manufacturing_order_id_enterprise
    ON public.complex_manufacturing_order_manufacturing_order USING btree
    (id ASC NULLS LAST, enterprise ASC NULLS LAST)
;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD COLUMN complex_manufacturing_order_manufacturing_order_output bigint;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    VALIDATE CONSTRAINT complex_manufacturing_order_manufacturing_order_sale_order_deta;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    VALIDATE CONSTRAINT complex_manufacturing_order_manufacturing_order_purchase_order_;

ALTER TABLE public.complex_manufacturing_order_manufacturing_order
    ADD CONSTRAINT complex_manufacturing_order_manufacturing_order_complex_manufacturing_order_manufacturing_order_output FOREIGN KEY (complex_manufacturing_order_manufacturing_order_output, enterprise)
    REFERENCES public.complex_manufacturing_order_manufacturing_order (id, enterprise) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;