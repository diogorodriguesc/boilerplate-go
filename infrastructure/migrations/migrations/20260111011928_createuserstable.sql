-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    CONSTRAINT users_email_key UNIQUE (email)
);
-- +goose StatementEnd
