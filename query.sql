-- name: GetUser :one
SELECT id,
  role,
  deleted_date,
  pwd
FROM UserT
WHERE phoneNum = $1
LIMIT 1;
-- name: GetDriver :one
SELECT UserT.id as ID,
  insurances,
  registration,
  driverLicense,
  TruckLicense,
  nationalidnumber,
  -- percentage,
  cmpt.name as cmpName,
  usert.phoneNum,
  usert.name as userName,
  usert.belongCMP,
  usert.role,
  usert.initPwdChanged,
  DriverT.lastAlert,
  DriverT.Approved_date,
  UserT.Deleted_Date as Deleted_Date
FROM DriverT
  inner join usert on DriverT.id = UserT.id
  inner join cmpt on usert.belongCMP = cmpt.id
where DriverT.id = $1
LIMIT 1;
-- name: GetUserByID :one
SELECT UserT.id as ID,
  cmpt.name as Cmpname,
  usert.phoneNum as phoneNum,
  usert.name as Username,
  usert.belongCMP,
  usert.role,
  usert.initPwdChanged,
  UserT.Deleted_Date as Deleted_Date,
  insurances,
  registration,
  driverLicense,
  TruckLicense,
  nationalidnumber,
  --  percentage,
  plateNum,
  Approved_date
from UserT
  inner join cmpt on UserT.belongcmp = cmpt.id
  left join DriverT on driverT.id = usert.id
where UserT.id = $1
LIMIT 1;
-- name: GetUserSeed :one
SELECT seed,
  UserT.deleted_date
from UserT
  inner join cmpt on UserT.belongcmp = cmpt.id
where UserT.id = $1
LIMIT 1;
-- name: GetUserList :many
SELECT json_build_object(
    'cmpid',
    cmpt.id,
    'Cmpname',
    cmpt.name,
    'list',
    json_agg(
      json_build_object(
        'id',
        UT.id,
        'phoneNum',
        UT.phonenum,
        'Role',
        UT.Role,
        'Username',
        UT.name,
        'deleted_date',
        UT.deleted_date,
        'last_modified_date',
        UT.last_modified_date
      )
    )
  )
from UserT UT
  right join cmpt on UT.belongcmp = cmpt.id
where (
    UT.id = sqlc.narg('id')
    OR sqlc.narg('id') IS NULL
  )
  AND (
    phoneNum = sqlc.narg('phoneNum')::Text
    OR sqlc.narg('phoneNum')::Text IS NULL
  )
  AND (
    UT.name like sqlc.narg('name')
    OR sqlc.narg('name') IS NULL
  )
  AND (
    belongcmp = sqlc.narg('belongcmp')
    OR sqlc.narg('belongcmp') IS NULL
  )
  AND (
    cmpt.Name like sqlc.narg('belongcmpName')
    OR sqlc.narg('belongcmpName') IS NULL
  )
  AND (
    (
      UT.create_date > sqlc.narg('create_date_start')
      OR sqlc.narg('create_date_start') IS NULL
    )
    AND (
      UT.create_date < sqlc.narg('create_date_end')
      OR sqlc.narg('create_date_end') IS NULL
    )
  )
  AND (
    (
      UT.deleted_date > sqlc.narg('deleted_date_start')
      OR sqlc.narg('deleted_date_start') IS NULL
    )
    AND (
      UT.deleted_date < sqlc.narg('deleted_date_end')
      OR sqlc.narg('deleted_date_end') IS NULL
    )
  )
  AND (
    (
      UT.last_modified_date > sqlc.narg('last_modified_date_start')
      OR sqlc.narg('last_modified_date_start') IS NULL
    )
    AND (
      UT.last_modified_date < sqlc.narg('last_modified_date_end')
      OR sqlc.narg('last_modified_date_end') IS NULL
    )
  )
  AND(UT.Deleted_Date is null)
group by cmpt.id;
-- name: CreateAdmin :one
INSERT INTO UserT(
    pwd,
    name,
    role,
    belongcmp,
    phoneNum,
    phoneNumInD
  )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
-- name: CreateUser :one
INSERT INTO UserT(
    pwd,
    name,
    role,
    belongcmp,
    phoneNum,
    phoneNumInD
  )
