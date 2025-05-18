-- +goose Up
CREATE TABLE IF NOT EXISTS "USERS" (
    "id" SERIAL PRIMARY KEY,
    "user_name" VARCHAR(255) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL,
    "password_reset" VARCHAR(255) UNIQUE, 
    "role" VARCHAR(255) NOT NULL,
    "session" VARCHAR(255)
);

-- +goose Down
DROP TABLE IF EXISTS "USERS";
