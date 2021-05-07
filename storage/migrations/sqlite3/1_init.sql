-- +migrate Up

CREATE TABLE IF NOT EXISTS "boundaries" (
    "id"                         integer primary key autoincrement,
    "type"                       varchar(191) not null,
    "stats_retention"            integer not null,
    "subscribers_limit"          integer not null,
    "campaigns_limit"            integer not null,
    "templates_limit"            integer not null,
    "groups_limit"               integer not null,
    "schedule_campaigns_enabled" integer not null,
    "saml_enabled"               integer not null,
    "team_members_limit"         integer not null,
    "created_at"                 datetime not null,
    "updated_at"                 datetime not null
);

INSERT INTO "boundaries" (
  "type", 
  "stats_retention", 
  "subscribers_limit", 
  "campaigns_limit", 
  "templates_limit", 
  "groups_limit", 
  "schedule_campaigns_enabled", 
  "saml_enabled", 
  "team_members_limit", 
  "created_at", 
  "updated_at"
) VALUES ("nolimit", 0, 0, 0, 0, 0, 1, 1, 0, datetime('now'), datetime('now'));

INSERT INTO "boundaries" (
  "type", 
  "stats_retention", 
  "subscribers_limit", 
  "campaigns_limit", 
  "templates_limit", 
  "groups_limit", 
  "schedule_campaigns_enabled", 
  "saml_enabled", 
  "team_members_limit", 
  "created_at", 
  "updated_at"
) VALUES ("free", 0, 0, 3, 0, 0, 0, 0, 0, datetime('now'), datetime('now'));

INSERT INTO "boundaries" (
  "type", 
  "stats_retention", 
  "subscribers_limit", 
  "campaigns_limit", 
  "templates_limit", 
  "groups_limit", 
  "schedule_campaigns_enabled", 
  "saml_enabled", 
  "team_members_limit", 
  "created_at", 
  "updated_at"
) VALUES ("db_test", 0, 0, 2, 0, 0, 1, 1, 0, datetime('now'), datetime('now'));

CREATE TABLE IF NOT EXISTS "users" (
    "id"          integer primary key autoincrement,
    "uuid"        varchar(36) not null UNIQUE,
    "username"    varchar(191) not null UNIQUE,
    "password"    varchar(191),
    "source"      varchar(191) not null,
    "active"      integer,
    "verified"    integer,
    "boundary_id" integer,
    "created_at"  datetime,
    "updated_at"  datetime,
    foreign key ("boundary_id") references boundaries("id")
);

CREATE TABLE IF NOT EXISTS "roles" (
    "id"    integer primary key autoincrement,
    "name"  varchar(100)
);

INSERT INTO "roles" ("name") VALUES ("admin");

CREATE TABLE IF NOT EXISTS "users_roles" (
    "user_id" integer,
    "role_id" integer,
    primary key ("user_id", "role_id"),
    foreign key ("user_id") references users("id"),
    foreign key ("role_id") references roles("id")
);

CREATE INDEX IF NOT EXISTS idx_user ON "users_roles" (user_id);
CREATE INDEX IF NOT EXISTS idx_role ON "users_roles" (role_id);

CREATE TABLE IF NOT EXISTS "sessions" (
    "id"         integer primary key autoincrement,
    "user_id"    integer not null,
    "session_id" varchar(191) not null,
    "created_at" datetime not null,
    "updated_at" datetime not null,
    UNIQUE("session_id"),
    foreign key ("user_id") references users("id")
);

CREATE TABLE IF NOT EXISTS "api_keys" (
    "id"         integer primary key autoincrement,
    "user_id"    integer not null,
    "secret_key" varchar(191) not null,
    "active"     integer not null,
    "created_at" datetime not null,
    "updated_at" datetime not null,
    UNIQUE("secret_key"),
    foreign key ("user_id") references users("id")
);

