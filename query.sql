-- name: GetUser :one
SELECT id, role, deleted_date FROM  UserT
WHERE phoneNum=$1 AND pwd=$2 LIMIT 1;

-- name: CreateAdmin :one
INSERT INTO UserT(
    pwd, name, role, belongcmp, phoneNum
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id;

-- name: CreateCmpAdmin :one
INSERT INTO UserT(
    pwd, name, role, belongcmp, phoneNum
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id;

-- name: CreateDriver :one
with createUser as (
    INSERT INTO UserT(
        pwd, name, role, belongcmp, phoneNum
        ) VALUES (
        $1, $2, $3, $4, $5
    )
    RETURNING id
)

insert into driverT (userid, percentage, nationalidnumber) 
	(select o.id, v.percentage, v.nationalidnumber
	from createUser o
	cross join(
		values 
		($6, $7)
	) as v (percentage, nationalidnumber))
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

