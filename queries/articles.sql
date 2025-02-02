-- name: InsertArticle :exec
INSERT INTO articles (pod_id, article_text)
VALUES ($1, $2);

-- name: GetArticleByPodId :one
SELECT article_text FROM articles WHERE pod_id = $1 LIMIT 1;

-- name: GetArticleOwner :one
SELECT p.created_by FROM articles a INNER JOIN pods p ON a.pod_id = p.id WHERE a.pod_id = $1 LIMIT 1;