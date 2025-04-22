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
-- name: GetAllUserFromCMP :many
SELECT usert.id,
  usert.name
from userT
where belongCMP = $1;
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
      order by UT.id
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
-- name: GetJobCmp :many
WITH user_jobs AS (
  SELECT UserT.id AS UserID,
    UserT.name AS UserName,
    Cmpt.id AS CmpID,
    Cmpt.name AS CmpName,
    ClaimJobT.id AS ID,
    JobsT.id AS JobID,
    JobsT.fromLoc AS FromLoc,
    JobsT.mid AS Mid,
    JobsT.toLoc AS ToLoc,
    ClaimJobT.Create_Date AS CreateDate,
    ClaimJobT.Approved_date AS ApprovedDate,
    ClaimJobT.Finished_Date AS FinishDate,
    JobsT.price AS JobPrice
  FROM ClaimJobT
    INNER JOIN JobsT ON JobsT.id = ClaimJobT.JobID
    INNER JOIN UserT ON UserT.id = ClaimJobT.DriverID
    INNER JOIN Cmpt ON UserT.belongCMP = Cmpt.id
  WHERE ClaimJobT.Deleted_date IS NULL
    AND ClaimJobT.Approved_date BETWEEN $1 AND $2
    AND (
      Cmpt.id = sqlc.narg('cmpId')
      OR sqlc.narg('cmpId') IS NULL
    )
    AND (
      UserT.id = sqlc.narg('userId')
      OR sqlc.narg('userId') IS NULL
    )
),
user_gas AS (
  SELECT UserT.id AS UserID,
    UserT.name AS UserName,
    Cmpt.id AS CmpID,
    Cmpt.name AS CmpName,
    SUM(GasInfoT.totalPrice) AS GasTotal
  FROM GasT
    INNER JOIN GasInfoT ON GasInfoT.gasID = GasT.id
    INNER JOIN UserT ON UserT.id = GasT.DriverID
    INNER JOIN Cmpt ON UserT.belongCMP = Cmpt.id
  WHERE GasT.Deleted_date IS NULL
    AND GasT.approved_date BETWEEN $1 AND $2
    AND (
      Cmpt.id = sqlc.narg('cmpId')
      OR sqlc.narg('cmpId') IS NULL
    )
    AND (
      UserT.id = sqlc.narg('userId')
      OR sqlc.narg('userId') IS NULL
    )
  GROUP BY UserT.id,
    UserT.name,
    Cmpt.id,
    Cmpt.name
),
user_repair AS (
  SELECT UserT.id AS UserID,
    UserT.name AS UserName,
    Cmpt.id AS CmpID,
    Cmpt.name AS CmpName,
    SUM(RepairInfoT.totalPrice) AS RepairTotal
  FROM RepairT
    INNER JOIN RepairInfoT ON RepairInfoT.repairID = RepairT.id
    INNER JOIN UserT ON UserT.id = RepairT.DriverID
    INNER JOIN Cmpt ON UserT.belongCMP = Cmpt.id
  WHERE RepairT.Deleted_date IS NULL
    AND RepairT.approved_date BETWEEN $1 AND $2
    AND (
      Cmpt.id = sqlc.narg('cmpId')
      OR sqlc.narg('cmpId') IS NULL
    )
    AND (
      UserT.id = sqlc.narg('userId')
      OR sqlc.narg('userId') IS NULL
    )
  GROUP BY UserT.id,
    UserT.name,
    Cmpt.id,
    Cmpt.name
),
user_summary AS (
  SELECT uj.CmpID,
    uj.CmpName,
    uj.UserID,
    uj.UserName,
    COUNT(uj.ID) AS JobCount,
    SUM(uj.JobPrice) AS JobTotal,
    COALESCE(ug.GasTotal, 0) AS GasTotal,
    COALESCE(ur.RepairTotal, 0) AS RepairTotal
  FROM user_jobs uj
    LEFT JOIN user_gas ug ON uj.UserID = ug.UserID
    LEFT JOIN user_repair ur ON uj.UserID = ur.UserID
  GROUP BY uj.CmpID,
    uj.CmpName,
    uj.UserID,
    uj.UserName,
    ug.GasTotal,
    ur.RepairTotal
)
SELECT cmpt.CmpID AS ID,
  cmpt.CmpName AS Name,
  SUM(cmpt.JobCount) AS Count,
  SUM(cmpt.JobTotal) AS JobTotal,
  SUM(cmpt.GasTotal) AS GasTotal,
  SUM(cmpt.RepairTotal) AS RepairTotal,
  JSON_AGG(
    JSON_BUILD_OBJECT(
      'UserID',
      cmpt.UserID,
      'UserName',
      cmpt.UserName,
      'JobCount',
      cmpt.JobCount,
      'JobTotal',
      cmpt.JobTotal,
      'GasTotal',
      cmpt.GasTotal,
      'RepairTotal',
      cmpt.RepairTotal
    )
  ) AS Users
