-- name: InsertQuiz :one
INSERT INTO quizzes (pod_id)
VALUES ($1)
RETURNING id;

-- name: GetQuizByPodId :one
SELECT id,pod_id FROM quizzes WHERE pod_id = $1 LIMIT 1;

-- name: GetQuizOwner :one
SELECT p.created_by FROM quizzes q INNER JOIN pods p ON q.pod_id = p.id WHERE q.pod_id = $1 LIMIT 1;