VALUES ($1, $2, $3, $4, $5, $5)
RETURNING id;
-- name: CreateDriverInfo :one
insert into driverT (id, nationalidnumber, plateNum)
values ($1, $2, $3)
RETURNING id;
-- name: UserHasModified :exec
UPDATE UserT
set last_modified_date = NOW()
WHERE id = $1;
-- name: NewSeed :exec
UPDATE UserT
set seed = $2,
  last_modified_date = NOW()
WHERE id = $1;
-- name: UpdateDriver :exec
UPDATE DriverT
set -- percentage = COALESCE($2, percentage),
  nationalidnumber = COALESCE($2, nationalidnumber),
  plateNum = COALESCE($3, plateNum)
WHERE id = $1;
-- name: UpdateUserPassword :exec
UPDATE UserT
set pwd = $2,
  initPwdChanged = True,
  last_modified_date = NOW()
WHERE id = $1;
-- name: UpdateUser :exec
UPDATE UserT
set phoneNum = COALESCE($2, phoneNum),
  name = COALESCE($3, name),
  belongCMP = COALESCE($4, belongCMP),
  role = COALESCE($5, role),
  last_modified_date = NOW()
WHERE UserT.id = $1;
-- name: UpdateDriverPic :exec
UPDATE DriverT
set insurances = COALESCE($2, insurances),
  registration = COALESCE($3, registration),
  driverLicense = COALESCE($4, driverLicense),
  truckLicense = COALESCE($5, truckLicense),
  approved_date = NULL
WHERE DriverT.id = $1;
-- name: ApproveDriver :exec
UPDATE DriverT
set approved_date = NOW()
where id = $1;
-- name: DeleteUser :exec
UPDATE UserT
set deleted_date = NOW(),
  phoneNum = null,
  last_modified_date = NOW()
WHERE id = $1;
-- name: GetCmp :one
SELECT *
FROM cmpt
  inner join usert on cmpt.id = usert.belongcmp
  AND (
    usert.role = 200
    OR usert.role = 100
  )
where cmpt.id = $1;
-- name: GetAllCmp :many
SELECT *
from cmpt;
-- name: NewCmp :one
INSERT INTO cmpt (name)
values ($1)
RETURNING id;
-- name: UpdateCmp :exec
UPDATE cmpt
set name = $2,
  last_modified_date = NOW()
WHERE id = $1;
-- name: DeleteCmp :exec
UPDATE cmpt
set deleted_date = NOW(),
  last_modified_date = NOW()
WHERE id = $1;
-- name: GetAllJobsClient :many
SELECT JobsT.ID,
  JobsT.From_Loc,
  JobsT.Mid,
  JobsT.To_Loc,
  JobsT.Price,
  JobsT.Remaining,
  JobsT.Belongcmp,
  JobsT.Source,
  JobsT.Memo,
  -- JobsT.Close_Date,
  JobsT.deleted_date
from JobsT
  inner join cmpt on JobsT.belongcmp = cmpt.id
where (
    JobsT.id = sqlc.narg('id')
    OR sqlc.narg('id') IS NULL
  )
  AND (
    JobsT.From_Loc = sqlc.narg('FromLoc')
    OR sqlc.narg('FromLoc') IS NULL
  )
  AND (
    JobsT.Mid = sqlc.narg('Mid')
    OR sqlc.narg('Mid') IS NULL
  )
  AND (
    JobsT.To_Loc = sqlc.narg('ToLoc')
    OR sqlc.narg('ToLoc') IS NULL
  )
  AND (
    belongcmp = sqlc.narg('belongcmp')
    OR sqlc.narg('belongcmp') IS NULL
  )
  AND (remaining <> 0)
  AND -- (JobsT.close_date is NULL)AND
  (JobsT.deleted_date is NULL);
-- AND
-- ((JobsT.create_date > sqlc.narg('create_date_start') OR sqlc.narg('create_date_start') IS NULL)
--  AND (JobsT.create_date < sqlc.narg('create_date_end') OR sqlc.narg('create_date_end') IS NULL)) AND
-- ((JobsT.deleted_date > sqlc.narg('deleted_date_start') OR sqlc.narg('deleted_date_start') IS NULL)
--  AND (JobsT.deleted_date < sqlc.narg('deleted_date_end') OR sqlc.narg('deleted_date_end') IS NULL)) AND
-- ((JobsT.last_modified_date > sqlc.narg('last_modified_date_start') OR sqlc.narg('last_modified_date_start') IS NULL) 
-- AND (JobsT.last_modified_date < sqlc.narg('last_modified_date_end') OR sqlc.narg('last_modified_date_end') IS NULL));
-- name: GetAllJobsAdmin :many
SELECT *
from JobsT
  inner join cmpt on JobsT.belongcmp = cmpt.id
