-- +migrate Up

CREATE TABLE IF NOT EXISTS `tokens` (
    `id`         bigint unsigned primary key AUTO_INCREMENT NOT NULL,
    `user_id`    integer unsigned NOT NULL,
    `token`      varchar(191) NOT NULL UNIQUE,
    `type`       varchar(191) NOT NULL,
    `expires_at` datetime(6) NOT NULL,
    `created_at` datetime(6) NOT NULL,
    `updated_at` datetime(6) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `tokens`;