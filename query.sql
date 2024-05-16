-- name: GetUser :one
SELECT id, role, deleted_date,pwd FROM  UserT
WHERE phoneNum=$1  LIMIT 1;

-- name: GetDriver :one
SELECT UserT.id as ID, insurances, registration, driverLicense, TruckLicense, nationalidnumber, percentage, cmpt.name as cmpName, usert.phoneNum, usert.name as userName, usert.belongCMP, usert.role, usert.initPwdChanged, DriverT.lastAlert, DriverT.Approved_date, UserT.Deleted_Date as Deleted_Date FROM  DriverT inner join usert on DriverT.id = UserT.id inner join cmpt on usert.belongCMP = cmpt.id where
DriverT.id = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT
UserT.id as ID, cmpt.name as Cmpname, usert.phoneNum as phoneNum, usert.name as Username, usert.belongCMP, usert.role, usert.initPwdChanged, UserT.Deleted_Date as Deleted_Date 
from UserT 
inner join cmpt on UserT.belongcmp = cmpt.id 
where UserT.id=$1 LIMIT 1;

-- name: GetUserSeed :one
SELECT seed from UserT 
inner join cmpt on UserT.belongcmp = cmpt.id 
where UserT.id=$1 LIMIT 1;

-- name: GetUserList :many
SELECT UserT.id as ID, phoneNum, UserT.name as Username, cmpt.name as Cmpname , role, UserT.create_date, UserT.deleted_date, UserT.last_modified_date 
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

-- name: UserHasModified :exec
UPDATE UserT set 
  last_modified_date = NOW()
WHERE id = $1;

-- name: NewSeed :exec
UPDATE UserT set 
  seed = $2,
  last_modified_date = NOW()
WHERE id = $1;

-- name: UpdateDriver :exec
UPDATE DriverT set 
  percentage = COALESCE($2, percentage),
  nationalidnumber = COALESCE($2, nationalidnumber)
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE UserT set 
  pwd = $2,
  initPwdChanged = True,
  last_modified_date = NOW()
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE UserT set 
  phoneNum = COALESCE($2, phoneNum),
  name = COALESCE($3, name),
  belongCMP = COALESCE($4, belongCMP),
  role = COALESCE($5, role),
  last_modified_date = NOW()
WHERE UserT.id = $1;

-- name: UpdateDriverPic :exec
UPDATE DriverT set 
  insurances = COALESCE($2, insurances),
  registration = COALESCE($3, registration),
  driverLicense = COALESCE($4, driverLicense),
  truckLicense = COALESCE($5, truckLicense),
  approved_date = NULL
WHERE DriverT.id = $1;

-- name: ApproveDriver :exec
UPDATE DriverT set 
  UserT.last_modified_date = NOW(),
  approved_date =  NOW()
WHERE DriverT.id = $1;


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

-- name: GetAllJobsClient :many
SELECT  
    JobsT.ID,
    JobsT.From_Loc,
    JobsT.Mid,
    JobsT.To_Loc,
    JobsT.Price,
    JobsT.Remaining,
    JobsT.Belongcmp,
    JobsT.Source,
    JobsT.Jobdate,
    JobsT.Memo,
    JobsT.Close_Date,
    JobsT.deleted_date
from JobsT
inner join cmpt on JobsT.belongcmp = cmpt.id 
where 
(JobsT.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)AND
(JobsT.From_Loc= sqlc.narg('FromLoc') OR sqlc.narg('FromLoc') IS NULL)AND
(JobsT.Mid= sqlc.narg('Mid') OR sqlc.narg('Mid') IS NULL)AND
(JobsT.To_Loc= sqlc.narg('ToLoc') OR sqlc.narg('ToLoc') IS NULL)AND
(belongcmp = sqlc.narg('belongcmp') OR sqlc.narg('belongcmp') IS NULL)AND
(remaining <> 0)AND
(JobsT.close_date is NULL)AND
(JobsT.deleted_date is NULL);

-- AND
-- ((JobsT.create_date > sqlc.narg('create_date_start') OR sqlc.narg('create_date_start') IS NULL)
--  AND (JobsT.create_date < sqlc.narg('create_date_end') OR sqlc.narg('create_date_end') IS NULL)) AND
-- ((JobsT.deleted_date > sqlc.narg('deleted_date_start') OR sqlc.narg('deleted_date_start') IS NULL)
--  AND (JobsT.deleted_date < sqlc.narg('deleted_date_end') OR sqlc.narg('deleted_date_end') IS NULL)) AND
-- ((JobsT.last_modified_date > sqlc.narg('last_modified_date_start') OR sqlc.narg('last_modified_date_start') IS NULL) 
-- AND (JobsT.last_modified_date < sqlc.narg('last_modified_date_end') OR sqlc.narg('last_modified_date_end') IS NULL));

