-- name: InsertQuiz :one
INSERT INTO quizzes (pod_id, created_by)
VALUES ($1, $2)
RETURNING id;