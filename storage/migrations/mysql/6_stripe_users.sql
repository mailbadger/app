-- +migrate Up

CREATE TABLE IF NOT EXISTS `stripe_users` (
    `id`             INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    `user_id`        INTEGER UNSIGNED NOT NULL,
    `stripe_user_id` VARCHAR(191) NOT NULL,
    `created_at`     DATETIME(6) NOT NULL,
    `updated_at`     DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `stripe_users`;