CREATE TABLE IF NOT EXISTS "templates" (
    "id"           integer primary key autoincrement,
    "user_id"      integer unsigned NOT NULL,
    "name"         varchar(191)     NOT NULL,
    "subject_part" varchar(191)     NOT NULL,
    "text_part"    text,
    "created_at"   datetime,
    "updated_at"   datetime,
    foreign key ("user_id") references users("id")
);

CREATE TABLE IF NOT EXISTS "ses_keys" (
    "id"         integer primary key autoincrement,
    "user_id"    integer,
    "access_key" varchar(191) not null,
    "secret_key" varchar(191) not null,
    "region"     varchar(30) not null,
    "created_at" datetime,
    "updated_at" datetime,
    UNIQUE("user_id"),
    foreign key ("user_id") references users("id")
);

CREATE TABLE IF NOT EXISTS "campaigns" (
    "id"            integer primary key autoincrement,
    "user_id"       integer,
    "name"          varchar(191) not null,
    "template_id"   integer,
    "event_id"      varchar(27),
    "status"        varchar(191),
    "created_at"    datetime,
    "updated_at"    datetime,
    "completed_at"  datetime DEFAULT NULL,
    "deleted_at"    datetime DEFAULT NULL,
    "started_at"    datetime DEFAULT NULL,
    foreign key ("user_id") references users("id"),
    foreign key ("template_id") references templates("id")
);

CREATE INDEX IF NOT EXISTS idx_user ON "campaigns" (user_id);
CREATE INDEX IF NOT EXISTS idx_id_created_at ON "campaigns" (id, created_at);

CREATE TABLE IF NOT EXISTS "subscribers" (
    "id"          integer primary key autoincrement,
    "user_id"     integer,
    "name"        varchar(191),
    "email"       varchar(191) not null,
    "metadata"    json,
    "blacklisted" integer,
    "active"      integer,
    "created_at"  datetime,
    "updated_at"  datetime,
    UNIQUE("user_id", "email"),
    foreign key ("user_id") references users("id")
);

CREATE INDEX IF NOT EXISTS idx_user ON "subscribers" (user_id);
CREATE INDEX IF NOT EXISTS idx_id_created_at ON "subscribers" (id, created_at);
CREATE INDEX IF NOT EXISTS idx_user_blacklist_active ON "subscribers" (user_id, blacklisted, active);

CREATE TABLE IF NOT EXISTS "segments" (
    "id"          integer primary key autoincrement,
    "user_id"     integer,
    "name"        varchar(191),
    "created_at"  datetime,
    "updated_at"  datetime,
    foreign key ("user_id") references users("id")
);

CREATE INDEX IF NOT EXISTS idx_id_created_at ON "segments" (id, created_at);
CREATE INDEX IF NOT EXISTS idx_user ON "segments" (user_id);

CREATE TABLE IF NOT EXISTS "subscribers_segments" (
    "segment_id"    integer,
    "subscriber_id" integer,
    primary key ("segment_id", "subscriber_id"),
    foreign key ("segment_id") references segments("id"),
    foreign key ("subscriber_id") references subscribers("id")
);

CREATE INDEX IF NOT EXISTS idx_segment    ON "subscribers_segments" (segment_id);
CREATE INDEX IF NOT EXISTS idx_subscriber ON "subscribers_segments" (subscriber_id);

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
    "created_at"      datetime,
    foreign key ("user_id") references users("id"),
    foreign key ("campaign_id") references campaigns("id")
);

CREATE INDEX IF NOT EXISTS idx_id_created_at ON "bounces" (id, created_at);

CREATE TABLE IF NOT EXISTS "complaints" (
    "id"              integer primary key autoincrement,
    "campaign_id"     integer,
    "user_id"         integer,
    "recipient"       varchar(191),
    "type"            varchar(30),
    "user_agent"      varchar(191),
    "feedback_id"     varchar(191),
    "created_at"      datetime,
    foreign key ("user_id") references users("id"),
    foreign key ("campaign_id") references campaigns("id")
);

CREATE INDEX IF NOT EXISTS idx_id_created_at ON "complaints" (id, created_at);

