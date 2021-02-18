-- +migrate Up

CREATE TABLE IF NOT EXISTS `unsubscribed_subscriber`
(
    `id`         integer unsigned primary key NOT NULL,
    `email`      varchar(191)                 NOT NULL,
    `created_at` datetime(6)                  NOT NULL
);

-- +migrate Down

DROP TABLE `unsubscribed_subscriber`;