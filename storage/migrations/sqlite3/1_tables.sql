
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
  "updated_at" datetime
 );

CREATE TABLE IF NOT EXISTS "campaigns" (
  "id"           integer primary key autoincrement,
  "user_id"      integer,
  "name"         varchar(255) NOT NULL,
  "subject"      varchar(255) NOT NULL,
  "template_id"  integer,
  "status"       varchar(255),
  "created_at"   datetime,
  "scheduled_at" datetime,
  "completed_at" datetime
);

CREATE TABLE IF NOT EXISTS "subscribers" (
  "id"         integer primary key autoincrement,
  "name"       varchar(255),
  "email"      varchar(255),
  "created_at" datetime,
  "updated_at" datetime
);

CREATE TABLE IF NOT EXISTS "lists" (
  "id"                integer primary key autoincrement,
  "user_id"           integer,
  "name"              varchar(255),
  "total_subscribers" integer,
  "created_at"        datetime,
  "updated_at"        datetime
);
CREATE TABLE IF NOT EXISTS "subscribers_lists" (
  "list_id"       integer,
  "subscriber_id" integer
);

CREATE TABLE IF NOT EXISTS "fields" (
  "id"         integer primary key autoincrement,
  "name"       varchar(255),
  "list_id"    integer,
  "created_at" datetime,
  "updated_at" datetime
);

CREATE TABLE IF NOT EXISTS "subscribers_fields" (
  "id"            integer primary key autoincrement,
  "field_id"      integer,
  "subscriber_id" integer,
  "value"         varchar(255),
  "created_at"    datetime,
  "updated_at"    datetime
);

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
-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "users";
DROP TABLE "templates";
DROP TABLE "campaigns";
DROP TABLE "subscribers";
DROP TABLE "lists";
DROP TABLE "subscribers_lists";
DROP TABLE "fields";
DROP TABLE "subscribers_fields";
DROP TABLE "sent_emails";
DROP TABLE "bounces";
