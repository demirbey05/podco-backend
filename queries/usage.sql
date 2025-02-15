-- name: GetRemainingCredits :one
SELECT credits from usage WHERE user_id = $1;

-- name: IsCreditExist :one
SELECT EXISTS(SELECT 1 FROM usage WHERE user_id = $1);

-- name: InsertCredit :exec
INSERT INTO usage (user_id, credits) VALUES ($1, $2);

-- name: UpdateCredit :one
UPDATE usage SET credits = $2 WHERE user_id = $1 RETURNING credits;

-- name: DeleteCredit :exec
DELETE FROM usage WHERE user_id = $1;

-- name: DecrementCredit :one
UPDATE usage SET credits = credits - $1 WHERE user_id = $2 RETURNING credits;