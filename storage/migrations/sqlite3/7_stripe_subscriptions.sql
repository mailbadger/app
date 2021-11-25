-- +migrate Up

CREATE TABLE IF NOT EXISTS "stripe_subscriptions" (
    "id"              integer primary key autoincrement,
    "user_id"         integer unsigned not null,
    "stripe_id"       varchar(191) unique not null,
    "status"          varchar(191) not null,
    "trial_ends_at"   datetime,
    "ends_at"         datetime,
    "created_at"      datetime not null,
    "updated_at"      datetime not null,
    foreign key ("user_id") references users("id")
);

-- +migrate Down

DROP TABLE "stripe_subscriptions";