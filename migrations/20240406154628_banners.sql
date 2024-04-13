-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS banner_schema;

CREATE TABLE banner_schema.banners(
    banner_id   BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT,
    text TEXT,
    url TEXT,
    feature_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    is_active BOOLEAN,
    version INTEGER
);

CREATE TABLE banner_schema.banners_X_tags(
    id     BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    banner_id INTEGER,
    tag_id INTEGER,
    FOREIGN KEY (banner_id) REFERENCES banner_schema.banners(banner_id)
);

CREATE TABLE banner_schema.banners_versions(
    id     BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    banner_id INTEGER,
    title TEXT,
    text TEXT,
    url TEXT,
    feature_id INTEGER,
    tag_ids BIGINT ARRAY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    is_active BOOLEAN,
    version INTEGER,
    FOREIGN KEY (banner_id) REFERENCES banner_schema.banners(banner_id)
);

CREATE INDEX idx_banner_id_banners ON banner_schema.banners(banner_id);
CREATE INDEX idx_banner_id_banners_X_tags ON banner_schema.banners_X_tags(banner_id);
CREATE INDEX idx_banner_id_versions ON banner_schema.banners_versions(banner_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF  EXISTS  banner_schema.banners_X_tags;
DROP TABLE IF  EXISTS  banner_schema.versions;
DROP TABLE IF  EXISTS  banner_schema.banners;

-- +goose StatementEnd
