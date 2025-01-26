-- name: InsertArticle :exec
INSERT INTO articles (pod_id, created_by, article_text)
VALUES ($1, $2, $3);

-- name: GetArticleByPodId :one
SELECT article_text FROM articles WHERE pod_id = $1 LIMIT 1;