where (
    JobsT.id = sqlc.narg('id')
    OR sqlc.narg('id') IS NULL
  )
  AND (
    JobsT.From_Loc = sqlc.narg('FromLoc')
    OR sqlc.narg('FromLoc') IS NULL
  )
  AND (
    JobsT.Mid = sqlc.narg('Mid')
    OR sqlc.narg('Mid') IS NULL
  )
  AND (
    JobsT.To_Loc = sqlc.narg('ToLoc')
    OR sqlc.narg('ToLoc') IS NULL
  )
  AND (
    belongcmp = sqlc.narg('belongcmp')
    OR sqlc.narg('belongcmp') IS NULL
  )
  AND (
    remaining <> sqlc.narg('remaining')
    OR sqlc.narg('remaining') IS NULL
  )
  AND -- ((JobsT.close_date> sqlc.narg('close_date_start') OR sqlc.narg('close_date_start') IS NULL)
  --  AND (JobsT.create_date < sqlc.narg('close_date_end') OR sqlc.narg('close_date_end') IS NULL)) AND
  (
    (
      JobsT.create_date > sqlc.narg('create_date_start')
      OR sqlc.narg('create_date_start') IS NULL
    )
    AND (
      JobsT.create_date < sqlc.narg('create_date_end')
      OR sqlc.narg('create_date_end') IS NULL
    )
  )
  AND (
    (
      JobsT.last_modified_date > sqlc.narg('last_modified_date_start')
      OR sqlc.narg('last_modified_date_start') IS NULL
    )
    AND (
      JobsT.last_modified_date < sqlc.narg('last_modified_date_end')
      OR sqlc.narg('last_modified_date_end') IS NULL
    )
  )
  AND(JobsT.deleted_date is NULL);
;
-- AND (
--   (
--     JobsT.deleted_date > sqlc.narg('deleted_date_start')
--     OR sqlc.narg('deleted_date_start') IS NULL
--   )
--   AND (
--     JobsT.deleted_date < sqlc.narg('deleted_date_end')
--     OR sqlc.narg('deleted_date_end') IS NULL
--   )
-- )
-- name: GetAllJobsByCmp :many
SELECT *
from JobsT
where belongcmp = $1;
-- name: GetJobById :one
SELECT *
from JobsT
where id = $1
LIMIT 1;
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
UPDATE JobsT
set from_loc = $1,
  mid = $2,
  to_loc = $3,
  price = $4,
  remaining = $5,
  belongCMP = $6,
  source = $7,
  -- jobDate = $8,
  memo = $8,
  -- close_date = $9,
  last_modified_date = NOW()
where id = $9
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
    -- jobDate,
    memo -- close_date
  )
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $5,
    $6,
    $7,
    $8 -- $9
    -- $10
  )
RETURNING id;
-- name: GetAllClaimedJobs :many
SELECT ClaimJobT.id as id,
  JobsT.id as JobID,
  UserT.id as UserID,
  JobsT.From_Loc,
  JobsT.mid,
  JobsT.To_Loc,
  JobsT.Price,
  ClaimJobT.Create_Date,
  usert.name as userName,
  cmpt.name as cmpname,
  cmpT.id as cmpID,
  ClaimJobT.Approved_date as ApprovedDate,
  ClaimJobT.Finished_Date as FinishDate,
  ClaimJobT.finishPic
from ClaimJobT
  inner join JobsT on JobsT.id = ClaimJobT.JobId
  inner join UserT on UserT.id = ClaimJobT.Driverid
  inner join Cmpt on UserT.belongCMP = cmpt.id
