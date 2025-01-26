-- name: InsertQuiz :one
INSERT INTO quizzes (pod_id, created_by)
VALUES ($1, $2)
RETURNING id;

-- name: GetQuizByPodId :one
SELECT id,pod_id FROM quizzes WHERE pod_id = $1 LIMIT 1;