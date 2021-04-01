-- +migrate Up

CREATE TABLE IF NOT EXISTS `campaign_failed_logs`
(
    `id`          varbinary(27) primary key,
    `user_id`     integer unsigned NOT NULL,
    `campaign_id` integer unsigned NOT NULL,
    `description` varchar(191)     NOT NULL,
    `created_at`  datetime(6),
    FOREIGN KEY (`user_id`) REFERENCES users (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `campaign_failed_logs`;
