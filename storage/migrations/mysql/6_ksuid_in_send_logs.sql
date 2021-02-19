-- +migrate Up

ALTER TABLE `send_logs` MODIFY `id` varbinary(27);
ALTER TABLE `send_logs` DROP COLUMN `uuid`;

-- +migrate Down

ALTER TABLE `send_logs` MODIFY `id` bigint unsigned;
ALTER TABLE `send_logs` ADD COLUMN `uuid` varchar(36) unique NOT NULL AFTER `id`;
