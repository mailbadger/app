-- +migrate Up

CREATE TABLE IF NOT EXISTS `users` (
  `id`         integer unsigned primary key AUTO_INCREMENT NOT NULL,
  `uuid`       varchar(36) NOT NULL UNIQUE,
  `username`   varchar(191) NOT NULL UNIQUE,
  `password`   varchar(191) NOT NULL,
  `source`     varchar(191) NOT NULL,
  `active`     integer NOT NULL,
  `verified`   integer NOT NULL,
  `created_at` datetime(6) NOT NULL,
  `updated_at` datetime(6) NOT NULL
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `sessions` (
  `id` integer unsigned primary key AUTO_INCREMENT NOT NULL,
  `user_id`    integer unsigned NOT NULL,
  `session_id` varchar(191) NOT NULL UNIQUE,
  `created_at` datetime(6) NOT NULL,
  `updated_at` datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `api_keys` (
  `id`         integer unsigned primary key AUTO_INCREMENT NOT NULL,
  `user_id`    integer unsigned NOT NULL,
  `secret_key` varchar(191) NOT NULL UNIQUE,
  `active`     tinyint(1) NOT NULL,
  `created_at` datetime(6) NOT NULL,
  `updated_at` datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `ses_keys` (
  `id`         integer unsigned primary key AUTO_INCREMENT NOT NULL,
  `user_id`    integer unsigned NOT NULL UNIQUE,
  `access_key` varchar(191) NOT NULL,
  `secret_key` varchar(191) NOT NULL,
  `region`     varchar(30) NOT NULL,
  `created_at` datetime(6) NOT NULL,
  `updated_at` datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `campaigns` (
  `id`            integer unsigned primary key AUTO_INCREMENT NOT NULL,
  `user_id`       integer unsigned NOT NULL,
  `name`          varchar(191) NOT NULL,
  `template_id`   integer unsigned,
  `status`        varchar(191) NOT NULL,
  `created_at`    datetime(6) NOT NULL,
  `updated_at`    datetime(6) NOT NULL,
  `scheduled_at`  datetime(6) DEFAULT NULL,
  `completed_at`  datetime(6) DEFAULT NULL,
  `deleted_at`    datetime(6) DEFAULT NULL,
  `started_at`    datetime(6) DEFAULT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  FOREIGN KEY (`template_id`) REFERENCES templates(`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  INDEX id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `subscribers` (
  `id`          bigint unsigned primary key AUTO_INCREMENT NOT NULL,
  `user_id`     integer unsigned NOT NULL,
  `name`        varchar(191) DEFAULT NULL,
  `email`       varchar(191) NOT NULL,
  `metadata`    JSON,
  `blacklisted` tinyint(1) NOT NULL,
  `active`      tinyint(1) NOT NULL,
  `created_at`  datetime(6) NOT NULL,
  `updated_at`  datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  INDEX id_created_at (`id`, `created_at`),
  UNIQUE(`user_id`, `email`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `segments` (
  `id`          integer unsigned primary key AUTO_INCREMENT NOT NULL,
  `user_id`     integer unsigned NOT NULL,
  `name`        varchar(191) NOT NULL,
  `created_at`  datetime(6) NOT NULL,
  `updated_at`  datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  INDEX id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `subscribers_segments` (
  `segment_id`    integer unsigned NOT NULL,
  `subscriber_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`segment_id`, `subscriber_id`),
  FOREIGN KEY (`segment_id`) REFERENCES segments(`id`),
  FOREIGN KEY (`subscriber_id`) REFERENCES subscribers(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `bounces` (
  `id`              bigint unsigned primary key AUTO_INCREMENT NOT NULL,
  `campaign_id`     integer unsigned NOT NULL,
  `user_id`         integer unsigned NOT NULL,
  `recipient`       varchar(191) NOT NULL,
  `type`            varchar(30) NOT NULL,
  `sub_type`        varchar(30) NOT NULL,
  `action`          varchar(191) NOT NULL,
  `status`          varchar(191) NOT NULL,
  `diagnostic_code` varchar(191) NOT NULL,
  `feedback_id`     varchar(191) NOT NULL,
  `created_at`      datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
  INDEX id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `complaints` (
  `id`          bigint unsigned primary key AUTO_INCREMENT NOT NULL,
  `campaign_id` integer unsigned NOT NULL,
  `user_id`     integer unsigned NOT NULL,
  `recipient`   varchar(191) NOT NULL,
  `type`        varchar(30) NOT NULL,
  `user_agent`  varchar(191) NOT NULL,
  `feedback_id` varchar(191) NOT NULL,
  `created_at`  datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
  INDEX id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `clicks` (
  `id`          bigint unsigned primary key AUTO_INCREMENT NOT NULL,
  `campaign_id` integer unsigned NOT NULL,
  `user_id`     integer unsigned NOT NULL,
  `ip_address`  varchar(50) NOT NULL,
  `recipient`   varchar(191) NOT NULL,
  `user_agent`  varchar(191) NOT NULL,
  `link`        varchar(191) NOT NULL,
  `created_at`  datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
  INDEX id_created_at (`id`, `created_at`),
  INDEX link (`campaign_id`, `user_id`, `link`),
  INDEX link_recipients (`campaign_id`, `user_id`, `recipient`, `link`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `opens` (
  `id`          bigint unsigned primary key AUTO_INCREMENT NOT NULL,
  `campaign_id` integer unsigned NOT NULL,
  `user_id`     integer unsigned NOT NULL,
  `recipient`   varchar(191) NOT NULL,
  `ip_address`  varchar(50) NOT NULL,
  `user_agent`  varchar(191) NOT NULL,
  `created_at`  datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
  INDEX id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `deliveries` (
  `id`                     bigint unsigned primary key AUTO_INCREMENT NOT NULL,
  `campaign_id`            integer unsigned NOT NULL,
  `user_id`                integer unsigned NOT NULL,
  `recipient`              varchar(191) NOT NULL,
  `processing_time_millis` integer NOT NULL,
  `smtp_response`          varchar(191) NOT NULL,
  `reporting_mta`          varchar(191) NOT NULL,
  `remote_mta_ip`          varchar(50) NOT NULL,
  `created_at`             datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
  INDEX id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `send_logs` (
  `id`              bigint unsigned primary key AUTO_INCREMENT NOT NULL,
  `uuid`            varchar(36) NOT NULL,
  `user_id`         integer unsigned NOT NULL,
  `subscriber_id`   integer unsigned NOT NULL,
  `campaign_id`     integer unsigned NOT NULL,
  `status`          varchar(191) NOT NULL,
  `description`     varchar(191) NOT NULL,
  `created_at`      datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
  INDEX u_s_c_id (`user_id`, `uuid`,`campaign_id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `sends` (
  `id`                 bigint unsigned primary key AUTO_INCREMENT NOT NULL,
  `user_id`            integer unsigned NOT NULL,
  `campaign_id`        integer unsigned NOT NULL,
  `message_id`         varchar(191) NOT NULL,
  `source`             varchar(191) NOT NULL,
  `sending_account_id` varchar(191) NOT NULL,
  `destination`        varchar(191) NOT NULL,
  `created_at`         datetime(6) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
  INDEX id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `subscribers_segments`;
DROP TABLE `subscriber_metadata`;
DROP TABLE `segments`;
DROP TABLE `subscribers`;
DROP TABLE `bounces`;
DROP TABLE `sends`;
DROP TABLE `send_logs`;
DROP TABLE `clicks`;
DROP TABLE `complaints`;
DROP TABLE `deliveries`;
DROP TABLE `ses_keys`;
DROP TABLE `opens`;
DROP TABLE `campaigns`;
DROP TABLE `users`;
DROP TABLE `sessions`;