WHERE ClaimJobT.Deleted_date is null
  and (
    ClaimJobT.driverid = sqlc.narg('uid')
    OR sqlc.narg('uid') IS NULL
  )
  and (
    claimjobt.jobID = sqlc.narg('jobid')
    OR sqlc.narg('jobid') IS NULL
  )
  and (
    usert.belongCMP = sqlc.narg('cmpID')
    OR sqlc.narg('cmpID') IS NULL
  )
  and (
    claimjobt.id = sqlc.narg('cjID')
    OR sqlc.narg('cjID') IS NULL
  )
  and (
    (
      sqlc.narg('cat') = 'pending'
      AND claimjobt.Approved_date IS NULL
    )
    OR (sqlc.narg('cat') IS NULL)
  )
  and (
    to_char(date(claimjobt.create_date), 'YYYY-MM') = to_char(date(sqlc.narg('ym')), 'YYYY-MM')
    OR sqlc.narg('ym') IS NULL
  )
  and (claimjobt.deleted_date IS NULL) -- and (
  --   claimjobt.approved_date = sqlc.narg('ym')
  --   OR sqlc.narg('ym') IS NULL
  -- )
;
-- name: GetClaimedJobByDriverID :many
SELECT ClaimJobT.id as id,
  JobsT.id as JobID,
  UserT.id as UserID,
  JobsT.From_Loc,
  JobsT.mid,
  JobsT.To_Loc,
  ClaimJobT.Create_Date,
  usert.name as userName,
  cmpt.name as cmpname,
  cmpT.id as cmpID,
  ClaimJobT.Approved_date as ApprovedDate,
  ClaimJobT.Finished_Date as FinishDate
from ClaimJobT
  inner join JobsT on JobsT.id = ClaimJobT.JobId
  inner join UserT on UserT.id = ClaimJobT.Driverid
  inner join Cmpt on UserT.belongCMP = cmpt.id
WHERE ClaimJobT.Deleted_date is null
  and UserT.id = $1;
-- name: GetClaimedJobByCmp :many
SELECT ClaimJobT.id as id,
  JobsT.id as JobID,
  UserT.id as UserID,
  JobsT.From_Loc,
  JobsT.mid,
  JobsT.To_Loc,
  ClaimJobT.Create_Date,
  usert.name as userName,
  cmpt.name as cmpname,
  cmpT.id as cmpID,
  ClaimJobT.Approved_date as ApprovedDate,
  ClaimJobT.Finished_Date as FinishDate
from ClaimJobT
  inner join JobsT on JobsT.id = ClaimJobT.JobId
  inner join UserT on UserT.id = ClaimJobT.Driverid
  inner join Cmpt on UserT.belongCMP = cmpt.id
WHERE ClaimJobT.Deleted_date is null
  and UserT.belongCMP = $1;
-- name: GetClaimedJobByID :one
SELECT ClaimJobT.id as id,
  JobsT.id as JobID,
  UserT.id as UserID,
  JobsT.From_Loc,
  finished_date,
  finishPic,
  JobsT.mid,
  JobsT.To_Loc,
  ClaimJobT.Create_Date,
  usert.name as userName,
  cmpt.name as cmpname,
  cmpT.id as cmpID,
  ClaimJobT.Approved_date as ApprovedDate,
  --  DriverT.percentage  as driverPercentage,
  -- ClaimJobT.percentage as percentage, 
  price
from ClaimJobT
  inner join JobsT on JobsT.id = ClaimJobT.JobId
  inner join UserT on UserT.id = ClaimJobT.Driverid
  inner join Cmpt on UserT.belongCMP = cmpt.id
  inner join DriverT on driverT.id = UserT.id
WHERE ClaimJobT.id = $1;
-- name: ClaimJob :one 
INSERT into ClaimJobT (jobID, driverID)
values ($1, $2)
RETURNING id;
-- name: DecreaseRemaining :one
Update JobsT
set remaining = remaining - 1,
  last_modified_date = NOW()
where id = $1
RETURNING remaining;
-- name: DeleteClaimedJob :exec 
Update ClaimJobT
Set deleted_by = $2,
  deleted_date = NOW(),
  last_modified_date = NOW()
where id = $1;
-- name: IncreaseRemaining :one
Update JobsT
set remaining = remaining + 1,
  last_modified_date = NOW()
where id = $1
RETURNING remaining;
-- name: FinishClaimedJob :exec
Update ClaimJobT
Set finishPic = $3,
  finished_date = NOW(),
  -- percentage = (SELECT percentage from driverT where driverT.id = (SELECT driverID from ClaimJobT where ClaimJobT.id = $1)),
  last_modified_date = NOW()
WHERE id = $1
  and ClaimJobT.Driverid = $2;
-- name: ApproveFinishedJob :exec
Update ClaimJobT
set memo = $3,
  Approved_By = $2,
  approved_date = NOW(),
  last_modified_date = NOW()
