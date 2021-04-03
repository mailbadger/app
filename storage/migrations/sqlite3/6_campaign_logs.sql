-- +migrate Up

CREATE TABLE IF NOT EXISTS `campaign_failed_logs`
(
    `id`          varchar(27) primary key,
    `user_id`     integer unsigned NOT NULL,
    `campaign_id` integer unsigned NOT NULL,
    `description` varchar(191)     NOT NULL,
    `created_at`  datetime
);

-- +migrate Down

DROP TABLE `campaign_failed_logs`;
