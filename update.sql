ALTER TABLE public.config
    ADD COLUMN undo_manufacturing_order_seconds smallint NOT NULL DEFAULT 120;