-- name: GetAllJobsAdmin :many
SELECT  *
from JobsT
inner join cmpt on JobsT.belongcmp = cmpt.id 
where 
(JobsT.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)AND
(JobsT.From_Loc= sqlc.narg('FromLoc') OR sqlc.narg('FromLoc') IS NULL)AND
(JobsT.Mid= sqlc.narg('Mid') OR sqlc.narg('Mid') IS NULL)AND
(JobsT.To_Loc= sqlc.narg('ToLoc') OR sqlc.narg('ToLoc') IS NULL)AND
(belongcmp = sqlc.narg('belongcmp') OR sqlc.narg('belongcmp') IS NULL)AND
(remaining <> sqlc.narg('remaining') OR sqlc.narg('remaining') IS NULL)AND
((JobsT.close_date> sqlc.narg('close_date_start') OR sqlc.narg('close_date_start') IS NULL)
 AND (JobsT.create_date < sqlc.narg('close_date_end') OR sqlc.narg('close_date_end') IS NULL)) AND
((JobsT.create_date > sqlc.narg('create_date_start') OR sqlc.narg('create_date_start') IS NULL)
 AND (JobsT.create_date < sqlc.narg('create_date_end') OR sqlc.narg('create_date_end') IS NULL)) AND
((JobsT.deleted_date > sqlc.narg('deleted_date_start') OR sqlc.narg('deleted_date_start') IS NULL)
 AND (JobsT.deleted_date < sqlc.narg('deleted_date_end') OR sqlc.narg('deleted_date_end') IS NULL)) AND
((JobsT.last_modified_date > sqlc.narg('last_modified_date_start') OR sqlc.narg('last_modified_date_start') IS NULL) 
AND (JobsT.last_modified_date < sqlc.narg('last_modified_date_end') OR sqlc.narg('last_modified_date_end') IS NULL));



-- name: GetAllJobsByCmp :many
SELECT * from JobsT where belongcmp = $1;

-- name: GetJobById :one
SELECT * from JobsT where id = $1 LIMIT 1;

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
    close_date = $10,
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
    close_date
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
SELECT ClaimJobT.id as id, JobsT.id as JobID, UserT.id as UserID, JobsT.From_Loc, JobsT.mid, JobsT.To_Loc, ClaimJobT.Create_Date, usert.name as userName, cmpt.name as cmpname, cmpT.id as cmpID, ClaimJobT.Approved_date as ApprovedDate, ClaimJobT.Finished_Date as FinishDate from ClaimJobT inner join JobsT on JobsT.id = ClaimJobT.JobId inner join UserT on UserT.id = ClaimJobT.Driverid inner join Cmpt on UserT.belongCMP = cmpt.id WHERE ClaimJobT.Deleted_date is  null;

-- name: GetClaimedJobByID :one
SELECT ClaimJobT.id as id, JobsT.id as JobID, UserT.id as UserID, JobsT.From_Loc, finished_date, finishPic, JobsT.mid, JobsT.To_Loc, ClaimJobT.Create_Date, usert.name as userName, cmpt.name as cmpname, cmpT.id as cmpID, ClaimJobT.Approved_date as ApprovedDate, DriverT.percentage  as driverPercentage, ClaimJobT.percentage as percentage, price from ClaimJobT inner join JobsT on JobsT.id = ClaimJobT.JobId inner join UserT on UserT.id = ClaimJobT.Driverid inner join Cmpt on UserT.belongCMP = cmpt.id inner join DriverT on driverT.id = UserT.id WHERE ClaimJobT.id = $1;

-- name: ClaimJob :one 
INSERT into ClaimJobT (
    jobID,
    driverID
) values (
    $1,
    $2
) RETURNING id;

-- name: DecreaseRemaining :one
Update JobsT set remaining = remaining - 1, last_modified_date = NOW() where id = $1 RETURNING remaining;

-- name: DeleteClaimedJob :exec 
Update ClaimJobT Set
    deleted_by = $2,
    deleted_date = NOW(),
    last_modified_date = NOW()
    where id = $1;

-- name: IncreaseRemaining :one
Update JobsT set remaining = remaining + 1, last_modified_date = NOW() where id = $1 RETURNING remaining;

-- name: FinishClaimedJob :exec
Update ClaimJobT Set
    finishPic =$3,
    finished_date = NOW(),
    percentage = (SELECT percentage from driverT where driverT.id = (SELECT driverID from ClaimJobT where ClaimJobT.id = $1)),
    last_modified_date = NOW()
WHERE id = $1 and ClaimJobT.Driverid = $2;

-- name: ApproveFinishedJob :exec
Update ClaimJobT set Approved_By = $2, approved_date = NOW(), last_modified_date = NOW() where id = $1;

