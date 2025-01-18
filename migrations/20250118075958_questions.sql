-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    quizzes_id INT REFERENCES quizzes(id),
    question_text TEXT NOT NULL,
    options TEXT[] NOT NULL,
    correct_option INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS questions;
-- +goose StatementEnd
