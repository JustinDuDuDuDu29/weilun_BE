-- name: GetUser :one
SELECT id, role, deleted_date FROM  UserT
WHERE userName=$1 AND pwd=$2 LIMIT 1;

-- name: CreateUser :one
INSERT INTO UserT(
    userName, pwd, role
) VALUES (
  $1, $2, $3
)
RETURNING id;

-- name: UpdateUser :exec
UPDATE UserT
  set userName= $2,
  pwd = $3,
  role = $4,
  last_modified_date = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
UPDATE UserT
  set deleted_date= NOW(),
  last_modified_date = NOW()
WHERE id = $1;

