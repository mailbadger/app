-- +migrate Up

CREATE TABLE IF NOT EXISTS `jobs` (
    `id`                  INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    `name`                VARCHAR(30) UNIQUE NOT NULL,
    `last_processed_date` DATE(6) UNSIGNED NOT NULL,
    `status`              VARCHAR(30) NOT NULL,
    `created_at`          DATETIME(6) NOT NULL,
    `updated_at`          DATETIME(6) NOT NULL,
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

INSERT INTO `jobs` (
    `name`,
    `last_processed_id`,
    `created_at`,
    `updated_at`
) VALUES ("subscriber_metrics", 0, NOW(), NOW())

CREATE TABLE IF NOT EXISTS `subscriber_metrics` (
    `id`                BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
    `user_id`           INTEGER UNSIGNED NOT NULL,
    `created`           INTEGER UNSIGNED NOT NULL,
    `unsubscribed`      INTEGER UNSIGNED NOT NULL,
    `deleted`           INTEGER UNSIGNED NOT NULL,
    `date`              DATE NOT NULL,
    PRIMARY KEY (`user_id`, `date`),
    FOREIGN KEY (`user_id`) REFERENCES users (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `processes`;
DROP TABLE `accumulated_subscribes`;
