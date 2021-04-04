-- +migrate Up

CREATE TABLE IF NOT EXISTS `tokens` (
  `id`         integer primary key autoincrement,
  `user_id`    integer,
  `token`      varchar(191) NOT NULL UNIQUE,
  `type`       varchar(191) NOT NULL,
  `expires_at` datetime NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
);

-- +migrate Down

DROP TABLE `tokens`;