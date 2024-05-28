// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sql

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Alertt struct {
	ID               int64
	Alert            string
	Belongcmp        int64
	CreateDate       time.Time
	DeletedDate      sql.NullTime
	LastModifiedDate time.Time
}

type Claimjobt struct {
	ID               int64
	Jobid            int64
	Driverid         int64
	Percentage       sql.NullInt16
	FinishedDate     sql.NullTime
	Finishpic        sql.NullString
	Memo             sql.NullString
	CreateDate       time.Time
	DeletedDate      sql.NullTime
	DeletedBy        sql.NullInt64
	LastModifiedDate time.Time
	ApprovedBy       sql.NullInt64
	ApprovedDate     sql.NullTime
}

type Cmpincharget struct {
	ID               int64
	Userid           int64
	Cmpid            int64
	CreateDate       time.Time
	DeletedDate      sql.NullTime
	LastModifiedDate time.Time
}

type Cmpt struct {
	ID               int64
	Name             string
	CreateDate       time.Time
	DeletedDate      sql.NullTime
	LastModifiedDate time.Time
}

type Drivert struct {
	ID               int64
	Platenum         string
	Insurances       sql.NullString
	Registration     sql.NullString
	Driverlicense    sql.NullString
	Trucklicense     sql.NullString
	Nationalidnumber interface{}
	Percentage       int16
	Lastalert        sql.NullInt64
	ApprovedDate     sql.NullTime
}

type Jobst struct {
	ID               int64
	FromLoc          string
	Mid              sql.NullString
	ToLoc            string
	Price            int16
	Estimated        int16
	Remaining        int16
	Belongcmp        int64
	Source           string
	Jobdate          time.Time
	Memo             sql.NullString
	CloseDate        sql.NullTime
	CreateDate       time.Time
	DeletedDate      sql.NullTime
	LastModifiedDate time.Time
}

type Logint struct {
	ID         int64
	Userid     int64
	CreateDate time.Time
}

type Repairpict struct {
	ID               int64
	RepairID         int64
	Pic              sql.NullString
	CreateDate       time.Time
	ApprovedDate     sql.NullTime
	DeletedDate      sql.NullTime
	LastModifiedDate time.Time
}

type Repairt struct {
	ID               int64
	Type             string
	Driverid         int64
	Repairinfo       json.RawMessage
	CreateDate       time.Time
	ApprovedDate     sql.NullTime
	DeletedDate      sql.NullTime
	LastModifiedDate time.Time
}

type Usert struct {
	ID               int64
	Phonenum         interface{}
	Pwd              string
	Name             string
	Belongcmp        int64
	Seed             sql.NullString
	Role             int16
	Initpwdchanged   bool
	CreateDate       time.Time
	DeletedDate      sql.NullTime
	LastModifiedDate time.Time
}
