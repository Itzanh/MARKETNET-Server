ALTER TABLE public.api_key
    ADD COLUMN token uuid NOT NULL;
ALTER TABLE public.config
    ADD COLUMN enable_api_key boolean NOT NULL DEFAULT false;