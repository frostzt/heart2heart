-- Create refresh tokens table
CREATE TABLE "refresh_tokens" (
    "rtid" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "refresh_token" VARCHAR(36) NOT NULL,
    "expires" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT "user_id_fk" FOREIGN KEY("user_id") REFERENCES "users"("uid")
);

-- Indexes
CREATE INDEX "rtuid_idx" ON "refresh_tokens"("user_id");
CREATE INDEX "refresh_token_idx" ON "refresh_tokens"("refresh_token");
