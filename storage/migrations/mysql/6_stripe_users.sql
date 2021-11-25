-- +migrate Up

CREATE TABLE IF NOT EXISTS `stripe_customers` (
    `id`                   INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    `user_id`              INTEGER UNSIGNED NOT NULL,
    `customer_id`          VARCHAR(191) UNIQUE NOT NULL,
    `created_at`           DATETIME(6) NOT NULL,
    `updated_at`           DATETIME(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `stripe_customers`;