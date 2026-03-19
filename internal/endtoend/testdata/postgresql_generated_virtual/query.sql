-- name: GetCustomer :one
SELECT * FROM customers WHERE id = $1;

-- name: CreateCustomer :one
INSERT INTO customers (id, height_cm)
VALUES ($1, $2)
RETURNING *;
