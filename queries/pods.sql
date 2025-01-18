-- name: GetPodByLink :many
select * from pods where link = $1;

-- name: InsertPod :one
INSERT INTO pods (link)
VALUES ($1)
RETURNING id;


