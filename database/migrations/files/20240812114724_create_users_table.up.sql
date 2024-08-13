-- create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "admin" boolean NOT NULL DEFAULT false,
  "created_at" bigint NOT NULL DEFAULT EXTRACT(epoch FROM now()),
  "updated_at" bigint NOT NULL DEFAULT EXTRACT(epoch FROM now()),
  "deleted_at" bigint NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
