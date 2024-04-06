-- name: GetUser :one
SELECT id, role, deleted_date FROM  UserT
WHERE phoneNum=$1 AND pwd=$2 LIMIT 1;

-- name: GetUserByID :one
SELECT * from UserT where id=$1;

-- name: CreateAdmin :one
INSERT INTO UserT(
    pwd, name, role, belongcmp, phoneNum
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id;

-- name: CreateUser :one
INSERT INTO UserT(
    pwd, name, role, belongcmp, phoneNum
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id;

-- name: CreateDriverInfo :one
insert into driverT (userid, percentage, nationalidnumber) 
    values ($1, $2, $3)
RETURNING userid;

-- name: UpdateUser :exec
UPDATE UserT
  set 
  pwd = $2,
  role = $3,
  last_modified_date = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
UPDATE UserT
  set deleted_date= NOW(),
  last_modified_date = NOW()
WHERE id = $1;


-- name: GetCmp :one
SELECT * FROM cmpt
inner join usert
on cmpt.id = usert.belongcmp AND (usert.role=200 OR usert.role=100)
where cmpt.id = $1;

-- name: GetAllCmp :many
SELECT * from cmpt;

-- name: NewCmp :one
INSERT INTO cmpt (name) values ($1) RETURNING id;

-- name: UpdateCmp :exec
UPDATE cmpt
  set deleted_date= NOW(),
  last_modified_date = NOW()
WHERE id = $1;

-- name: DeleteCmp :exec
UPDATE cmpt
  set deleted_date= NOW(),
  last_modified_date = NOW()
WHERE id = $1;
