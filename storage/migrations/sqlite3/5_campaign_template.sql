-- +migrate Up

ALTER TABLE campaigns
    DROP COLUMN template_name,
    ADD  template_id INT,
    FOREIGN KEY(tmp_id) REFERENCES templates(id)

-- +migrate Down
