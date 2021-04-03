-- +migrate Up

CREATE TABLE IF NOT EXISTS `scheduled_campaigns`
(
    `id`           varchar(27) primary key,
    `campaign_id`  integer,
    `scheduled_at` datetime,
    "created_at"   datetime,
    "updated_at"   datetime
    );

-- +migrate Down
DROP TABLE `scheduled_campaigns`;