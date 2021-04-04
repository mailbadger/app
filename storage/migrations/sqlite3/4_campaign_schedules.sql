-- +migrate Up

CREATE TABLE IF NOT EXISTS `campaign_schedules`
(
    `id`           varchar(27) primary key,
    `campaign_id`  integer,
    `scheduled_at` datetime,
    "created_at"   datetime,
    "updated_at"   datetime
    );

-- +migrate Down
DROP TABLE `campaign_schedules`;