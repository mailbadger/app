-- +migrate Up

CREATE TABLE IF NOT EXISTS `processes` (
    `id`                INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `name`              VARCHAR(30) UNIQUE NOT NULL,
    `last_processed_id` BIGINT UNSIGNED NOT NULL,
    `created_at`        datetime(6)            NOT NULL,
    `updated_at`        datetime(6)            NOT NULL,
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `accumulated_subscribes` (
    `id`                BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `user_id`           INTEGER UNSIGNED NOT NULL,
    `total_subscribers` BIGINT UNSIGNED NOT NULL,
    `created`           BIGINT UNSIGNED NOT NULL,
    `unsubscribed`      BIGINT UNSIGNED NOT NULL,
    `deleted`           BIGINT UNSIGNED NOT NULL,
    `date`              date         NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `processes`;
DROP TABLE `accumulated_subscribes`;
