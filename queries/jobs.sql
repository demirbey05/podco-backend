-- name: InsertJob :one
INSERT INTO jobs(pod_id)
VALUES($1)
RETURNING id;

-- name: GetJobStatusByPodID :one
SELECT job_status 
FROM jobs 
WHERE pod_id = $1;

-- name: GetJobStatusByID :one
SELECT job_status 
FROM jobs 
WHERE id = $1;

-- name: UpdateJobStatusByID :exec
UPDATE jobs
SET job_status = $2
WHERE id = $1;