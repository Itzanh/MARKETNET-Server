CREATE UNIQUE INDEX manufacturing_order_type_components_component
    ON public.manufacturing_order_type_components USING btree
    (manufacturing_order_type ASC NULLS LAST, product ASC NULLS LAST)
;