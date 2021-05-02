-- +migrate Up

CREATE TABLE IF NOT EXISTS `reports` (
    `id`         integer unsigned primary key AUTO_INCREMENT NOT NULL,
    `user_id`    integer unsigned                            NOT NULL,
    `resource`   varchar(191)                                NOT NULL,
    `file_name`  varchar(191)                                NOT NULL,
    `type`       varchar(191)                                NOT NULL,
    `status`     varchar(191)                                NOT NULL,
    `note`       varchar(191),
    `created_at` datetime(6)                                 NOT NULL,
    `updated_at` datetime(6)                                 NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users (`id`),
    INDEX idx_user_id_resource (`user_id`, `resource`, `file_name`, `type`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `reports`;