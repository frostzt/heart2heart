-- Create users table
CREATE TABLE "users" (
    "uid" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "birth_date" DATE NOT NULL,
    "username" VARCHAR(255) NOT NULL UNIQUE,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" VARCHAR(60) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ
);

-- In case of any update made to the row update the `updated_at` timestamp
CREATE TRIGGER set_updated_at BEFORE UPDATE ON "users"
  FOR EACH ROW EXECUTE PROCEDURE trigger_set_update_at_timestamp();
