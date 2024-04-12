-- name: CreateCharge :one
INSERT INTO charges_table (
    container_id,
    user_id,
    total_cpu_usage,
    total_memory_usage,
    total_net_ingress_usage,
    total_net_egress_usage,
    timestamp,
    total_cost
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
-- name: GetCharge :one
SELECT *
FROM charges_table
WHERE id = $1
LIMIT 1;
-- name: ListCharges :many
SELECT *
FROM charges_table
WHERE user_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;
-- name: UpdateCharge :one
UPDATE charges_table
SET total_cost = $2
WHERE id = $1
RETURNING *;
-- name: DeleteCharge :exec
DELETE FROM charges_table
WHERE id = $1;