-- name: GetCurrentClaimedJob :one
SELECT t2.id, t1.*  FROM ClaimJobT t2 inner join JobsT t1 on t1.id = t2.jobID where t2.driverID = $1 and (t2.deleted_date IS NULL and t2.finished_date IS NULL) order by t2.create_date LIMIT 1;

-- name: GetDriverRevenue :many
SELECT coalesce(sum(t1.percentage*t2.price), 0) as earn
, coalesce((select count(*) from ClaimJobT t1 where t1.driverID = $1 
 and (t1.finished_date IS NOT NULL and approved_date IS NOT NULL and t1.deleted_date IS NULL) 
and t1.finished_date between $2 and $3), 0) as count
from ClaimJobT t1 inner join JobsT t2 on t1.jobID = t2.id
where t1.driverID = $1 
and (t1.finished_date IS NOT NULL and approved_date IS NOT NULL and t1.deleted_date IS NULL) 
and t1.finished_date between $2 and $3;

-- name: CreateNewRepair :one
INSERT into repairT (type, driverID, repairInfo) values ($1, $2, $3) RETURNING id;

-- name: GetRepair :many
SELECT repairT.id as ID, UserT.id as Driverid, UserT.Name as Drivername, repairT.type as type, repairT.Repairinfo as Repairinfo, repairT.Create_Date as CreateDate, repairT.Approved_Date as ApprovedDate
from repairT 
inner join UserT on UserT.id = repairT.driverID
where 
(repairT.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)AND
(repairT.driverID = sqlc.narg('driverID') OR sqlc.narg('driverID') IS NULL)AND
(UserT.name = sqlc.narg('name') OR sqlc.narg('name') IS NULL)AND
(UserT.belongcmp = sqlc.narg('belongcmp') OR sqlc.narg('belongcmp') IS NULL)AND
((repairT.create_date > sqlc.narg('create_date_start') OR sqlc.narg('create_date_start') IS NULL)
 AND (repairT.create_date < sqlc.narg('create_date_end') OR sqlc.narg('create_date_end') IS NULL)) AND
repairT.deleted_date is null AND
((repairT.last_modified_date > sqlc.narg('last_modified_date_start') OR sqlc.narg('last_modified_date_start') IS NULL) 
AND (repairT.last_modified_date < sqlc.narg('last_modified_date_end') OR sqlc.narg('last_modified_date_end') IS NULL));

-- name: ApproveRepair :exec
Update repairT set approved_date = NOW(), last_modified_date = NOW() where id =$1;

-- name: DeleteRepair :exec
Update repairT set deleted_date = NOW(), last_modified_date = NOW() where id =$1;

-- name: CreateAlert :one
INSERT INTO AlertT (alert, belongCMP) values ($1, $2) RETURNING id;

-- name: UpdateAlert :exec
Update AlertT Set
alert = $2,
last_modified_date = NOW()
where id = $1;

-- name: DeleteAlert :exec
Update AlertT Set
deleted_date = NOW(),
last_modified_date = NOW()
where id = $1;

-- name: UpdateLastAlert :exec
Update driverT set
lastAlert = $2
where id = $1;

-- name: GetLastAlert :one
SELECT lastAlert from driverT where id = $1;

-- name: GetAlert :many
SELECT alertT.id as ID, cmpt.name as cmpName, alertT.belongCMP as cmpID, alertT.alert as alert, alertT.create_date as Createdate, alertT.Deleted_Date as Deletedate from alertT
inner join cmpt on alertT.Belongcmp = cmpt.id
where 
(alertT.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)AND
(belongCMP = sqlc.narg('belongCMP') OR sqlc.narg('belongCMP') IS NULL)AND
(alert like sqlc.narg('alert') OR sqlc.narg('alert') IS NULL)AND
((alertT.create_date > sqlc.narg('create_date_start') OR sqlc.narg('create_date_start') IS NULL)
 AND (alertT.create_date < sqlc.narg('create_date_end') OR sqlc.narg('create_date_end') IS NULL)) AND
((alertT.deleted_date > sqlc.narg('deleted_date_start') OR sqlc.narg('deleted_date_start') IS NULL)
 AND (alertT.deleted_date < sqlc.narg('deleted_date_end') OR sqlc.narg('deleted_date_end') IS NULL)) AND
((alertT.last_modified_date > sqlc.narg('last_modified_date_start') OR sqlc.narg('last_modified_date_start') IS NULL) 
AND (alertT.last_modified_date < sqlc.narg('last_modified_date_end') OR sqlc.narg('last_modified_date_end') IS NULL))
order by alertT.id desc;

-- name: GetAlertByCmp :many
SELECT *
from alertT
where belongCMP = $1 order by id desc;

-- name: GetRepairById :one
SELECT * from repairT where id = $1;
