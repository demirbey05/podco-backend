-- name: InsertQuestion :one
INSERT INTO questions (quizzes_id, question_text, options, correct_option)
VALUES ($1, $2, $3, $4)
RETURNING id;