where id = $1;
-- name: GetCurrentClaimedJob :one
SELECT t2.id as claimID,
  t2.create_date as claimDate,
  t1.from_loc,
  t1.mid,
  t1.to_loc,
  t1.price,
  t1.source,
  t1.memo,
  t1.id
FROM ClaimJobT t2
  inner join JobsT t1 on t1.id = t2.jobID
where t2.driverID = $1
  and (
    t2.deleted_date IS NULL
    and t2.finished_date IS NULL
  )
order by t2.create_date
LIMIT 1;
-- name: GetDriverRevenueByCmp :many
SELECT coalesce(sum(t2.PRICE), 0) as earn,
  coalesce(count(t2.ID), 0) as count
from ClaimJobT t1
  inner join JobsT t2 on t1.jobID = t2.id
  inner join UserT t3 on t1.driverID = t3.id
where t3.belongCMP = $1
  and (
    t1.finished_date IS NOT NULL
    and approved_date IS NOT NULL
    and t1.deleted_date IS NULL
  )
  and date(t1.finished_date) >= date($2)
  and date(t1.finished_date) <= date($3);
-- name: GetDriverRevenue :many
SELECT coalesce(sum(t2.PRICE), 0) as earn,
  coalesce(count(t2.ID), 0) as count
from ClaimJobT t1
  inner join JobsT t2 on t1.jobID = t2.id
where t1.driverID = $1
  and (
    t1.finished_date IS NOT NULL
    and approved_date IS NOT NULL
    and t1.deleted_date IS NULL
  )
  and date(t1.finished_date) >= date($2)
  and date(t1.finished_date) <= date($3);
-- name: CreateNewRepair :one
INSERT into repairT (type, driverID, repairInfo, pic, place)
values ($1, $2, $3, $4, $5)
RETURNING id;
-- name: GetRepair :many
SELECT repairT.id as ID,
  UserT.id as Driverid,
  UserT.Name as Drivername,
  cmpt.name as cmpName,
  repairT.type as type,
  repairT.Repairinfo as Repairinfo,
  repairT.Create_Date as CreateDate,
  repairT.Approved_Date as ApprovedDate,
  repairT.pic as pic,
  repairT.place as place,
  driverT.plateNum as plateNum
from repairT
  inner join UserT on UserT.id = repairT.driverID
  inner join driverT on UserT.id = driverT.id
  inner join cmpt on cmpT.id = UserT.belongCMP
where (
    repairT.id = sqlc.narg('id')
    OR sqlc.narg('id') IS NULL
  )
  AND (
    repairT.driverID = sqlc.narg('driverID')
    OR sqlc.narg('driverID') IS NULL
  )
  AND (
    UserT.name = sqlc.narg('name')
    OR sqlc.narg('name') IS NULL
  )
  AND (
    UserT.belongcmp = sqlc.narg('belongcmp')
    OR sqlc.narg('belongcmp') IS NULL
  ) -- AND (
  --   (
  --     repairT.create_date > sqlc.narg('create_date_start')
  --     OR sqlc.narg('create_date_start') IS NULL
  --   )
  --   AND (
  --     repairT.create_date < sqlc.narg('create_date_end')
  --     OR sqlc.narg('create_date_end') IS NULL
  --   )
  -- )
  AND repairT.deleted_date is null -- AND (
  --   (
  --     repairT.last_modified_date > sqlc.narg('last_modified_date_start')
  --     OR sqlc.narg('last_modified_date_start') IS NULL
  --   )
  --   AND (
  --     repairT.last_modified_date < sqlc.narg('last_modified_date_end')
  --     OR sqlc.narg('last_modified_date_end') IS NULL
  --   )
  -- )
  and (
    (
      sqlc.arg('cat') = 'pending'
      AND repairT.Approved_date IS NULL
    )
    OR (sqlc.narg('cat') IS NULL)
  )
  and (
    to_char(date(repairT.create_date), 'YYYY-MM') = to_char(date(sqlc.narg('ym')), 'YYYY-MM')
    OR sqlc.narg('ym') IS NULL
  );
-- name: GetRepairXXX :many
SELECT repairT.id as ID,
  UserT.id as Driverid,
  UserT.Name as Drivername,
  repairT.type as type,
  repairT.Repairinfo as Repairinfo,
  repairT.Create_Date as CreateDate,
  repairT.Approved_Date as ApprovedDate,
  repairT.pic as pic
