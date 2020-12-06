-- +migrate Up

CREATE TABLE IF NOT EXISTS "templates"
(
    "id"         integer primary key autoincrement,
    `name`       varchar(191) NOT NULL,
    "subject"    varchar(191) NOT NULL,
    `text_part`  text,
    "created_at" datetime,
    "updated_at" datetime
);

-- +migrate Down

DROP TABLE "templates";