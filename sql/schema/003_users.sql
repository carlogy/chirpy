-- +goose Up
ALTER TABLE users
ADD COLUMN IF NOT EXISTS hashed_password VARCHAR NOT NULL DEFAULT 'unset';
