-- +goose Up
CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "email" varchar not null unique
);

-- +goose Down
DROP TABLE IF EXISTS users;