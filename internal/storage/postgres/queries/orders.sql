-- name: CreateOrder :one
INSERT INTO orders (user_id, amount, note)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListOrdersByUser :many
SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC;