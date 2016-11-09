
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS "users" (
  "id"       integer primary key autoincrement,
  "username" varchar(255) NOT NULL UNIQUE,
  "password" varchar(255),
  "api_key"  varchar(255) NOT NULL UNIQUE,
  "auth_key" varchar(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "templates" (
  "id"         integer primary key autoincrement,
  "user_id"    integer,
  "name"       varchar(255),
  "content"    text,
  "created_at" datetime,
  "updated_at" datetime,
  UNIQUE("user_id", "name")
 );

CREATE INDEX IF NOT EXISTS i_user ON "templates" (user_id);

CREATE TABLE IF NOT EXISTS "campaigns" (
  "id"           integer primary key autoincrement,
  "user_id"      integer,
  "name"         varchar(255) NOT NULL,
  "subject"      varchar(255) NOT NULL,
  "template_id"  integer,
  "status"       varchar(255),
  "created_at"   datetime,
  "updated_at"   datetime,
  "scheduled_at" datetime,
  "completed_at" datetime
);

CREATE INDEX IF NOT EXISTS i_user ON "campaigns" (user_id);

CREATE TABLE IF NOT EXISTS "subscribers" (
  "id"         integer primary key autoincrement,
  "user_id"    integer,
  "name"       varchar(255) NOT NULL,
  "email"      varchar(255) NOT NULL,
  "created_at" datetime,
  "updated_at" datetime,
  UNIQUE("user_id", "email")
);

CREATE TABLE IF NOT EXISTS "lists" (
  "id"         integer primary key autoincrement,
  "user_id"    integer,
  "name"       varchar(255),
  "created_at" datetime,
  "updated_at" datetime
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
  "key"           varchar(255),
  "value"         varchar(255),
  "created_at"    datetime,
  "updated_at"    datetime
);

CREATE INDEX IF NOT EXISTS i_subscriber ON "subscriber_metadata" (subscriber_id);

CREATE TABLE IF NOT EXISTS "sent_emails" (
  "id"          integer primary key autoincrement,
  "campaign_id" integer,
  "user_id"     integer,
  "token"       varchar(255),
  "status"      varchar(255) NOT NULL,
  "ip"          varchar(255),
  "latitude"    real,
  "longitude"   real,
  "opens"       integer,
  "created_at"  datetime,
  "updated_at"  datetime
);

CREATE INDEX IF NOT EXISTS i_user     ON "sent_emails" (user_id);
CREATE INDEX IF NOT EXISTS i_campaign ON "sent_emails" (campaign_id);

CREATE TABLE IF NOT EXISTS "bounces" (
  "id"         integer primary key autoincrement,
  "recipient"  varchar(255),
  "sender"     varchar(255),
  "type"       varchar(255),
  "sub_type"   varchar(255),
  "action"     varchar(255),
  "created_at" datetime,
  "updated_at" datetime
);

CREATE TABLE IF NOT EXISTS "events" (
  "id"            integer primary key autoincrement,
  "campaign_id"   integer,
  "subscriber_id" integer,
  "message"       varchar(255),
  "created_at"    datetime,
  "updated_at"    datetime
);

CREATE INDEX IF NOT EXISTS i_campaign   ON "events" (campaign_id);
CREATE INDEX IF NOT EXISTS i_subscriber ON "events" (subscriber_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
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
DROP TABLE "events";
