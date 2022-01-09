-- +migrate Up
CREATE TABLE IF NOT EXISTS `subscriber_metrics` (
    `user_id` INTEGER UNSIGNED NOT NULL,
    `created` INTEGER UNSIGNED NOT NULL,
    `unsubscribed` INTEGER UNSIGNED NOT NULL,
    `datetime` DATETIME NOT NULL,
    PRIMARY KEY (`user_id`, `datetime`),
    FOREIGN KEY (`user_id`) REFERENCES users (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE `subscriber_metrics`;
