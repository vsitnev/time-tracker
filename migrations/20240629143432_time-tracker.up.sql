CREATE SCHEMA IF NOT EXISTS md;

CREATE TABLE IF NOT EXISTS md."users" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(36) NOT NULL,
    surname VARCHAR(36) NOT NULL,
    patronymic VARCHAR(36) NOT NULL,
    passport_number VARCHAR(11) NOT NULL UNIQUE,
    "address" VARCHAR(256) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_user_passport_number ON md."users"(passport_number);

CREATE TABLE IF NOT EXISTS md.tasks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    "description" VARCHAR(128) NOT NULL,
    duration INT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES md."users"(id) ON DELETE CASCADE
);