from repairT
  inner join UserT on UserT.id = repairT.driverID
where (
    repairT.id = sqlc.narg('id')
    OR sqlc.narg('id') IS NULL
  )
  AND (
    repairT.driverID = sqlc.narg('driverID')
    OR sqlc.narg('driverID') IS NULL
  )
  AND (
    UserT.name = sqlc.narg('name')
    OR sqlc.narg('name') IS NULL
  )
  AND (
    UserT.belongcmp = sqlc.narg('belongcmp')
    OR sqlc.narg('belongcmp') IS NULL
  )
  AND (
    (
      repairT.create_date > sqlc.narg('create_date_start')
      OR sqlc.narg('create_date_start') IS NULL
    )
    AND (
      repairT.create_date < sqlc.narg('create_date_end')
      OR sqlc.narg('create_date_end') IS NULL
    )
  )
  AND repairT.deleted_date is null
  AND (
    (
      repairT.last_modified_date > sqlc.narg('last_modified_date_start')
      OR sqlc.narg('last_modified_date_start') IS NULL
    )
    AND (
      repairT.last_modified_date < sqlc.narg('last_modified_date_end')
      OR sqlc.narg('last_modified_date_end') IS NULL
    )
  );
-- name: ApproveRepair :exec
Update repairT
set approved_date = NOW(),
  last_modified_date = NOW()
where id = $1;
-- name: DeleteRepair :exec
Update repairT
set deleted_date = NOW(),
  last_modified_date = NOW()
where id = $1;
-- name: CreateAlert :one
INSERT INTO AlertT (alert, belongCMP)
values ($1, $2)
RETURNING id;
-- name: UpdateAlert :exec
Update AlertT
Set alert = $2,
  last_modified_date = NOW()
where id = $1;
-- name: DeleteAlert :exec
Update AlertT
Set deleted_date = NOW(),
  last_modified_date = NOW()
where id = $1;
-- name: UpdateLastAlert :exec
Update driverT
set lastAlert = $2
where id = $1;
-- name: GetLastAlert :one
SELECT lastAlert
from driverT
where id = $1;
-- name: GetAlert :many
SELECT alertT.id as ID,
  cmpt.name as cmpName,
  alertT.belongCMP as cmpID,
  alertT.alert as alert,
  alertT.create_date as Createdate,
  alertT.Deleted_Date as Deletedate
from alertT
  inner join cmpt on alertT.Belongcmp = cmpt.id
where (
    alertT.id = sqlc.narg('id')
    OR sqlc.narg('id') IS NULL
  )
  AND (
    belongCMP = sqlc.narg('belongCMP')
    OR sqlc.narg('belongCMP') IS NULL
  )
  AND (
    alert like sqlc.narg('alert')
    OR sqlc.narg('alert') IS NULL
  )
  AND (
    (
      alertT.create_date > sqlc.narg('create_date_start')
      OR sqlc.narg('create_date_start') IS NULL
    )
    AND (
      alertT.create_date < sqlc.narg('create_date_end')
      OR sqlc.narg('create_date_end') IS NULL
    )
  )
  AND (
    (
      alertT.deleted_date > sqlc.narg('deleted_date_start')
      OR sqlc.narg('deleted_date_start') IS NULL
    )
    AND (
      alertT.deleted_date < sqlc.narg('deleted_date_end')
      OR sqlc.narg('deleted_date_end') IS NULL
    )
  )
  AND (
    (
      alertT.last_modified_date > sqlc.narg('last_modified_date_start')
      OR sqlc.narg('last_modified_date_start') IS NULL
    )
    AND (
      alertT.last_modified_date < sqlc.narg('last_modified_date_end')
      OR sqlc.narg('last_modified_date_end') IS NULL
    )
  )
order by alertT.id desc;
-- name: GetAlertByCmp :many
SELECT *
from alertT
where belongCMP = $1
order by id desc;
-- name: GetRepairById :one
SELECT repairT.id as id,
  repairT.driverID as uid,
  repairT.repairInfo,
  repairT.pic,
  repairT.Approved_Date,
  repairT.Create_Date,
  usert.name,
  cmpt.id as cmpid,
  cmpt.name
from repairT
  inner join usert on usert.id = repairT.driverID
  inner join cmpt on usert.belongCMP = cmpt.id
