ALTER TABLE public.api_key
    ADD COLUMN auth character(1) NOT NULL DEFAULT 'P';

ALTER TABLE public.api_key
    ALTER COLUMN token DROP NOT NULL;

ALTER TABLE public.api_key
    ADD COLUMN basic_auth_user character varying(20);

ALTER TABLE public.api_key
    ADD COLUMN basic_auth_password character varying(20);

CREATE UNIQUE INDEX api_key_basic_auth
    ON public.api_key USING btree
    (basic_auth_user ASC NULLS LAST, basic_auth_password ASC NULLS LAST)

    WHERE auth = 'B';

DROP INDEX public.api_key_token;

CREATE UNIQUE INDEX api_key_token
    ON public.api_key USING btree
    (token ASC NULLS LAST)

    WHERE token IS NOT NULL;
