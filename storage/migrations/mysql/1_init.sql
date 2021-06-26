-- +migrate Up

CREATE TABLE IF NOT EXISTS `boundaries` (
    `id`                         INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `type`                       VARCHAR(191) NOT NULL,
    `stats_retention`            INTEGER NOT NULL,
    `subscribers_limit`          INTEGER NOT NULL,
    `campaigns_limit`            INTEGER NOT NULL,
    `templates_limit`            INTEGER NOT NULL,
    `groups_limit`               INTEGER NOT NULL,
    `schedule_campaigns_enabled` TINYINT(1) NOT NULL,
    `saml_enabled`               TINYINT(1) NOT NULL,
    `team_members_limit`         INTEGER NOT NULL,
    `created_at`                 DATETIME(6) NOT NULL,
    `updated_at`                 DATETIME(6) NOT NULL
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

INSERT INTO `boundaries` (
  `type`, 
  `stats_retention`, 
  `subscribers_limit`, 
  `campaigns_limit`, 
  `templates_limit`, 
  `groups_limit`, 
  `schedule_campaigns_enabled`, 
  `saml_enabled`, 
  `team_members_limit`, 
  `created_at`, 
  `updated_at`
  ) VALUES ("nolimit", 0, 0, 0, 0, 0, 1, 1, 0, NOW(), NOW());

INSERT INTO `boundaries` (
  `type`, 
  `stats_retention`, 
  `subscribers_limit`, 
  `campaigns_limit`, 
  `templates_limit`, 
  `groups_limit`, 
  `schedule_campaigns_enabled`, 
  `saml_enabled`, 
  `team_members_limit`, 
  `created_at`, 
  `updated_at`
  ) VALUES ("free", 0, 0, 3, 0, 0, 0, 0, 0, NOW(), NOW());

CREATE TABLE IF NOT EXISTS `users` (
    `id`          INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `uuid`        VARCHAR(36) NOT NULL UNIQUE,
    `username`    VARCHAR(191) NOT NULL UNIQUE,
    `password`    VARCHAR(191) NOT NULL,
    `source`      VARCHAR(191) NOT NULL,
    `active`      INTEGER NOT NULL,
    `verified`    INTEGER NOT NULL,
    `boundary_id` INTEGER UNSIGNED NOT NULL,
    `created_at`  DATETIME(6) NOT NULL,
    `updated_at`  DATETIME(6) NOT NULL,
    FOREIGN KEY (`boundary_id`) REFERENCES boundaries(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `roles` (
    `id`    INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `name`  VARCHAR(100) NOT NULL
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

INSERT INTO `roles` (`name`) VALUES ("admin");
INSERT INTO `roles` (`name`) VALUES ("billing");

CREATE TABLE IF NOT EXISTS `users_roles` (
    `user_id` INTEGER UNSIGNED NOT NULL,
    `role_id` INTEGER UNSIGNED NOT NULL,
    PRIMARY KEY (`user_id`, `role_id`),
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`role_id`) REFERENCES roles(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `sessions` (
    `id`         INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `user_id`    INTEGER UNSIGNED NOT NULL,
    `session_id` VARCHAR(191) NOT NULL UNIQUE,
    `created_at` DATETIME(6) NOT NULL,
    `updated_at` DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `api_keys` (
    `id`         INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `user_id`    INTEGER UNSIGNED NOT NULL,
    `secret_key` VARCHAR(191) NOT NULL UNIQUE,
    `active`     TINYINT(1) NOT NULL,
    `created_at` DATETIME(6) NOT NULL,
    `updated_at` DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `ses_keys` (
    `id`         INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `user_id`    INTEGER UNSIGNED NOT NULL UNIQUE,
    `access_key` VARCHAR(191) NOT NULL,
    `secret_key` VARCHAR(191) NOT NULL,
    `region`     VARCHAR(30) NOT NULL,
    `created_at` DATETIME(6) NOT NULL,
    `updated_at` DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `templates` (
    `id`           integer unsigned primary key AUTO_INCREMENT NOT NULL,
    `user_id`      integer unsigned                            NOT NULL,
    `name`         varchar(191)                                NOT NULL,
    `subject_part` varchar(191)                                NOT NULL,
    `text_part`    text,
    `created_at`   datetime(6)                                 NOT NULL,
    `updated_at`   datetime(6)                                 NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `campaigns` (
    `id`            INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `user_id`       INTEGER UNSIGNED NOT NULL,
    `name`          VARCHAR(191) NOT NULL,
    `template_id`   INTEGER UNSIGNED,
    `event_id`      VARBINARY(27) DEFAULT NULL,
    `status`        VARCHAR(191) NOT NULL,
    `created_at`    DATETIME(6) NOT NULL,
    `updated_at`    DATETIME(6) NOT NULL,
    `completed_at`  DATETIME(6) DEFAULT NULL,
    `deleted_at`    DATETIME(6) DEFAULT NULL,
    `started_at`    DATETIME(6) DEFAULT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`template_id`) REFERENCES templates(`id`) ON DELETE SET NULL ON UPDATE CASCADE,
    INDEX idx_id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `subscribers` (
    `id`          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `user_id`     INTEGER UNSIGNED NOT NULL,
    `name`        VARCHAR(191) DEFAULT NULL,
    `email`       VARCHAR(191) NOT NULL,
    `metadata`    JSON,
    `blacklisted` TINYINT(1) NOT NULL,
    `active`      TINYINT(1) NOT NULL,
    `created_at`  DATETIME(6) NOT NULL,
    `updated_at`  DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    INDEX idx_id_created_at (`id`, `created_at`),
    INDEX idx_user_blacklist_active (`user_id`, `blacklisted`, `active`),
    UNIQUE(`user_id`, `email`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `segments` (
    `id`          INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `user_id`     INTEGER UNSIGNED NOT NULL,
    `name`        VARCHAR(191) NOT NULL,
    `created_at`  DATETIME(6) NOT NULL,
    `updated_at`  DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    INDEX idx_id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `subscribers_segments` (
    `segment_id`    INTEGER UNSIGNED NOT NULL,
    `subscriber_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`segment_id`, `subscriber_id`),
    FOREIGN KEY (`segment_id`) REFERENCES segments(`id`),
    FOREIGN KEY (`subscriber_id`) REFERENCES subscribers(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `bounces` (
    `id`              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `campaign_id`     INTEGER UNSIGNED NOT NULL,
    `user_id`         INTEGER UNSIGNED NOT NULL,
    `recipient`       VARCHAR(191) NOT NULL,
    `type`            VARCHAR(30) NOT NULL,
    `sub_type`        VARCHAR(30) NOT NULL,
    `action`          VARCHAR(191) NOT NULL,
    `status`          VARCHAR(191) NOT NULL,
    `diagnostic_code` VARCHAR(191) NOT NULL,
    `feedback_id`     VARCHAR(191) NOT NULL,
    `created_at`      DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
    INDEX idx_id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `complaints` (
    `id`          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `campaign_id` INTEGER UNSIGNED NOT NULL,
    `user_id`     INTEGER UNSIGNED NOT NULL,
    `recipient`   VARCHAR(191) NOT NULL,
    `type`        VARCHAR(30) NOT NULL,
    `user_agent`  VARCHAR(191) NOT NULL,
    `feedback_id` VARCHAR(191) NOT NULL,
    `created_at`  DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
    INDEX idx_id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `clicks` (
    `id`          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `campaign_id` INTEGER UNSIGNED NOT NULL,
    `user_id`     INTEGER UNSIGNED NOT NULL,
    `ip_address`  VARCHAR(50) NOT NULL,
    `recipient`   VARCHAR(191) NOT NULL,
    `user_agent`  VARCHAR(191) NOT NULL,
    `link`        VARCHAR(191) NOT NULL,
    `created_at`  DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
    INDEX idx_id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `opens` (
    `id`          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `campaign_id` INTEGER UNSIGNED NOT NULL,
    `user_id`     INTEGER UNSIGNED NOT NULL,
    `recipient`   VARCHAR(191) NOT NULL,
    `ip_address`  VARCHAR(50) NOT NULL,
    `user_agent`  VARCHAR(191) NOT NULL,
    `created_at`  DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
    INDEX idx_id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `deliveries` (
    `id`                     BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `campaign_id`            INTEGER UNSIGNED NOT NULL,
    `user_id`                INTEGER UNSIGNED NOT NULL,
    `recipient`              VARCHAR(191) NOT NULL,
    `processing_time_millis` INTEGER NOT NULL,
    `smtp_response`          VARCHAR(191) NOT NULL,
    `reporting_mta`          VARCHAR(191) NOT NULL,
    `remote_mta_ip`          VARCHAR(50) NOT NULL,
    `created_at`             DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
    INDEX idx_id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `send_logs` (
    `id`              varbinary(27) primary key NOT NULL,
    `user_id`         integer unsigned NOT NULL,
    `event_id`        varbinary(27) NOT NULL,
    `subscriber_id`   integer unsigned NOT NULL,
    `campaign_id`     integer unsigned NOT NULL,
    `status`          varchar(191) NOT NULL,
    `message_id`      varchar(191),
    `description`     varchar(191) NOT NULL,
    `created_at`      datetime(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
    INDEX idx_id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `sends` (
    `id`                 BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `user_id`            INTEGER UNSIGNED NOT NULL,
    `campaign_id`        INTEGER UNSIGNED NOT NULL,
    `message_id`         VARCHAR(191) NOT NULL,
    `source`             VARCHAR(191) NOT NULL,
    `sending_account_id` VARCHAR(191) NOT NULL,
    `destination`        VARCHAR(191) NOT NULL,
    `created_at`         DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns(`id`),
    INDEX idx_id_created_at (`id`, `created_at`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `subscriber_events` (
    `id`               VARBINARY(27) PRIMARY KEY NOT NULL,
    `user_id`          INTEGER UNSIGNED NOT NULL,
    `subscriber_email` VARCHAR(191) NOT NULL,
    `event_type`       VARCHAR(50) NOT NULL,
    `created_at`       DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users (`id`),
    INDEX idx_user_id_created_at (`user_id`, `created_at`)
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
DROP TABLE `templates`;
DROP TABLE `sessions`;
DROP TABLE `boundaries`;
DROP TABLE `users_roles`;
DROP TABLE `roles`;
DROP TABLE `users`;