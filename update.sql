CREATE INDEX accounting_movement_date_created
    ON public.accounting_movement USING btree
    (date_created DESC NULLS LAST)
;

ALTER TABLE public.accounting_movement
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.accounting_movement_detail
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.api_key
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.charges
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

CREATE INDEX collection_operation_status_enterprise
    ON public.collection_operation USING btree
    (status ASC NULLS LAST, enterprise ASC NULLS LAST)
;

ALTER TABLE public.collection_operation
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.collection_operation
    ALTER COLUMN date_expiration TYPE timestamp(3) with time zone ;

CREATE INDEX collection_operation_date_created
    ON public.collection_operation USING btree
    (date_created ASC NULLS LAST)
;

ALTER TABLE public.config
    ALTER COLUMN limit_accounting_date TYPE timestamp(0) with time zone ;

ALTER TABLE public.connection_filter
    ALTER COLUMN time_start TYPE time(0) with time zone ;

ALTER TABLE public.connection_filter
    ALTER COLUMN time_end TYPE time(0) with time zone ;

ALTER TABLE public.connection_log
    ALTER COLUMN date_connected TYPE timestamp(3) with time zone ;

ALTER TABLE public.connection_log
    ALTER COLUMN date_disconnected TYPE timestamp(3) with time zone ;

ALTER TABLE public.customer
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.document
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.document
    ALTER COLUMN date_updated TYPE timestamp(3) with time zone ;

ALTER TABLE public.document_container
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.login_tokens
    ALTER COLUMN date_last_used TYPE timestamp(3) with time zone ;

ALTER TABLE public.logs
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.manufacturing_order
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.manufacturing_order
    ALTER COLUMN date_last_update TYPE timestamp(3) with time zone ;

ALTER TABLE public.manufacturing_order
    ALTER COLUMN date_manufactured TYPE timestamp(3) with time zone ;

ALTER TABLE public.manufacturing_order
    ALTER COLUMN date_tag_printed TYPE timestamp(3) with time zone ;

ALTER TABLE public.payment_transaction
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.payment_transaction
    ALTER COLUMN date_expiration TYPE timestamp(3) with time zone ;

ALTER TABLE public.payments
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.product
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.ps_address
    ALTER COLUMN date_add TYPE timestamp(0) with time zone ;

ALTER TABLE public.ps_address
    ALTER COLUMN date_upd TYPE timestamp(0) with time zone ;

ALTER TABLE public.ps_customer
    ALTER COLUMN date_add TYPE timestamp(0) with time zone ;

ALTER TABLE public.ps_customer
    ALTER COLUMN date_upd TYPE timestamp(0) with time zone ;

ALTER TABLE public.ps_order
    ALTER COLUMN date_add TYPE timestamp(0) with time zone ;

ALTER TABLE public.ps_order
    ALTER COLUMN date_upd TYPE timestamp(0) with time zone ;

ALTER TABLE public.ps_product
    ALTER COLUMN date_add TYPE timestamp(0) with time zone ;

ALTER TABLE public.ps_product
    ALTER COLUMN date_upd TYPE timestamp(0) with time zone ;

ALTER TABLE public.purchase_delivery_note
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.purchase_invoice
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.purchase_order
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.purchase_order
    ALTER COLUMN date_paid TYPE timestamp(3) with time zone ;

ALTER TABLE public.sales_delivery_note
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.sales_invoice
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.sales_order
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.sales_order
    ALTER COLUMN date_payment_accepted TYPE timestamp(3) with time zone ;

ALTER TABLE public.shipping
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.shipping
    ALTER COLUMN date_sent TYPE timestamp(3) with time zone ;

ALTER TABLE public.shipping_tag
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.suppliers
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.transactional_log
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public."user"
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public."user"
    ALTER COLUMN date_last_pwd TYPE timestamp(3) with time zone ;

ALTER TABLE public."user"
    ALTER COLUMN date_last_login TYPE timestamp(3) with time zone ;

ALTER TABLE public.warehouse_movement
    ALTER COLUMN date_created TYPE timestamp(3) with time zone ;

ALTER TABLE public.wc_customers
    ALTER COLUMN date_created TYPE timestamp(0) with time zone ;

ALTER TABLE public.wc_orders
    ALTER COLUMN date_created TYPE timestamp(0) with time zone ;

ALTER TABLE public.wc_products
    ALTER COLUMN date_created TYPE timestamp(0) with time zone ;