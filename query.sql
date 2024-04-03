-- name: GetUser :one
SELECT id, role, deleted_date FROM  UserT
WHERE userName=$1 AND pwd=$2 LIMIT 1;

-- name: CreateAdmin :one
INSERT INTO UserT(
    userName, pwd, name, role, belongcmp, phoneNum
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id;

-- name: CreateCmpAdmin :one
INSERT INTO UserT(
    userName, pwd, name, role, belongcmp, phoneNum
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id;

-- name: CreateDriver :one
with createUser as (
    INSERT INTO UserT(
        userName, pwd, name, role, belongcmp, phoneNum
        ) VALUES (
        $1, $2, $3, $4, $5, $6
    )
    RETURNING id
)

insert into driverT (userid, percentage, nationalidnumber) 
	(select o.id, v.percentage, v.nationalidnumber
	from createUser o
	cross join(
		values 
		($7, $8)
	) as v (percentage, nationalidnumber))
RETURNING userid;



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

