-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS usage (
  id SERIAL PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL,
  credits INT NOT NULL DEFAULT 15000
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS usage;
-- +goose StatementEnd
