-- +migrate Up

CREATE TABLE IF NOT EXISTS "stripe_customers" (
    "id"                         integer primary key autoincrement,
    "user_id"                    integer unsigned NOT NULL,
    "customer_id"                varchar(191) NOT NULL,
    "created_at"                 datetime NOT NULL,
    "updated_at"                 datetime NOT NULL,
    foreign key ("user_id") references users("id")
);

-- +migrate Down

DROP TABLE "stripe_customers";