CREATE TABLE IF NOT EXISTS "clicks" (
    "id"              integer primary key autoincrement,
    "campaign_id"     integer,
    "user_id"         integer,
    "recipient"       varchar(191),
    "ip_address"      varchar(50),
    "user_agent"      varchar(191),
    "link"            varchar(191),
    "created_at"      datetime,
    foreign key ("user_id") references users("id"),
    foreign key ("campaign_id") references campaigns("id")
);

CREATE INDEX IF NOT EXISTS idx_id_created_at ON "clicks" (id, created_at);
CREATE INDEX IF NOT EXISTS idx_link ON "clicks" (campaign_id, user_id, link);
CREATE INDEX IF NOT EXISTS idx_link_recipients ON "clicks" (campaign_id, user_id, recipient, link);

CREATE TABLE IF NOT EXISTS "opens" (
    "id"              integer primary key autoincrement,
    "campaign_id"     integer,
    "user_id"         integer,
    "recipient"       varchar(191),
    "ip_address"      varchar(50),
    "user_agent"      varchar(191),
    "created_at"      datetime,
    foreign key ("user_id") references users("id"),
    foreign key ("campaign_id") references campaigns("id")
);

CREATE INDEX IF NOT EXISTS idx_id_created_at ON "opens" (id, created_at);

CREATE TABLE IF NOT EXISTS "deliveries" (
    "id"                     integer primary key autoincrement,
    "campaign_id"            integer,
    "user_id"                integer,
    "recipient"              varchar(191),
    "processing_time_millis" integer,
    "smtp_response"          varchar(191),
    "reporting_mta"          varchar(191),
    "remote_mta_ip"          varchar(50),
    "created_at"             datetime,
    foreign key ("user_id") references users("id"),
    foreign key ("campaign_id") references campaigns("id")
);

CREATE INDEX IF NOT EXISTS idx_id_created_at ON "deliveries" (id, created_at);

CREATE TABLE IF NOT EXISTS "send_logs" (
  "id"            varchar(27) primary key,
  "user_id"       integer NOT NULL,
  "event_id"      varchar(27) NOT NULL,
  "campaign_id"   integer NOT NULL,
  "subscriber_id" integer NOT NULL,
  "status"        varchar(191) NOT NULL,
  "message_id"    varchar(191),
  "description"   varchar(191),
  "created_at"    datetime,
  foreign key ("user_id") references users("id"),
  foreign key ("campaign_id") references campaigns("id")
);

CREATE INDEX IF NOT EXISTS idx_id_created_at ON "send_logs" (id, created_at);

CREATE TABLE IF NOT EXISTS "sends" (
    "id"                 integer primary key autoincrement,
    "user_id"            integer,
    "campaign_id"        integer,
    "message_id"         varchar(191) not null,
    "source"             varchar(191),
    "sending_account_id" varchar(191),
    "destination"        varchar(191),
    "created_at"         datetime,
    foreign key ("user_id") references users("id"),
    foreign key ("campaign_id") references campaigns("id")
);

CREATE INDEX IF NOT EXISTS idx_id_created_at ON "sends" (id, created_at);

CREATE TABLE IF NOT EXISTS "subscriber_events" (
    "id"               VARBINARY(27) PRIMARY KEY NOT NULL,
    "user_id"          INTEGER UNSIGNED NOT NULL,
    "subscriber_email" VARCHAR(191) NOT NULL,
    "event_type"       VARCHAR(50) NOT NULL,
    "created_at"       DATETIME(6) NOT NULL,
    foreign key ("user_id") references users("id")
);

CREATE INDEX IF NOT EXISTS idx_id_created_at ON "subscriber_events" (id, created_at);
CREATE INDEX IF NOT EXISTS idx_event_type ON "subscriber_events" (event_type);

-- +migrate Down

DROP TABLE "sessions";
DROP TABLE "templates";
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
DROP TABLE "boundaries";
DROP TABLE "users_roles";
DROP TABLE "roles";
DROP TABLE "users";