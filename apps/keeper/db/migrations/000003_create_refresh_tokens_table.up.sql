-- Create refresh tokens table
CREATE TABLE "refresh_tokens" (
    "rtid" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "refresh_token" VARCHAR(255) NOT NULL,
    "expires" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT "user_id_fk" FOREIGN KEY("user_id") REFERENCES "users"("uid")
);

-- Indexes
CREATE INDEX "rtuid_idx" ON "refresh_tokens"("user_id");
