ALTER TABLE public.config
    ADD COLUMN smtp_reply_to character varying(50) DEFAULT '';