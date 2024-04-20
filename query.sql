-- name: GetUser :one
SELECT id, role, deleted_date FROM  UserT
WHERE phoneNum=$1 AND pwd=$2 LIMIT 1;

-- name: GetUserByID :one
SELECT UserT.id, phoneNum, UserT.name, cmpt.name, role, UserT.create_date, UserT.deleted_date, UserT.last_modified_date 
from UserT 
inner join cmpt on UserT.belongcmp = cmpt.id 
where UserT.id=$1 LIMIT 1;

-- name: GetUserList :many
SELECT UserT.id, phoneNum, UserT.name, cmpt.name, role, UserT.create_date, UserT.deleted_date, UserT.last_modified_date 
from UserT 
inner join cmpt on UserT.belongcmp = cmpt.id 
where 
(UserT.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)AND
(phoneNum = sqlc.narg('phoneNum')::Text OR sqlc.narg('phoneNum')::Text IS NULL)AND
(UserT.name = sqlc.narg('name') OR sqlc.narg('name') IS NULL)AND
(belongcmp = sqlc.narg('belongcmp') OR sqlc.narg('belongcmp') IS NULL)AND
((UserT.create_date > sqlc.narg('create_date_start') OR sqlc.narg('create_date_start') IS NULL)
 AND (UserT.create_date < sqlc.narg('create_date_end') OR sqlc.narg('create_date_end') IS NULL)) AND
((UserT.deleted_date > sqlc.narg('deleted_date_start') OR sqlc.narg('deleted_date_start') IS NULL)
 AND (UserT.deleted_date < sqlc.narg('deleted_date_end') OR sqlc.narg('deleted_date_end') IS NULL)) AND
((UserT.last_modified_date > sqlc.narg('last_modified_date_start') OR sqlc.narg('last_modified_date_start') IS NULL) 
AND (UserT.last_modified_date < sqlc.narg('last_modified_date_end') OR sqlc.narg('last_modified_date_end') IS NULL));

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
insert into driverT (id, percentage, nationalidnumber) 
    values ($1, $2, $3)
RETURNING id;

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
    set
    name = $2,
    last_modified_date = NOW()
WHERE id = $1;

-- name: DeleteCmp :exec
UPDATE cmpt
  set deleted_date= NOW(),
  last_modified_date = NOW()
WHERE id = $1;


-- name: GetAllJobs :many
SELECT * from JobsT;


-- name: GetAllJobsByCmp :many
SELECT * from JobsT where belongcmp = $1;

-- name: SetJobNoMore :exec
UPDATE JobsT 
  set finished_date = NOW(),
  last_modified_date = NOW()
WHERE id = $1;
 
-- name: DeleteJob :exec
UPDATE JobsT 
  set deleted_date = NOW(),
  last_modified_date = NOW()
WHERE id = $1;


-- name: UpdateJob :one
UPDATE JobsT set 
    from_loc = $1,
    mid = $2,
    to_loc = $3,
    price = $4,
    remaining = $5,
    belongCMP = $6,
    source = $7,
    jobDate = $8,
    memo = $9,
    end_date = $10,
    last_modified_date = NOW()
where id = $11
 RETURNING id;


-- name: CreateJob :one
INSERT INTO JobsT (
    from_loc,
    mid,
    to_loc,
    price,
    estimated,
    remaining,
    belongCMP,
    source,
    jobDate,
    memo,
    end_date
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
) RETURNING id;

-- name: GetAllClaimedJobs :many
SELECT * from ClaimJobT;

-- name: GetClaimedJobByID :one
SELECT * from ClaimJobT where id = $1;

-- name: ClaimJob :one 
INSERT into ClaimJobT (
    jobID,
    driverID
) values (
    $1,
    $2
) RETURNING id;

-- name: DecreaseRemaining :exec
Update JobsT set remaining = remaining - 1, last_modified_date = NOW() where id = $1;

-- name: DeleteClaimedJob :exec 
Update ClaimJobT Set
    deleted_by = $2,
    deleted_date = NOW(),
    last_modified_date = NOW()
    where id = $1;

-- name: IncreaseRemaining :exec
Update JobsT set remaining = remaining + 1, last_modified_date = NOW() where id = $1;

-- name: FinishClamedJob :exec
Update ClaimJobT Set
    finished_date = NOW(),
    percentage = (SELECT percentage from driverT where usert.id = (SELECT driverID from ClaimJobT where ClaimJobT.id = $1)),
    last_modified_date = NOW()
WHERE id = $1 and ClaimJobT.Driverid = $2;

-- name: ApproveFinishedJob :exec
Update ClaimJobT set Approved_By = $2, approved_date = NOW(), last_modified_date = NOW() where id = $1;

-- name: GetClaimedJob :one
SELECT t2.id, t1.*  FROM ClaimJobT t2 inner join JobsT t1 on t1.id = t2.jobID where t2.driverID = $1 and (t2.deleted_date IS NULL and t2.finished_date IS NULL) order by t2.create_date LIMIT 1;

-- name: GetDriverRevenue :many
SELECT t1.percentage*t2.price as earn from ClaimJobT t1 inner join JobsT t2 on t1.jobID = t2.id where t1.driverID = $1 and (t1.finished_date IS NOT NULL 
    and approved_date IS NOT NULL and deleted_date IS NOT NULL) and t1.finished_date 
    between $2 and $3;
