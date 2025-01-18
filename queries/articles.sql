-- name: InsertArticle :exec
INSERT INTO articles (pod_id, created_by, article_text)
VALUES ($1, $2, $3);