-- name: CreateEntry :one
insert into entries(
    account_id,
    amount
)
VALUES (
  $1, $2
) RETURNING *;

-- name: GetEntry :one
select * from entries
where id=$1 LIMIT 1;

-- name: ListEntries :many
select * from entries
where account_id =$1
order by id
LIMIT $2
OFFSET $3;