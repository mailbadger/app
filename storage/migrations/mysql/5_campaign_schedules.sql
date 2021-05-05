-- +migrate Up

CREATE TABLE IF NOT EXISTS `campaign_schedules` (
    `id`            varbinary(27)       primary key,
    `user_id`       integer unsigned        NOT NULL,
    `campaign_id`   integer unsigned UNIQUE NOT NULL,
    `scheduled_at`  datetime(6)         NOT NULL,
    `source`        varchar(191)        NOT NULL,
    `from_name`     varchar(60)         NOT NULL,
    `segment_ids`           JSON        NOT NULL,
    `default_template_data` JSON        NOT NULL,
    `created_at`    datetime(6)         NOT NULL,
    `updated_at`    datetime(6)         NOT NULL,
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns (`id`)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `campaign_schedules`;
