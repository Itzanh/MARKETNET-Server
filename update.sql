CREATE TABLE public.connection_log
(
    id bigint NOT NULL,
    date_connected timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    date_disconnected timestamp(3) without time zone,
    "user" smallint NOT NULL,
    ok boolean NOT NULL,
    ip_address character varying(15) NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT connection_log_user FOREIGN KEY ("user")
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE public.connection_log
    OWNER to postgres;

CREATE OR REPLACE FUNCTION set_connection_log_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(connection_log.id) END AS id FROM connection_log) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_connection_log_id
before insert on connection_log
for each row execute procedure set_connection_log_id();

CREATE TABLE public.connection_filter
(
    id smallint NOT NULL,
    name character varying(100) NOT NULL,
    type character(1) NOT NULL,
    ip_address character varying(15),
    time_start time(0) without time zone,
    time_end time(0) without time zone,
    PRIMARY KEY (id)
);

ALTER TABLE public.connection_filter
    OWNER to postgres;

CREATE TABLE public.connection_filter_user
(
    connection_filter smallint NOT NULL,
    "user" smallint NOT NULL,
    PRIMARY KEY (connection_filter, "user"),
    CONSTRAINT connection_filter_user_connection_filter FOREIGN KEY (connection_filter)
        REFERENCES public.connection_filter (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT connection_filter_user_user FOREIGN KEY ("user")
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE public.connection_filter_user
    OWNER to postgres;

CREATE OR REPLACE FUNCTION set_connection_filter_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(connection_filter.id) END AS id FROM connection_filter) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_connection_filter_id
before insert on connection_filter
for each row execute procedure set_connection_filter_id();

ALTER TABLE public.config
    ADD COLUMN connection_log boolean NOT NULL DEFAULT false;

ALTER TABLE public.config
    ADD COLUMN filter_connections boolean NOT NULL DEFAULT false;