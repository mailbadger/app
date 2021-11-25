-- +migrate Up

CREATE TABLE IF NOT EXISTS `stripe_subscriptions` (
    `id`              INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    `user_id`         INTEGER UNSIGNED NOT NULL,
    `stripe_id`       VARCHAR(191) UNIQUE NOT NULL,
    `status`          VARCHAR(191) NOT NULL,
    `trial_ends_at`   DATETIME(6),
    `ends_at`         DATETIME(6),
    `created_at`      DATETIME(6) NOT NULL,
    `updated_at`      DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `stripe_subscriptions`;