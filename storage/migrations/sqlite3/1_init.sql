-- +migrate Up

CREATE TABLE IF NOT EXISTS "users" (
  "id"       integer primary key autoincrement,
  "username" varchar(191) NOT NULL UNIQUE,
  "password" varchar(191),
  "api_key"  varchar(191) NOT NULL UNIQUE,
  "auth_key" varchar(191) NOT NULL,
  "created_at" datetime,
  "updated_at" datetime
);

CREATE TABLE IF NOT EXISTS "ses_keys" (
  "id"         integer primary key autoincrement,
  "user_id"    integer,
  "access_key" varchar(191) NOT NULL,
  "secret_key" varchar(191) NOT NULL,
  "region"     varchar(30) NOT NULL,
  "created_at" datetime,
  "updated_at" datetime,
  UNIQUE("user_id")
);

CREATE TABLE IF NOT EXISTS "campaigns" (
  "id"            integer primary key autoincrement,
  "user_id"       integer,
  "name"          varchar(191) NOT NULL,
  "template_name" varchar(191) NOT NULL,
  "status"        varchar(191),
  "created_at"    datetime,
  "updated_at"    datetime,
  "scheduled_at"  datetime DEFAULT NULL,
  "completed_at"  datetime DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS i_user ON "campaigns" (user_id);

CREATE TABLE IF NOT EXISTS "subscribers" (
  "id"         integer primary key autoincrement,
  "user_id"    integer,
  "name"       varchar(191) NOT NULL,
  "email"      varchar(191) NOT NULL,
  "blacklisted" integer,
  "active"      integer,
  "created_at" datetime,
  "updated_at" datetime,
  UNIQUE("user_id", "email")
);

CREATE INDEX IF NOT EXISTS i_user ON "subscribers" (user_id);
CREATE INDEX IF NOT EXISTS i_user_blacklist_active ON "subscribers" (user_id, blacklisted, active);

CREATE TABLE IF NOT EXISTS "lists" (
  "id"          integer primary key autoincrement,
  "user_id"     integer,
  "name"        varchar(191),
  "created_at"  datetime,
  "updated_at"  datetime
);

CREATE INDEX IF NOT EXISTS i_user ON "lists" (user_id);

CREATE TABLE IF NOT EXISTS "subscribers_lists" (
  "list_id"       integer,
  "subscriber_id" integer,
  UNIQUE("list_id", "subscriber_id")
);

CREATE INDEX IF NOT EXISTS i_list       ON "subscribers_lists" (list_id);
CREATE INDEX IF NOT EXISTS i_subscriber ON "subscribers_lists" (subscriber_id);

CREATE TABLE IF NOT EXISTS "subscriber_metadata" (
  "id"            integer primary key autoincrement,
  "subscriber_id" integer,
  "key"           varchar(191),
  "value"         varchar(191),
  "created_at"    datetime,
  "updated_at"    datetime
);

CREATE INDEX IF NOT EXISTS i_subscriber ON "subscriber_metadata" (subscriber_id);

CREATE TABLE IF NOT EXISTS "bounces" (
  "id"              integer primary key autoincrement,
  "campaign_id"     integer,
  "user_id"         integer,
  "recipient"       varchar(191),
  "type"            varchar(30),
  "sub_type"        varchar(30),
  "action"          varchar(191),
  "status"          varchar(191),
  "diagnostic_code" varchar(191),
  "feedback_id"     varchar(191),
  "created_at"      datetime
);

CREATE TABLE IF NOT EXISTS "complaints" (
  "id"              integer primary key autoincrement,
  "campaign_id"     integer,
  "user_id"         integer,
  "recipient"       varchar(191),
  "type"            varchar(30),
  "user_agent"      varchar(191),
  "feedback_id"     varchar(191),
  "created_at"      datetime
);

CREATE TABLE IF NOT EXISTS "clicks" (
  "id"              integer primary key autoincrement,
  "campaign_id"     integer,
  "user_id"         integer,
  "ip_address"      varchar(50),
  "user_agent"      varchar(191),
  "link"            varchar(191),
  "created_at"      datetime
);

CREATE TABLE IF NOT EXISTS "opens" (
  "id"              integer primary key autoincrement,
  "campaign_id"     integer,
  "user_id"         integer,
  "ip_address"      varchar(50),
  "user_agent"      varchar(191),
  "created_at"      datetime
);

CREATE TABLE IF NOT EXISTS "send_bulk_logs" (
  "id"            integer primary key autoincrement,
  "uuid"          varchar(36) NOT NULL,
  "user_id"       integer,
  "campaign_id"   integer,
  "message_id"    varchar(191) NOT NULL,
  "status"        varchar(191) NOT NULL,
  "created_at"    datetime
);

CREATE INDEX IF NOT EXISTS i_user_campaign ON "send_bulk_logs" (user_id, campaign_id);
CREATE INDEX IF NOT EXISTS i_uuid ON "send_bulk_logs" (uuid);

-- +migrate Down

DROP TABLE "users";
DROP TABLE "templates";
DROP TABLE "campaigns";
DROP TABLE "lists";
DROP TABLE "subscribers";
DROP TABLE "subscribers_lists";
DROP TABLE "list_metadata";
DROP TABLE "subscriber_metadata";
DROP TABLE "sent_emails";
DROP TABLE "bounces";
DROP TABLE "send_bulk_logs";
