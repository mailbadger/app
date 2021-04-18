-- +migrate Up

CREATE TABLE IF NOT EXISTS "campaign_schedules"
(
    "id"                    varchar(27) primary key,
    "user_id"               integer,
    "campaign_id"           integer,
    "scheduled_at"          datetime,
    "source"                varchar,
    "from_name"             varchar,
    "segment_ids"           varchar,
    "default_template_data" varchar,
    "created_at"            datetime,
    "updated_at"            datetime,
    foreign key ("campaign_id") references campaigns("id")
    );

-- +migrate Down
DROP TABLE "campaign_schedules";