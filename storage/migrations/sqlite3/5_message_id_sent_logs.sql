-- +migrate Up

ALTER TABLE "send_logs"
    ADD COLUMN "message_id" varchar(191);

-- +migrate Down

ALTER TABLE "send_logs"
    DROP COLUMN "message_id";