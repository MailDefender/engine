-- Create "mailbox_histories" table
CREATE TABLE "public"."mailbox_histories" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "action_type" text NULL,
  "name" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_mailbox_histories_deleted_at" to table: "mailbox_histories"
CREATE INDEX "idx_mailbox_histories_deleted_at" ON "public"."mailbox_histories" ("deleted_at");
-- Create "notifications" table
CREATE TABLE "public"."notifications" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "type" text NULL,
  "channel" text NULL,
  "recipient" text NULL,
  "content" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_notifications_deleted_at" to table: "notifications"
CREATE INDEX "idx_notifications_deleted_at" ON "public"."notifications" ("deleted_at");
-- Create "reputations" table
CREATE TABLE "public"."reputations" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NULL,
  "status" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_reputations_deleted_at" to table: "reputations"
CREATE INDEX "idx_reputations_deleted_at" ON "public"."reputations" ("deleted_at");
-- Create "rules" table
CREATE TABLE "public"."rules" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "uuid" character varying(100) NOT NULL DEFAULT NULL::character varying,
  "category" text NULL,
  "destination" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_rules_uuid" UNIQUE ("uuid")
);
-- Create index "idx_rules_deleted_at" to table: "rules"
CREATE INDEX "idx_rules_deleted_at" ON "public"."rules" ("deleted_at");
-- Create "validation_tokens" table
CREATE TABLE "public"."validation_tokens" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NULL,
  "token" text NULL,
  "admin_token" boolean NULL,
  "generation_date" timestamptz NULL,
  "expiry_date" timestamptz NULL,
  "validated" boolean NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_validation_tokens_token" UNIQUE ("token")
);
-- Create index "idx_validation_tokens_deleted_at" to table: "validation_tokens"
CREATE INDEX "idx_validation_tokens_deleted_at" ON "public"."validation_tokens" ("deleted_at");
-- Create "criterions" table
CREATE TABLE "public"."criterions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "rule_id" bigint NULL,
  "parent_criterion_id" bigint NULL,
  "type" text NULL,
  "input" text NULL,
  "values" text[] NULL,
  "count" bigint NULL,
  "should_exact_count" boolean NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_criterions_sub_criteria" FOREIGN KEY ("id") REFERENCES "public"."criterions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_rules_criteria" FOREIGN KEY ("rule_id") REFERENCES "public"."rules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_criterions_deleted_at" to table: "criterions"
CREATE INDEX "idx_criterions_deleted_at" ON "public"."criterions" ("deleted_at");
-- Create "messages" table
CREATE TABLE "public"."messages" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "message_id" text NULL,
  "sender_name" text NULL,
  "sender_email" text NULL,
  "subject" text NULL,
  "received_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_messages_message_id" UNIQUE ("message_id")
);
-- Create index "idx_messages_deleted_at" to table: "messages"
CREATE INDEX "idx_messages_deleted_at" ON "public"."messages" ("deleted_at");
-- Create "message_histories" table
CREATE TABLE "public"."message_histories" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "action_type" text NULL,
  "message_hash" text NULL,
  "source" text NULL,
  "destination" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_message_histories_message" FOREIGN KEY ("message_hash") REFERENCES "public"."messages" ("message_id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_message_histories_deleted_at" to table: "message_histories"
CREATE INDEX "idx_message_histories_deleted_at" ON "public"."message_histories" ("deleted_at");
-- Create "pending_messages" table
CREATE TABLE "public"."pending_messages" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "message_id" text NULL,
  "sender_email" text NULL,
  "subject" text NULL,
  "mailbox" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_pending_messages_message" FOREIGN KEY ("message_id") REFERENCES "public"."messages" ("message_id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_pending_messages_deleted_at" to table: "pending_messages"
CREATE INDEX "idx_pending_messages_deleted_at" ON "public"."pending_messages" ("deleted_at");
