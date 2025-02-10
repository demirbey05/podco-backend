-- name: GetPodByLink :many
select * from pods where link = $1;

-- name: InsertPod :one
INSERT INTO pods (link,title,created_by)
VALUES ($1,$2,$3)
RETURNING id;


-- name: GetPodsByUserID :many
SELECT * FROM pods WHERE created_by = $1;

-- name: UpdatePodIsPublic :exec
UPDATE pods SET is_public = $1 WHERE id = $2;
