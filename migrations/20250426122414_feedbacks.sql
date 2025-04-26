-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS feedbacks (
    created_by VARCHAR(255) NOT NULL,
    feedback JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feedbacks;
-- +goose StatementEnd
