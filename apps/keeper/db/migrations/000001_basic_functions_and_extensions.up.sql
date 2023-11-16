-- Function trigger for update timestamp
CREATE OR REPLACE FUNCTION trigger_set_update_at_timestamp()
RETURNS TRIGGER AS $trigger_set_update_at_timestamp$
  BEGIN
    IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
      NEW.updated_at = now();
      RETURN NEW;
    ELSE
      RETURN OLD;
    END IF;
  END;
$trigger_set_update_at_timestamp$ LANGUAGE plpgsql;

-- Load extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";
CREATE EXTENSION IF NOT EXISTS "btree_gist";
CREATE EXTENSION IF NOT EXISTS "hstore";
