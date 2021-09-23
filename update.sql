CREATE TABLE public.report_template
(
    enterprise integer NOT NULL,
    key character varying(50) COLLATE pg_catalog."default" NOT NULL,
    html text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT report_template_pkey PRIMARY KEY (enterprise, key),
    CONSTRAINT report_template_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);