FROM (
    SELECT CmpID,
      CmpName,
      UserID,
      UserName,
      JobCount,
      JobTotal,
      GasTotal,
      RepairTotal
    FROM user_summary
  ) AS cmpt
GROUP BY cmpt.CmpID,
  cmpt.CmpName;
-- name: GetJobCmpX :many
WITH user_jobs AS (
  SELECT UserT.id AS Userid,
    UserT.name AS Username,
    cmpt.id AS Cmpid,
    cmpt.name AS Cmpname,
    ClaimJobT.id AS ID,
    JobsT.id AS Jobid,
    JobsT.fromLoc AS Fromloc,
    JobsT.mid AS Mid,
    JobsT.toLoc AS Toloc,
    ClaimJobT.Create_Date AS CreateDate,
    ClaimJobT.Approved_date AS Approveddate,
    ClaimJobT.Finished_Date AS Finishdate,
    JobsT.price
  FROM ClaimJobT
    INNER JOIN JobsT ON JobsT.id = ClaimJobT.JobId
    INNER JOIN UserT ON UserT.id = ClaimJobT.Driverid
    INNER JOIN Cmpt ON UserT.belongCMP = Cmpt.id
  WHERE ClaimJobT.Deleted_date IS NULL
    AND ClaimJobT.Approved_date BETWEEN $1 AND $2
    AND (
      cmpt.id = sqlc.narg('cmpId')
      OR sqlc.narg('cmpId') IS NULL
    )
)
SELECT cmpt.cmpID AS ID,
  cmpt.cmpName AS Name,
  SUM(cmpt.jobCount) AS count,
  -- Summing the job count from the inner query
  SUM(cmpt.price) AS total,
  JSON_AGG(
    JSON_BUILD_OBJECT(
      'UserID',
      cmpt.UserID,
      'UserName',
      cmpt.UserName,
      'job',
      jobsList
    )
  ) AS jobs
FROM (
    SELECT cmpID,
      cmpName,
      UserID,
      UserName,
      COUNT(*) AS jobCount,
      -- Count the jobs per user in the inner query
      JSON_AGG(
        JSON_BUILD_OBJECT(
          'ID',
          ID,
          'Username',
          UserName,
          'Jobid',
          JobID,
          'Fromloc',
          Fromloc,
          'Mid',
          Mid,
          'Toloc',
          Toloc,
          'CreateDate',
          CreateDate,
          'Approveddate',
          Approveddate,
          'Finishdate',
          Finishdate
        )
      ) AS jobsList,
      SUM(price) AS price
    FROM user_jobs
    GROUP BY cmpID,
      cmpName,
      UserID,
      UserName
  ) AS cmpt
GROUP BY cmpt.cmpID,
  cmpt.cmpName;
