-- +migrate Up

ALTER TABLE `campaigns`
    DROP COLUMN template_name,
    ADD template_id integer unsigned,
    ADD FOREIGN KEY (`template_id`) REFERENCES templates (`id`);
    -- +migrate StatementEnd