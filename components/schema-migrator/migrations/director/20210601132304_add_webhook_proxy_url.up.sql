BEGIN;

ALTER TABLE webhooks
    ADD COLUMN proxy_url VARCHAR(255);

COMMIT;
