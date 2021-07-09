-- +migrate Up

CREATE TABLE IF NOT EXISTS `jobs` (
    `id`                  INTEGER UNSIGNED PRIMARY KEY,
    `name`                VARCHAR(30) UNIQUE,
    `last_processed_date` DATE,
    `status`              VARCHAR(30),
    `created_at`          DATETIME(6),
    `updated_at`          DATETIME(6)
);

INSERT INTO `jobs` (
    `name`,
    `last_processed_date`,
    `status`,
    `created_at`,
    `updated_at`
) VALUES ("subscriber_metrics", "2006-01-02", "idle", datetime('now'), datetime('now'));

CREATE TABLE IF NOT EXISTS `subscriber_metrics` (
    `id`           integer UNSIGNED UNIQUE,
    `user_id`      INTEGER UNSIGNED,
    `created`      INTEGER UNSIGNED,
    `unsubscribed` INTEGER UNSIGNED,
    `deleted`      INTEGER UNSIGNED,
    `date`         DATE,
    PRIMARY KEY (`user_id`, `date`),
    FOREIGN KEY (`user_id`) REFERENCES users (`id`)
);

-- +migrate Down

DROP TABLE `processes`;
DROP TABLE `accumulated_subscribes`;
