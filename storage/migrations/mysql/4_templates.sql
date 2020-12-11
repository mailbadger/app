-- +migrate Up

CREATE TABLE IF NOT EXISTS `templates`
(
    `id`           integer unsigned primary key AUTO_INCREMENT NOT NULL,
    `user_id`      integer unsigned                            NOT NULL,
    `name`         varchar(191)                                NOT NULL,
    `subject_part` varchar(191)                                NOT NULL,
    `text_part`    text,
    `created_at`   datetime(6)                                 NOT NULL,
    `updated_at`   datetime(6)                                 NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users (`id`)
) CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `templates`;