-- name: GetJobCmpXX :many
SELECT cmpt.id AS cmpID,
  cmpt.name AS cmpName,
  COUNT(*) AS jobCount,
  SUM(JobsT.price) AS totalAmount,
  JSON_AGG(
    JSON_BUILD_OBJECT(
      'userID',
      UserT.id,
      'userName',
      UserT.name,
      'jobs',
      (
        SELECT JSON_AGG(
            JSON_BUILD_OBJECT(
              'id',
              ClaimJobT.id,
              'jobID',
              JobsT.id,
              'fromLoc',
              JobsT.fromLoc,
              'mid',
              JobsT.mid,
              'toLoc',
              JobsT.toLoc,
              'createDate',
              ClaimJobT.Create_Date,
              'approvedDate',
              ClaimJobT.Approved_date,
              'finishDate',
              ClaimJobT.Finished_Date
            )
          )
        FROM ClaimJobT AS innerClaim
        WHERE innerClaim.DriverID = UserT.id
          AND innerClaim.Deleted_date IS NULL
          AND innerClaim.Approved_date BETWEEN $1 AND $2
      )
    )
  ) AS users
FROM ClaimJobT
  INNER JOIN JobsT ON JobsT.id = ClaimJobT.JobId
  INNER JOIN UserT ON UserT.id = ClaimJobT.Driverid
  INNER JOIN Cmpt ON UserT.belongCMP = Cmpt.id
WHERE ClaimJobT.Deleted_date IS NULL
  AND (
    cmpt.id = sqlc.narg('cmpId')
    OR sqlc.narg('cmpId') IS NULL
  )
  AND ClaimJobT.Approved_date BETWEEN $1 AND $2
GROUP BY cmpt.id,
  cmpt.name;
-- name: GetJobCmpXXX :many
SELECT cmpt.id,
  cmpt.name,
  COUNT(*) as count,
  sum(price) as total
FROM cmpt
  LEFT JOIN userT ON userT.belongCMP = cmpt.id
  LEFT JOIN claimjobt ON claimjobt.driverID = userT.id
  LEFT JOIN JobsT on claimjobt.jobID = JobsT.id
WHERE claimjobt.Approved_date BETWEEN $1 AND $2
  AND (
    cmpt.id = sqlc.narg('cmpId')
    OR sqlc.narg('cmpId') IS NULL
  )
GROUP BY cmpt.id,
  cmpt.name;
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
  JobsT.fromLoc,
  JobsT.Mid,
  JobsT.toLoc,
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
    JobsT.fromLoc = sqlc.narg('FromLoc')
    OR sqlc.narg('FromLoc') IS NULL
  )
  AND (
    JobsT.Mid = sqlc.narg('Mid')
    OR sqlc.narg('Mid') IS NULL
  )
  AND (
    JobsT.toLoc = sqlc.narg('ToLoc')
    OR sqlc.narg('ToLoc') IS NULL
  )
  AND (
    belongcmp = sqlc.narg('belongcmp')
    OR sqlc.narg('belongcmp') IS NULL
  )
  AND (remaining <> 0)
  AND (JobsT.deleted_date is NULL);
-- name: GetAllJobsSuper :many
SELECT cmpt.name AS cmpName,
  cmpt.id AS BELONGCMP,
  json_agg(
    json_build_object(
      'ID',
      JobsT.id,
      'Fromloc',
      JobsT.fromLoc,
      'Toloc',
      JobsT.toLoc,
      'Mid',
      JobsT.mid,
      'Price',
      JobsT.price,
      'estimated',
      JobsT.estimated,
      'Remaining',
      JobsT.remaining,
      'Source',
      JobsT.source,
      'Memo',
      JobsT.memo,
      'Belongcmp',
      cmpT.id,
      'create_date',
      JobsT.create_date
    )
  ) AS jobs
FROM JobsT
  INNER JOIN cmpt ON JobsT.belongcmp = cmpt.id
