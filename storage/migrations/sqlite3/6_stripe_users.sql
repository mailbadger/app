-- +migrate Up

CREATE TABLE IF NOT EXISTS "stripe_users" (
    "id"                         integer primary key autoincrement,
    "user_id"                    integer NOT NULL,
    "stripe_user_id"             varchar(191) NOT NULL,
    "created_at"                 datetime NOT NULL,
    "updated_at"                 datetime NOT NULL,
    foreign key ("user_id") references users("id")
);

-- +migrate Down

DROP TABLE "stripe_users";