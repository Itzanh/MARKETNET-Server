ALTER TABLE public.config
    ADD COLUMN cron_clear_logs character varying(25) NOT NULL DEFAULT '@monthly';
    
CREATE TABLE public.logs
(
    id bigint NOT NULL,
    date_created timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    title character varying(255) NOT NULL,
    info text NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE public.logs
    OWNER to postgres;

CREATE OR REPLACE FUNCTION set_logs_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(logs.id) END AS id FROM logs) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_logs_id
before insert on logs
for each row execute procedure set_logs_id();

GRANT ALL ON TABLE public.logs TO marketnet;