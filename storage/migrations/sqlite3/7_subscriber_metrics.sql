-- +migrate Up

CREATE TABLE IF NOT EXISTS `jobs` (
    `id`                  integer primary key autoincrement,
    `name`                varchar(30) unique,
    `last_processed_date` date,
    `status`              varchar(30),
    `created_at`          datetime,
    `updated_at`          datetime
);

INSERT INTO `jobs` (
    `name`,
    `last_processed_date`,
    `status`,
    `created_at`,
    `updated_at`
) VALUES ("subscriber_metrics",date('now'), "idle", datetime('now'), datetime('now'));

CREATE TABLE IF NOT EXISTS `subscriber_metrics` (
    `id`           integer unsigned unique,
    `user_id`      integer unsigned,
    `created`      integer unsigned,
    `unsubscribed` integer unsigned,
    `deleted`      integer unsigned,
    `date`         date,
    primary key (`user_id`, `date`),
    foreign key (`user_id`) references users (`id`)
);

-- +migrate Down

DROP TABLE `processes`;
DROP TABLE `accumulated_subscribes`;
