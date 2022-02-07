-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
select * from accounts
WHERE id =$1 LIMIT 1;

-- name: ListAccounts :many
select * from accounts
order by id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
update accounts
set balance =$2
where id=$1
RETURNING *;

-- name: DeleteAccount :exec
delete from accounts
where id=$1;