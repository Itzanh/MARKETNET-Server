CREATE OR REPLACE FUNCTION set_config_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id = (SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(config.id) END AS id FROM config) + 1;
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

create trigger set_config_id
before insert on config
for each row execute procedure set_config_id();