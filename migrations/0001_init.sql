-- +goose Up
CREATE TABLE IF NOT EXISTS "users" (
  "id"    bigserial PRIMARY KEY,
  "name"  text NOT NULL,
  "email" text NOT NULL
);

-- +goose Down
DROP TABLE "users";
