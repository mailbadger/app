-- +migrate Up

CREATE TABLE IF NOT EXISTS `subscriber_metrics` (
    `user_id`      integer unsigned,
    `created`      integer unsigned,
    `unsubscribed` integer unsigned,
    `deleted`      integer unsigned,
    `datetime`     datetime,
    primary key (`user_id`, `datetime`),
    foreign key (`user_id`) references users (`id`)
);

-- +migrate Down

DROP TABLE `subscriber_metrics`;
