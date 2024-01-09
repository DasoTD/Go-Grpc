-- name: CreateAccount :one
INSERT INTO accounts (
  username,
  firstname,
  lastname,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- -- name: GetAccountForUpdate :one
-- SELECT * FROM accounts
-- WHERE id = $1 LIMIT 1
-- FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts;