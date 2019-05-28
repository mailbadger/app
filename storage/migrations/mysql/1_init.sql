-- +migrate Up

CREATE TABLE IF NOT EXISTS `users` (
  `id`         integer primary key AUTO_INCREMENT NOT NULL,
  `uuid`       varchar(36) NOT NULL UNIQUE,
  `username`   varchar(191) NOT NULL UNIQUE,
  `password`   varchar(191),
  `source`     varchar(191),
  `active`     integer,
  `verified`   integer,
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE IF NOT EXISTS `ses_keys` (
  `id`         integer primary key AUTO_INCREMENT NOT NULL,
  `user_id`    integer,
  `access_key` varchar(191) NOT NULL,
  `secret_key` varchar(191) NOT NULL,
  `region`     varchar(30) NOT NULL,
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE IF NOT EXISTS `campaigns` (
  `id`            integer primary key AUTO_INCREMENT NOT NULL,
  `user_id`       integer,
  `name`          varchar(191) NOT NULL,
  `template_name` varchar(191) NOT NULL,
  `status`        varchar(191),
  `created_at`    datetime,
  `updated_at`    datetime,
  `scheduled_at`  datetime DEFAULT NULL,
  `completed_at`  datetime DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS `subscribers` (
  `id`          integer primary key AUTO_INCREMENT NOT NULL,
  `user_id`     integer,
  `name`        varchar(191) NOT NULL,
  `email`       varchar(191) NOT NULL,
  `blacklisted` integer,
  `active`      integer,
  `created_at`  datetime,
  `updated_at`  datetime
);

CREATE TABLE IF NOT EXISTS `lists` (
  `id`          integer primary key AUTO_INCREMENT NOT NULL,
  `user_id`   integer,
  `name`      varchar(191),
  `created_at`  datetime,
  `updated_at`  datetime
);

CREATE TABLE IF NOT EXISTS `subscribers_lists` (
  `list_id`       integer NOT NULL,
  `subscriber_id` integer NOT NULL,
  PRIMARY KEY (`list_id`, `subscriber_id`)
);

CREATE TABLE IF NOT EXISTS `subscriber_metadata` (
  `id`            integer primary key AUTO_INCREMENT NOT NULL,
  `subscriber_id` integer,
  `key`         varchar(191),
  `value`       varchar(191),
  `created_at`    datetime,
  `updated_at`    datetime
);

CREATE TABLE IF NOT EXISTS `bounces` (
  `id`              integer primary key AUTO_INCREMENT NOT NULL,
  `campaign_id`     integer,
  `user_id`         integer,
  `recipient`       varchar(191),
  `type`            varchar(30),
  `sub_type`        varchar(30),
  `action`          varchar(191),
  `status`          varchar(191),
  `diagnostic_code` varchar(191),
  `feedback_id`     varchar(191),
  `created_at`      datetime
);

CREATE TABLE IF NOT EXISTS `complaints` (
  `id`          integer primary key AUTO_INCREMENT NOT NULL,
  `campaign_id` integer,
  `user_id`     integer,
  `recipient`   varchar(191),
  `type`        varchar(30),
  `user_agent`  varchar(191),
  `feedback_id` varchar(191),
  `created_at`  datetime
);

CREATE TABLE IF NOT EXISTS `clicks` (
  `id`          integer primary key AUTO_INCREMENT NOT NULL,
  `campaign_id` integer,
  `user_id`     integer,
  `ip_address`  varchar(50),
  `user_agent`  varchar(191),
  `link`        varchar(191),
  `created_at`  datetime
);

CREATE TABLE IF NOT EXISTS `opens` (
  `id`          integer primary key AUTO_INCREMENT NOT NULL,
  `campaign_id` integer,
  `user_id`     integer,
  `ip_address`  varchar(50),
  `user_agent`  varchar(191),
  `created_at`  datetime
);

CREATE TABLE IF NOT EXISTS `deliveries` (
  `id`                     integer primary key AUTO_INCREMENT NOT NULL,
  `campaign_id`            integer,
  `user_id`                integer,
  `recipient`              varchar(191),
  `processing_time_millis` integer,
  `smtp_response`          varchar(191),
  `reporting_mta`          varchar(191),
  `remote_mta_ip`          varchar(50),
  `created_at`             datetime
);

CREATE TABLE IF NOT EXISTS `send_bulk_logs` (
  `id`          integer primary key AUTO_INCREMENT NOT NULL,
  `uuid`        varchar(36) NOT NULL,
  `user_id`     integer,
  `campaign_id` integer,
  `message_id`  varchar(191),
  `status`      varchar(191) NOT NULL,
  `created_at`  datetime
);

CREATE TABLE IF NOT EXISTS `sends` (
  `id`                 integer primary key AUTO_INCREMENT NOT NULL,
  `user_id`            integer,
  `campaign_id`        integer,
  `message_id`         varchar(191) NOT NULL,
  `source`             varchar(191),
  `source_arn`         varchar(191),
  `source_ip`          varchar(191),
  `sending_account_id` varchar(191),
  `destination`        varchar(191),
  `created_at`         datetime
);

-- +migrate Down

DROP TABLE `users`;
DROP TABLE `campaigns`;
DROP TABLE `lists`;
DROP TABLE `subscribers`;
DROP TABLE `subscribers_lists`;
DROP TABLE `subscriber_metadata`;
DROP TABLE `bounces`;
DROP TABLE `send_bulk_logs`;
DROP TABLE `sends`;
DROP TABLE `clicks`;
DROP TABLE `complaints`;
DROP TABLE `deliveries`;
DROP TABLE `ses_keys`;
DROP TABLE `opens`;