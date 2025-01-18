-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    pod_id INT NOT NULL REFERENCES pods(id),
    job_status int NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS jobs;
-- +goose StatementEnd
