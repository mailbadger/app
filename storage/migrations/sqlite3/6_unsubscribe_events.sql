-- +migrate Up

CREATE TABLE IF NOT EXISTS `unsubscribe_events`
(
    `id`         varchar(27) primary key,
    `user_id`    integer unsigned NOT NULL,
    `email`      varchar(191)     NOT NULL,
    `created_at` datetime         NOT NULL
);

-- +migrate Down

DROP TABLE `unsubscribe_events`;