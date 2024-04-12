-- name: CreateMutation :one
INSERT INTO mutations_table (
    user_id,
    mutation,
    balance,
    type,
    deposit_id,
    charge_id
  )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: GetMutation :one
SELECT *
FROM mutations_table
WHERE id = $1
LIMIT 1;
-- name: ListMutations :many
SELECT *
FROM mutations_table
WHERE user_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;
-- name: UpdateMutation :one
UPDATE mutations_table
SET balance = $2,
  type = $3
WHERE id = $1
RETURNING *;
-- name: DeleteMutation :exec
DELETE FROM mutations_table
WHERE id = $1;