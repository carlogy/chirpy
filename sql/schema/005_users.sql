-- +goose Up
ALTER TABLE users
ADD COLUMN IF NOT EXISTS is_chirpy_red boolean DEFAULT FALSE;
