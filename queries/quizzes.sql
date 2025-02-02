-- name: InsertQuiz :one
INSERT INTO quizzes (pod_id)
VALUES ($1)
RETURNING id;

-- name: GetQuizByPodId :one
SELECT id,pod_id FROM quizzes WHERE pod_id = $1 LIMIT 1;