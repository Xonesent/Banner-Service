-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS banner_schema;

CREATE TABLE banner_schema.banners(
    banner_id   BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT,
    text TEXT,
    url TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    is_active BOOLEAN
);

CREATE TABLE banner_schema.banners_X_tags(
    id     BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    banner_id INTEGER,
    tag_id INTEGER,
    feature_id INTEGER,
    FOREIGN KEY (banner_id) REFERENCES banner_schema.banners(banner_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF  EXISTS  banner_schema.banners;
DROP TABLE IF  EXISTS  banner_schema.banners_X_tags;
-- +goose StatementEnd
