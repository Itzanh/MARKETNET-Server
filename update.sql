ALTER TABLE public.journal
    ALTER COLUMN id TYPE integer;

ALTER TABLE public.account
    ALTER COLUMN journal TYPE integer;

ALTER TABLE public.accounting_movement_detail
    ALTER COLUMN journal TYPE integer;

ALTER TABLE public.config
    ALTER COLUMN customer_journal TYPE integer;

ALTER TABLE public.config
    ALTER COLUMN sales_journal TYPE integer;

ALTER TABLE public.config
    ALTER COLUMN supplier_journal TYPE integer;

ALTER TABLE public.config
    ALTER COLUMN purchase_journal TYPE integer;