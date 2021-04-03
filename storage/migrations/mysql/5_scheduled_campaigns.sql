-- +migrate Up

CREATE TABLE IF NOT EXISTS `scheduled_campaigns`
(
    `id`            varbinary(27)       primary key,
    `campaign_id`   integer unsigned    NOT NULL,
    `scheduled_at`  datetime(6)         NOT NULL,
    `description`   varchar(191)        NOT NULL,
    `created_at`    datetime(6),
    `updated_at`    datetime(6),
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns (`id`)
    ) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `scheduled_campaigns`;
