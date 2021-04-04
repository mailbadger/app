-- +migrate Up

CREATE TABLE IF NOT EXISTS `campaign_schedules`
(
    `id`            varbinary(27)       primary key,
    `campaign_id`   integer unsigned UNIQUE NOT NULL,
    `scheduled_at`  datetime(6)         NOT NULL,
    `created_at`    datetime(6)         NOT NULL,
    `updated_at`    datetime(6)         NOT NULL,
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns (`id`)
    ) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `campaign_schedules`;
