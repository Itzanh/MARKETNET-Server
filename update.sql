CREATE TABLE public.permission_dictionary
(
    enterprise integer NOT NULL,
    key character varying(150) NOT NULL,
    description character varying(250) NOT NULL,
    PRIMARY KEY (enterprise, key),
    CONSTRAINT permission_dictionary_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE UNIQUE INDEX group_id_enterprise
    ON public."group" USING btree
    (id ASC NULLS LAST, enterprise ASC NULLS LAST)
;

CREATE TABLE public.permission_dictionary_group
(
    "group" integer NOT NULL,
    permission_key character varying(150) NOT NULL,
    enterprise integer NOT NULL,
    PRIMARY KEY (enterprise, permission_key, "group"),
    CONSTRAINT permission_dictionary_group_enterprise FOREIGN KEY (enterprise)
        REFERENCES public.config (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT permission_dictionary_group_group FOREIGN KEY ("group", enterprise)
        REFERENCES public."group" (id, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT permission_dictionary_group_permission FOREIGN KEY (permission_key, enterprise)
        REFERENCES public.permission_dictionary (key, enterprise) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);