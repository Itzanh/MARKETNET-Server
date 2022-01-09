CREATE TABLE public.report_template_translation
(
    enterprise integer NOT NULL,
    key character varying(50) COLLATE pg_catalog."default" NOT NULL,
    language integer NOT NULL,
    translation character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT report_template_translation_pkey PRIMARY KEY (enterprise, key, language),
    CONSTRAINT report_template_translation_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT report_template_translation_language FOREIGN KEY (enterprise, language)
        REFERENCES public.language (enterprise, id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);