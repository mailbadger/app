-- +migrate Up

CREATE TABLE IF NOT EXISTS `scheduled_campaigns`
(
    `id`            varbinary(27)       primary key,
    `campaign_id`   integer unsigned    NOT NULL,
    `scheduled_at`  datetime(6)         NOT NULL,
    `created_at`    datetime(6),
    `updated_at`    datetime(6),
    FOREIGN KEY (`campaign_id`) REFERENCES campaigns (`id`)
    );

-- +migrate Down

DROP TABLE `scheduled_campaigns`;
