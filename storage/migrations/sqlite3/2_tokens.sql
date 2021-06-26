-- +migrate Up

CREATE TABLE IF NOT EXISTS "tokens" (
    "id"         integer primary key autoincrement,
    "user_id"    integer NOT NULL,
    "token"      varchar(191) NOT NULL UNIQUE,
    "type"       varchar(191) NOT NULL,
    "expires_at" datetime NOT NULL,
    "created_at" datetime NOT NULL,
    "updated_at" datetime NOT NULL,
    foreign key ("user_id") references users("id")
);

-- +migrate Down

DROP TABLE "tokens";