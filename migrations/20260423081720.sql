-- Modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "btc_xpub", DROP COLUMN "ltc_xpub", DROP COLUMN "sol_address", DROP COLUMN "btc_address_index", DROP COLUMN "ltc_address_index", ADD COLUMN "account_index" bigint NOT NULL DEFAULT 0;
-- Backfill unique account_index for existing rows
UPDATE "public"."users" SET "account_index" = sub.rn - 1
FROM (SELECT "id", ROW_NUMBER() OVER (ORDER BY "create_time") AS rn FROM "public"."users") sub
WHERE "public"."users"."id" = sub."id";
-- Drop the default (index is assigned in code at registration)
ALTER TABLE "public"."users" ALTER COLUMN "account_index" DROP DEFAULT;
-- Create index "users_account_index_key" to table: "users"
CREATE UNIQUE INDEX "users_account_index_key" ON "public"."users" ("account_index");
