-- +goose Up
-- One row per calendar month: how many recipients this service has mailed.
-- lib/backend/email reserves against it before every send (SMTP_MONTHLY_LIMIT).
CREATE TABLE email_monthly_sends (
    month date PRIMARY KEY,
    sent  integer NOT NULL
);

-- +goose Down
DROP TABLE email_monthly_sends;
