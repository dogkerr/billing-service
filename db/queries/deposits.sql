-- name: CreateDeposit :one
INSERT INTO deposits_table (user_id, value, status)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetDeposit :one
SELECT *
FROM deposits_table
WHERE id = $1
LIMIT 1;
-- name: ListDeposits :many
SELECT *
FROM deposits_table
WHERE user_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;
-- name: UpdateDeposit :one
UPDATE deposits_table
SET status = $2
WHERE id = $1
RETURNING *;
-- name: DeleteDeposit :exec
DELETE FROM deposits_table
WHERE id = $1;