where repairT.id = $1;
-- name: UploadRepairPic :exec
insert into RepairPicT (repair_id, pic)
values ($1, $2);
-- name: ApproveRepairPic :exec
Update RepairPicT
Set Approved_Date = NOW(),
  last_modified_date = NOW()
where repair_id = $1;
-- name: GetRevenueExcel :many
SELECT JSON_BUILD_OBJECT(
    'uid',
    MQ.UID,
    'username',
    MQ.USERNAME,
    'list',
    JSON_AGG(MQ.JSON_BUILD_OBJECT)
  )
FROM (
    SELECT SQ.UID,
      SQ.USERNAME AS USERNAME,
      JSON_BUILD_OBJECT(
        'date',
        SQ.APPROVEDDATE,
        'data',
        JSON_AGG(
          JSON_BUILD_OBJECT(
            -- 'id',
            -- sq.ID,
            'platenum',
            SQ.PLATENUM,
            'cmpName',
            SQ.CMPNAME,
            'fromLoc',
            SQ.FROMLOC,
            'mid',
            SQ.MID,
            'toloc',
            SQ.TOLOC,
            'count',
            SQ.TIMES,
            'jp',
            SQ.JP,
            'total',
            SQ.TOTALPRICE,
            'ss',
            SQ.JOBSOURCE
          )
        ),
        'gas',
        coalesce(SQ.GAS, 0)
      )::JSONB
    FROM (
        SELECT USERT.ID AS UID,
          DRIVERT.PLATENUM AS PLATENUM,
          USERT.NAME AS USERNAME,
          JOBST.BELONGCMP AS BELONGCMP,
          JOBST.FROM_LOC AS FROMLOC,
          coalesce(JOBST.MID, '') AS MID,
          JOBST.TO_LOC AS TOLOC,
          COUNT(*) AS TIMES,
          JOBST.PRICE AS JP,
          JOBST.PRICE * COUNT(*) AS TOTALPRICE,
          JOBST.SOURCE AS JOBSOURCE,
          CMPT.NAME AS CMPNAME,
          DATE (CLAIMJOBT.APPROVED_DATE) AS APPROVEDDATE,
          CLAIMJOBT.MEMO,
          RT.GAS
        FROM CLAIMJOBT
          LEFT JOIN JOBST ON CLAIMJOBT.JOBID = JOBST.ID
          LEFT JOIN USERT ON USERT.ID = CLAIMJOBT.DRIVERID
          LEFT JOIN CMPT ON CMPT.ID = USERT.BELONGCMP
          LEFT JOIN DRIVERT ON USERT.ID = DRIVERT.ID
          LEFT JOIN (
            SELECT DRIVERID,
              SUM(
                ((REPAIRINFO->>0)::JSON->>'price')::INT --* ((REPAIRINFO->>0)::JSON->>'quantity')::INT
              ) AS GAS,
              DATE (CREATE_DATE)
            FROM REPAIRT
            GROUP BY DRIVERID,
              DATE (CREATE_DATE)
          ) RT ON RT.DRIVERID = DRIVERT.ID
          AND DATE (RT.DATE) = DATE (CLAIMJOBT.APPROVED_DATE)
        WHERE (CLAIMJOBT.DELETED_DATE IS NULL)
          AND (CLAIMJOBT.APPROVED_DATE IS NOT NULL)
          and (
            ClaimJobT.Approved_date between $1 and $2
          )
        GROUP BY USERT.ID,
          JOBID,
          DRIVERT.PLATENUM,
          USERT.NAME,
          JOBST.FROM_LOC,
          JOBST.MID,
          JOBST.TO_LOC,
          JOBST.PRICE,
          CMPT.NAME,
          JOBST.BELONGCMP,
          JOBST.SOURCE,
          DATE (CLAIMJOBT.APPROVED_DATE),
          CLAIMJOBT.MEMO,
          RT.GAS
      ) SQ
    GROUP BY SQ.UID,
      SQ.USERNAME,
      SQ.APPROVEDDATE,
      SQ.GAS
  ) MQ
GROUP BY MQ.UID,
  MQ.USERNAME;
-- name: GetCJDate :many
SELECT to_char(create_date, 'YYYY-MM')
FROM public.claimjobt
where driverid = $1
group by to_char(create_date, 'YYYY-MM');
-- name: GetRepairDate :many
SELECT to_char(create_date, 'YYYY-MM')
FROM public.repairT
where driverid = $1
group by to_char(create_date, 'YYYY-MM');