-- +migrate Up

CREATE TABLE IF NOT EXISTS `unsubscribe_events`
(
    `id`         varbinary(27) primary key NOT NULL,
    `user_id`    integer unsigned          NOT NULL,
    `email`      varchar(191)              NOT NULL,
    `created_at` datetime(6)               NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users (`id`)
) CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `unsubscribe_events`;