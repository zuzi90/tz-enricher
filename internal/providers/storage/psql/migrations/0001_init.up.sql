-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id          serial     PRIMARY KEY,
    name        varchar    NOT NULL,
    surname     varchar,
    patronymic  varchar    ,
    age         int CHECK (age >=0) NOT NULL,
    gender      varchar    ,
    nationality varchar    ,
    is_deleted  bool        NOT NULL DEFAULT FALSE,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd