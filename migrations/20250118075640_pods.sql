-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pods (
    id SERIAL PRIMARY KEY,
    link TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

    
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pods;
-- +goose StatementEnd
