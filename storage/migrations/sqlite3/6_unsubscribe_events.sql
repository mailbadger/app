-- +migrate Up

CREATE TABLE IF NOT EXISTS `unsubscribe_events`
(
    `id`         integer primary key autoincrement,
    `email`      varchar(191) NOT NULL,
    `created_at` datetime  NOT NULL
);

-- +migrate Down

DROP TABLE `unsubscribe_events`;