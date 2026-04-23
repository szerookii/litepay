-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "webhook_secret" character varying NULL;