where (
    JobsT.id = sqlc.narg('id')
    OR sqlc.narg('id') IS NULL
  )
  AND (
    JobsT.fromLoc = sqlc.narg('FromLoc')
    OR sqlc.narg('FromLoc') IS NULL
  )
  AND (
    JobsT.Mid = sqlc.narg('Mid')
    OR sqlc.narg('Mid') IS NULL
  )
  AND (
    JobsT.toLoc = sqlc.narg('ToLoc')
    OR sqlc.narg('ToLoc') IS NULL
  )
  AND (
    belongcmp = sqlc.narg('belongcmp')
    OR sqlc.narg('belongcmp') IS NULL
  )
  AND (remaining <> 0)
  AND (
    remaining <> sqlc.narg('remaining')
    OR sqlc.narg('remaining') IS NULL
  )
  AND (
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
  AND(JobsT.deleted_date is NULL)
GROUP BY cmpt.name,
  cmpt.id;
-- name: GetAllJobsAdmin :many
SELECT *,
  cmpt.name as cmpName
from JobsT
  inner join cmpt on JobsT.belongcmp = cmpt.id
where (
    JobsT.id = sqlc.narg('id')
    OR sqlc.narg('id') IS NULL
  )
  AND (
    JobsT.fromLoc = sqlc.narg('FromLoc')
    OR sqlc.narg('FromLoc') IS NULL
  )
  AND (
    JobsT.Mid = sqlc.narg('Mid')
    OR sqlc.narg('Mid') IS NULL
  )
  AND (
    JobsT.toLoc = sqlc.narg('ToLoc')
    OR sqlc.narg('ToLoc') IS NULL
  )
  AND (
    belongcmp = sqlc.narg('belongcmp')
    OR sqlc.narg('belongcmp') IS NULL
  )
  AND (remaining <> 0)
  AND (
    remaining <> sqlc.narg('remaining')
    OR sqlc.narg('remaining') IS NULL
  )
  AND (
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
set fromLoc = $1,
  mid = $2,
  toLoc = $3,
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
    fromLoc,
    mid,
    toLoc,
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
SELECT ClaimJobT.id AS id,
  JobsT.id AS JobID,
  UserT.id AS UserID,
  JobsT.fromLoc AS fromloc,
  -- Aliased as fromloc
  JobsT.mid AS mid,
  JobsT.toLoc AS toloc,
  -- Aliased as toloc
  JobsT.Price,
  ClaimJobT.Create_Date,
  UserT.name AS userName,
  Cmpt.name AS cmpname,
  Cmpt.id AS cmpID,
  ClaimJobT.Approved_date AS ApprovedDate,
  ClaimJobT.Finished_Date AS FinishDate,
  ClaimJobT.finishPic
FROM ClaimJobT
  INNER JOIN JobsT ON JobsT.id = ClaimJobT.JobId
  INNER JOIN UserT ON UserT.id = ClaimJobT.Driverid
  INNER JOIN Cmpt ON UserT.belongCMP = Cmpt.id
WHERE ClaimJobT.Deleted_date IS NULL
  AND (
    ClaimJobT.driverid = sqlc.narg('uid')
    OR sqlc.narg('uid') IS NULL
  )
  AND (
    ClaimJobT.jobID = sqlc.narg('jobid')
    OR sqlc.narg('jobid') IS NULL
  )
  AND (
    UserT.belongCMP = sqlc.narg('cmpID')
    OR sqlc.narg('cmpID') IS NULL
  )
  AND (
    ClaimJobT.id = sqlc.narg('cjID')
    OR sqlc.narg('cjID') IS NULL
  )
  AND (
    (
      sqlc.narg('cat') = 'pending'
      AND ClaimJobT.Approved_date IS NULL
    )
    OR sqlc.narg('cat') IS NULL
  )
  AND (
    TO_CHAR(DATE(ClaimJobT.create_date), 'YYYY-MM') = TO_CHAR(DATE(sqlc.narg('ym')), 'YYYY-MM')
    OR sqlc.narg('ym') IS NULL
  )
  AND ClaimJobT.deleted_date IS NULL;
-- name: GetClaimedJobByDriverID :many
SELECT ClaimJobT.id as id,
  JobsT.id as JobID,
  UserT.id as UserID,
  JobsT.fromLoc,
  JobsT.mid,
  JobsT.toLoc,
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
  and UserT.id = $1
  and ClaimJobT.Approved_date between $2 and $3;
-- name: GetClaimedJobByCmp :many
SELECT ClaimJobT.id as id,
  JobsT.id as JobID,
  UserT.id as UserID,
  JobsT.fromLoc,
  JobsT.mid,
  JobsT.toLoc,
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
  and UserT.belongCMP = $1
  and ClaimJobT.Approved_date between $2 and $3;
-- name: GetClaimedJobByID :one
SELECT ClaimJobT.id as id,
  JobsT.id as JobID,
  UserT.id as UserID,
  JobsT.fromLoc,
  finished_date,
  finishPic,
  JobsT.mid,
  JobsT.toLoc,
  ClaimJobT.Create_Date,
  usert.name as userName,
  cmpt.name as cmpname,
  cmpT.id as cmpID,
  ClaimJobT.Approved_date as ApprovedDate,
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

-- name: ApproveMultipleJobs :exec
UPDATE ClaimJobT
SET 
  memo = sqlc.arg(memo),
  approved_by = sqlc.arg(approved_by),
  approved_date = NOW(),
  last_modified_date = NOW()
WHERE id = ANY(sqlc.arg(ids)::bigint[]);

-- name: GetCurrentClaimedJob :one
SELECT t2.id as claimID,
  t2.create_date as claimDate,
  t1.fromLoc,
  t1.mid,
  t1.toLoc,
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
-- name: CreateNewRepair :one
INSERT into repairT (driverID, pic, place)
values ($1, $2, $3)
RETURNING id;
-- name: CreateNewRepairInfo :one
INSERT into repairInfoT (repairID, itemName, quantity, totalPrice)
values ($1, $2, $3, $4)
RETURNING id;
-- name: GetRepair :many
SELECT repairT.id as ID,
  UserT.id as Driverid,
  UserT.Name as Drivername,
  cmpt.name as cmpName,
  cmpt.id as cmpID,
  -- Include repair information as JSON
  (
    SELECT json_agg(
        json_build_object(
          'id',
          repairinfoT.id,
          'itemName',
          repairinfoT.itemName,
          'quantity',
          repairinfoT.quantity,
          'totalPrice',
          repairinfoT.totalPrice,
          'create_date',
          repairinfoT.create_date
        )
      )
    FROM repairinfoT
    WHERE repairinfoT.repairID = repairT.id
  ) as Repairinfo,
  repairT.Create_Date as CreateDate,
  repairT.Approved_Date as ApprovedDate,
  repairT.pic as pic,
  repairT.place as place,
  driverT.plateNum as plateNum
FROM repairT
  INNER JOIN UserT ON UserT.id = repairT.driverID
  INNER JOIN driverT ON UserT.id = driverT.id
  INNER JOIN cmpt ON cmpt.id = UserT.belongCMP
WHERE (
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
  AND repairT.deleted_date IS NULL
  AND (
    (
      sqlc.arg('cat') = 'pending'
      AND repairT.Approved_date IS NULL
    )
    OR (sqlc.narg('cat') IS NULL)
  )
  AND (
    to_char(date(repairT.create_date), 'YYYY-MM') = to_char(date(sqlc.narg('ym')), 'YYYY-MM')
    OR sqlc.narg('ym') IS NULL
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
-- name: GetRepairDate :many
SELECT to_char(create_date, 'YYYY-MM')
FROM public.repairT
where driverid = $1
group by to_char(create_date, 'YYYY-MM');
-- name: GetRepairInfoById :many
SELECT *
from RepairInfoT
where repairID = $1;

-- name: GetRevenueExcel :many
WITH GasData AS (
  SELECT GasT.DRIVERID,
    SUM(GasInfoT.totalPrice) AS GAS,
    DATE(GasT.Approved_Date) AS GAS_DATE
  FROM GasT
    LEFT JOIN GasInfoT ON GasInfoT.gasid = GasT.id
  where GasT.approved_date is not null and GasT.deleted_date is null
  And GasT.Approved_Date  BETWEEN $1 AND $2
  GROUP BY GasT.DRIVERID,
    DATE(GasT.Approved_Date)
),
RepairData AS (
  SELECT RepairT.DRIVERID,
    SUM(RepairInfoT.totalPrice) AS REPAIR,
    DATE(RepairT.Approved_Date) AS REPAIR_DATE
  FROM RepairT
    LEFT JOIN RepairInfoT ON RepairInfoT.repairid = RepairT.id
  where RepairT.approved_date is not null and RepairT.deleted_date is null
  and RepairT.Approved_Date  BETWEEN $1 AND $2
  GROUP BY RepairT.DRIVERID,
    DATE(RepairT.Approved_Date)
),
JobData AS (
  SELECT USERT.ID AS UID,
    USERT.NAME AS USERNAME,
    DRIVERT.PLATENUM AS PLATENUM,
    JOBST.fromLoc AS FROMLOC,
    COALESCE(JOBST.MID, '') AS MID,
    JOBST.toLoc AS TOLOC,
    COUNT(CLAIMJOBT.JOBID) AS TIMES,
    JOBST.PRICE AS JP,
    JOBST.PRICE * COUNT(CLAIMJOBT.JOBID) AS TOTALPRICE,
    JOBST.SOURCE AS JOBSOURCE,
    CMPT.NAME AS CMPNAME,
    DATE(CLAIMJOBT.finished_date) AS finishDate
  FROM CLAIMJOBT
    LEFT JOIN JOBST ON CLAIMJOBT.JOBID = JOBST.ID
    LEFT JOIN USERT ON USERT.ID = CLAIMJOBT.DRIVERID
    LEFT JOIN DRIVERT ON USERT.ID = DRIVERT.ID
    LEFT JOIN CMPT ON CMPT.ID = USERT.BELONGCMP
  WHERE CLAIMJOBT.DELETED_DATE IS NULL
    AND CLAIMJOBT.finished_date BETWEEN $1 AND $2
    AND USERT.belongCMP = $3
  GROUP BY USERT.ID,
    USERT.NAME,
    DRIVERT.PLATENUM,
    JOBST.fromLoc,
    JOBST.MID,
    JOBST.toLoc,
    JOBST.PRICE,
    JOBST.SOURCE,
    CMPT.NAME,
    DATE(CLAIMJOBT.finished_date)
)
SELECT JSON_BUILD_OBJECT(
    'uid',
    MQ.UID,
    'username',
    MQ.USERNAME,
    'list',
    JSON_AGG(
      MQ.JSON_BUILD_OBJECT
      ORDER BY MQ.DATE ASC
    )
  )
FROM (
    SELECT COALESCE(JD.UID, GD.DRIVERID, RD.DRIVERID) AS UID,
      COALESCE(JD.USERNAME, U.NAME) AS USERNAME,
      COALESCE(JD.finishDate, GD.GAS_DATE, RD.REPAIR_DATE) AS DATE,
      JSON_BUILD_OBJECT(
        'date',
        COALESCE(JD.finishDate, GD.GAS_DATE, RD.REPAIR_DATE),
        'data',
        JSON_AGG(
          JSON_BUILD_OBJECT(
            'platenum',
            COALESCE(JD.PLATENUM, DR.PLATENUM, 'Unknown Plate'),
            'cmpName',
            JD.CMPNAME,
            'fromLoc',
            JD.FROMLOC,
            'mid',
            JD.MID,
            'toloc',
            JD.TOLOC,
            'count',
            JD.TIMES,
            'jp',
            JD.JP,
            'total',
            JD.TOTALPRICE,
            'ss',
            JD.JOBSOURCE
          )
        ),
        'gas',
        JSON_BUILD_OBJECT(
          'platenum',
          COALESCE(JD.PLATENUM, DR.PLATENUM, 'Unknown Plate'),
          'gasAmount',
          COALESCE(GD.GAS, 0)
        ),
        'repair',
        JSON_BUILD_OBJECT(
          'platenum',
          COALESCE(JD.PLATENUM, DR.PLATENUM, 'Unknown Plate'),
          'repairAmount',
          COALESCE(RD.REPAIR, 0)
        )
      ) AS JSON_BUILD_OBJECT
    FROM JobData JD
      FULL OUTER JOIN GasData GD ON JD.UID = GD.DRIVERID
      AND JD.finishDate = GD.GAS_DATE
      FULL OUTER JOIN RepairData RD ON COALESCE(JD.UID, GD.DRIVERID) = RD.DRIVERID
      AND COALESCE(JD.finishDate, GD.GAS_DATE) = RD.REPAIR_DATE
      LEFT JOIN USERT U ON COALESCE(JD.UID, GD.DRIVERID, RD.DRIVERID) = U.ID
      LEFT JOIN DRIVERT DR ON COALESCE(JD.UID, GD.DRIVERID, RD.DRIVERID) = DR.ID
    GROUP BY GD.GAS,
      RD.REPAIR,
      COALESCE(JD.UID, GD.DRIVERID, RD.DRIVERID),
      COALESCE(JD.USERNAME, U.NAME),
      COALESCE(JD.finishDate, GD.GAS_DATE, RD.REPAIR_DATE),
      JD.PLATENUM,
      DR.PLATENUM
  ) MQ
GROUP BY MQ.UID,
  MQ.USERNAME
ORDER BY MAX(MQ.DATE) ASC;

-- name: GetCJDate :many
SELECT to_char(create_date, 'YYYY-MM')
FROM public.claimjobt
where driverid = $1
group by to_char(create_date, 'YYYY-MM');
-- name: GetUserWithPendingJob :many
SELECT JSON_BUILD_OBJECT(
    'cmpID',
    CMPT.id,
    'cmpName',
    CMPT.name,
    'users',
    JSON_AGG(
      JSON_BUILD_OBJECT(
        'userID',
        UserT.id,
        'userName',
        UserT.name,
        'jobs',
        (
          SELECT JSON_AGG(
              JSON_BUILD_OBJECT(
                'ID',
                CLAIMJOBT.id,
                'Userid',
                CLAIMJOBT.driverid,
                'CreateDate',
                CLAIMJOBT.create_date,
                'userName',
                UserT.name,
                'Cmpname',
                CMPT.name,
                'Finishdate',
                CLAIMJOBT.finished_date,
                'Cmpid',
                CMPT.id,
                'Approveddate',
                CLAIMJOBT.approved_date,
                'jobId',
                CLAIMJOBT.jobid,
                'Fromloc',
                jobst.fromloc,
                'Mid',
                jobst.mid,
                'Toloc',
                jobst.toloc
              )
            )
          FROM CLAIMJOBT
            LEFT JOIN jobst ON jobst.id = CLAIMJOBT.jobid
          WHERE CLAIMJOBT.driverID = UserT.id
            AND CLAIMJOBT.approved_date IS NULL
            AND CLAIMJOBT.deleted_date IS NULL
        )
      )
    )
  ) AS result
FROM UserT
  LEFT JOIN CMPT ON CMPT.id = UserT.belongCMP
WHERE EXISTS (
    SELECT 1
    FROM CLAIMJOBT
    WHERE CLAIMJOBT.driverID = UserT.id
      AND CLAIMJOBT.approved_date IS NULL
      AND CLAIMJOBT.deleted_date IS NULL
  )
  AND (
    CMPT.id = sqlc.narg('cmpid')
    OR sqlc.narg('cmpid') IS NULL
  )
GROUP BY CMPT.id,
  CMPT.name;
  
-- name: CreateNewGas :one
INSERT into GasT (driverID, pic)
values ($1, $2)
RETURNING id;
-- name: CreateNewGasInfo :one
INSERT into GasInfoT (gasID, itemName, quantity, totalPrice)
values ($1, $2, $3, $4)
RETURNING id;
-- name: GetGas :many
SELECT gasT.id as ID,
  UserT.id as Driverid,
  UserT.Name as Drivername,
  cmpt.name as cmpName,
  cmpt.id as cmpID,
  -- Include repair information as JSON with default values if NULL

  COALESCE(
    (
      SELECT jsonb_agg(
        jsonb_build_object(
          'id', COALESCE(GasInfoT.id, 0),
          'itemName', COALESCE(GasInfoT.itemName, ''),
          'quantity', COALESCE(GasInfoT.quantity, 0),
          'totalPrice', COALESCE(GasInfoT.totalPrice, 0),
          'create_date', COALESCE(GasInfoT.create_date, '1970-01-01')
        )
      )
      FROM GasInfoT
      WHERE GasInfoT.gasID = gasT.id
    ), '[
        {
          "id": 0,
          "itemName": "",
          "quantity": 0,
          "totalPrice": 0,
          "create_date": "1970-01-01"
        }
      ]'::jsonb
  ) AS Repairinfo,
  gasT.Create_Date as CreateDate,
  gasT.Approved_Date as ApprovedDate,
  gasT.pic as pic,
  driverT.plateNum as plateNum
FROM gasT
  INNER JOIN UserT ON UserT.id = gasT.driverID
  INNER JOIN driverT ON UserT.id = driverT.id
  INNER JOIN cmpt ON cmpt.id = UserT.belongCMP
WHERE  (
    gasT.id = sqlc.narg('id')
    OR sqlc.narg('id') IS NULL
  )
  AND (
    gasT.driverID = sqlc.narg('driverID')
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
  AND gasT.deleted_date IS NULL
  AND (
    (
      sqlc.arg('cat') = 'pending'
      AND gasT.Approved_date IS NULL
    )
    OR (sqlc.narg('cat') IS NULL)
  )
  AND (
    to_char(date(gasT.create_date), 'YYYY-MM') = to_char(date(sqlc.narg('ym')), 'YYYY-MM')
    OR sqlc.narg('ym') IS NULL
  );

-- name: ApproveGas :exec
Update GasT
set approved_date = NOW(),
  last_modified_date = NOW()
where id = $1;
-- name: DeleteGasT :exec
Update gasT
set deleted_date = NOW(),
  last_modified_date = NOW()
where id = $1;
-- name: GetGasDate :many
SELECT to_char(create_date, 'YYYY-MM')
FROM public.GasT
where driverid = $1
group by to_char(create_date, 'YYYY-MM');
-- name: GetGasInfoById :many
SELECT *
from GasInfoT
where gasID = $1;
-- name: UpdateItem :exec
UPDATE RepairInfoT
SET totalPrice = $2,
  last_modified_date = NOW()
WHERE RepairInfoT.id = $1;
-- name: GetRepairCmpUser :many
SELECT JSON_BUILD_OBJECT(
    'cmpName',
    cmpt.name,
    'cmpId',
    cmpt.id,
    'users',
    JSON_AGG(
      JSON_BUILD_OBJECT(
        'driverID',
        UserT.id,
        'driverName',
        UserT.name
      )
    )
  ) AS result
FROM cmpt
  INNER JOIN (
    SELECT DISTINCT ON (UserT.id) UserT.id,
      UserT.name,
      UserT.belongCMP
    FROM repairT
      INNER JOIN UserT ON UserT.id = repairT.driverID
      INNER JOIN driverT ON UserT.id = driverT.id
    WHERE repairT.deleted_date IS NULL
      AND ((repairT.Approved_date IS NULL))
  ) AS UserT ON UserT.belongCMP = cmpt.id
WHERE (
    UserT.belongCMP = sqlc.narg('belongcmp')
    OR sqlc.narg('belongcmp') IS NULL
  )
GROUP BY cmpt.name,
  cmpt.id;
-- name: GetGasCmpUser :many
SELECT JSON_BUILD_OBJECT(
    'cmpName',
    cmpt.name,
    'cmpId',
    cmpt.id,
    'users',
    JSON_AGG(
      JSON_BUILD_OBJECT(
        'driverID',
        UserT.id,
        'driverName',
        UserT.name
      )
    )
  ) AS result
FROM cmpt
  INNER JOIN (
    SELECT DISTINCT ON (UserT.id) UserT.id,
      UserT.name,
      UserT.belongCMP
    FROM GasT
      INNER JOIN UserT ON UserT.id = GasT.driverID
      INNER JOIN driverT ON UserT.id = driverT.id
    WHERE GasT.deleted_date IS NULL
      AND ((GasT.Approved_date IS NULL))
  ) AS UserT ON UserT.belongCMP = cmpt.id
WHERE (
    UserT.belongCMP = sqlc.narg('belongcmp')
    OR sqlc.narg('belongcmp') IS NULL
  )
GROUP BY cmpt.name,
  cmpt.id;
-- name: UpdateGas :exec
UPDATE GasInfoT
SET totalPrice = $2,
  last_modified_date = NOW()
WHERE GasInfoT.id = $1;