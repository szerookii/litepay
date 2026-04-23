-- Create "users" table
CREATE TABLE "public"."users" (
 "id" uuid NOT NULL,
 "create_time" timestamptz NOT NULL,
 "update_time" timestamptz NOT NULL,
 "email" character varying NOT NULL,
 "password_hash" character varying NOT NULL,
 "api_key" character varying NOT NULL,
 "btc_xpub" character varying NULL,
 "ltc_xpub" character varying NULL,
 "sol_address" character varying NULL,
 "btc_address_index" bigint NOT NULL DEFAULT 0,
 "ltc_address_index" bigint NOT NULL DEFAULT 0,
 "webhook_url" character varying NULL,
 PRIMARY KEY ("id")
);
-- Create index "user_api_key" to table: "users"
CREATE UNIQUE INDEX "user_api_key" ON "public"."users" ("api_key");
-- Create index "user_email" to table: "users"
CREATE UNIQUE INDEX "user_email" ON "public"."users" ("email");
-- Create index "users_api_key_key" to table: "users"
CREATE UNIQUE INDEX "users_api_key_key" ON "public"."users" ("api_key");
-- Create "payments" table
CREATE TABLE "public"."payments" (
 "id" uuid NOT NULL,
 "create_time" timestamptz NOT NULL,
 "update_time" timestamptz NOT NULL,
 "wallet_address" character varying NOT NULL,
 "sol_reference" character varying NULL,
 "address_index" bigint NOT NULL DEFAULT 0,
 "amount_crypto" double precision NOT NULL,
 "currency_crypto" character varying NOT NULL,
 "amount_fiat" double precision NOT NULL,
 "currency_fiat" character varying NOT NULL,
 "transaction_hash" character varying NULL,
 "status" character varying NOT NULL DEFAULT 'PENDING',
 "expires_at" timestamptz NOT NULL,
 "user_id" uuid NOT NULL,
 PRIMARY KEY ("id"),
 CONSTRAINT "payments_users_payments" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "payment_sol_reference" to table: "payments"
CREATE INDEX "payment_sol_reference" ON "public"."payments" ("sol_reference");
-- Create index "payment_user_id" to table: "payments"
CREATE INDEX "payment_user_id" ON "public"."payments" ("user_id");
-- Create index "payment_wallet_address" to table: "payments"
CREATE INDEX "payment_wallet_address" ON "public"."payments" ("wallet_address");
