ALTER TABLE public.config
    ADD COLUMN smtp_starttls boolean NOT NULL DEFAULT false;