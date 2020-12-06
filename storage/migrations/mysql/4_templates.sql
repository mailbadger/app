-- +migrate Up

CREATE TABLE IF NOT EXISTS `templates`
(
    `id`         integer unsigned primary key AUTO_INCREMENT NOT NULL,
    `name`       varchar(191)                                NOT NULL,
    `subject`    varchar(191)                                NOT NULL,
    `text_part`  text,
    `created_at` datetime(6)                                 NOT NULL,
    `updated_at` datetime(6)                                 NOT NULL,
) CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `templates`;