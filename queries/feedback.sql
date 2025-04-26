
-- name: InsertFeedback :exec
INSERT INTO feedbacks (created_by, feedback)
VALUES ($1, $2);
