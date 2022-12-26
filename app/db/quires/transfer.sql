-- name: Createtransfer :one
INSERT INTO transfer (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: Gettransfer :one
SELECT * FROM transfer
WHERE id = $1 LIMIT 1;

-- name: Gettransfers :many
SELECT * FROM transfer
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: Updatetransfer :one
UPDATE transfer
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: Deletetransfer :exec
DELETE FROM transfer WHERE id = $1;