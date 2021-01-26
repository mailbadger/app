-- +migrate Up

CREATE TABLE IF NOT EXISTS "users" (
  "id"         integer primary key autoincrement,
  "uuid"       varchar(36) NOT NULL UNIQUE,
  "username"   varchar(191) NOT NULL UNIQUE,
  "password"   varchar(191),
  "source"     varchar(191) NOT NULL,
  "active"     integer,
  "verified"   integer,
  "created_at" datetime,
  "updated_at" datetime
);

CREATE TABLE IF NOT EXISTS "sessions" (
  `id` integer primary key autoincrement,
  `user_id`    integer NOT NULL,
  `session_id` varchar(191) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  UNIQUE("session_id")
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
  "template_id"   integer,
  "status"        varchar(191),
  "created_at"    datetime,
  "updated_at"    datetime,
  "scheduled_at"  datetime DEFAULT NULL,
  "completed_at"  datetime DEFAULT NULL,
  "deleted_at"    datetime DEFAULT NULL,
  "started_at"    datetime DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS i_user ON "campaigns" (user_id);

CREATE TABLE IF NOT EXISTS "subscribers" (
  "id"          integer primary key autoincrement,
  "user_id"     integer,
  "name"        varchar(191),
  "email"       varchar(191) NOT NULL,
  "metadata"    json,
  "blacklisted" integer,
  "active"      integer,
  "created_at"  datetime,
  "updated_at"  datetime,
  UNIQUE("user_id", "email")
);

CREATE INDEX IF NOT EXISTS i_user ON "subscribers" (user_id);
CREATE INDEX IF NOT EXISTS i_user_blacklist_active ON "subscribers" (user_id, blacklisted, active);

CREATE TABLE IF NOT EXISTS "segments" (
  "id"          integer primary key autoincrement,
  "user_id"     integer,
  "name"        varchar(191),
  "created_at"  datetime,
  "updated_at"  datetime
);

CREATE INDEX IF NOT EXISTS i_user ON "segments" (user_id);

CREATE TABLE IF NOT EXISTS "subscribers_segments" (
  "segment_id"    integer,
  "subscriber_id" integer,
  UNIQUE("segment_id", "subscriber_id")
);

CREATE INDEX IF NOT EXISTS i_segment    ON "subscribers_segments" (segment_id);
CREATE INDEX IF NOT EXISTS i_subscriber ON "subscribers_segments" (subscriber_id);

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
  "recipient"       varchar(191),
  "ip_address"      varchar(50),
  "user_agent"      varchar(191),
  "link"            varchar(191),
  "created_at"      datetime
);

CREATE TABLE IF NOT EXISTS "opens" (
  "id"              integer primary key autoincrement,
  "campaign_id"     integer,
  "user_id"         integer,
  `recipient`       varchar(191),
  "ip_address"      varchar(50),
  "user_agent"      varchar(191),
  "created_at"      datetime
);

CREATE TABLE IF NOT EXISTS "deliveries" (
  "id"                     integer primary key autoincrement,
  "campaign_id"            integer,
  "user_id"                integer,
  "recipient"              varchar(191),
  "processing_time_millis" integer,
  "smtp_response"          varchar(191),
  "reporting_mta"          varchar(191),
  "remote_mta_ip"          varchar(50),
  "created_at"             datetime
);

CREATE TABLE IF NOT EXISTS "send_logs" (
  "id"            integer primary key autoincrement,
  "uuid"          varchar(36) NOT NULL,
  "user_id"       integer,
  "campaign_id"   integer,
  "subscriber_id" integer,
  "status"        varchar(191),
  "description"   varchar(191),
  "created_at"    datetime
);

CREATE INDEX IF NOT EXISTS i_user_campaign ON "send_logs" (user_id, campaign_id);
CREATE INDEX IF NOT EXISTS i_uuid ON "send_logs" (uuid);

CREATE TABLE IF NOT EXISTS "sends" (
  "id"                 integer primary key autoincrement,
  "user_id"            integer,
  "campaign_id"        integer,
  "message_id"         varchar(191) NOT NULL,
  "source"             varchar(191),
  "sending_account_id" varchar(191),
  "destination"        varchar(191),
  "created_at"         datetime
);

-- +migrate Down

DROP TABLE "users";
DROP TABLE "sessions";
DROP TABLE "campaigns";
DROP TABLE "segments";
DROP TABLE "subscribers";
DROP TABLE "subscribers_segments";
DROP TABLE "bounces";
DROP TABLE "sends";
DROP TABLE "clicks";
DROP TABLE "complaints";
DROP TABLE "deliveries";
DROP TABLE "ses_keys";
DROP TABLE "opens";