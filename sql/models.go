// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Claimjobt struct {
	ID               int64
	Jobid            int64
	Driverid         int64
	Percentage       pgtype.Int2
	FinishedDate     pgtype.Timestamp
	CreateDate       pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	DeletedBy        pgtype.Int8
	LastModifiedDate pgtype.Timestamp
	ApprovedBy       pgtype.Int8
	ApprovedDate     pgtype.Timestamp
}

type Cmpincharget struct {
	ID               int64
	Userid           pgtype.Int8
	Cmpid            pgtype.Int8
	CreateDate       pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}

type Cmpt struct {
	ID               int64
	Name             string
	CreateDate       pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}

type Drivert struct {
	ID               int64
	Percentage       int16
	Nationalidnumber interface{}
}

type Jobst struct {
	ID               int64
	FromLoc          string
	Mid              pgtype.Text
	ToLoc            string
	Price            int16
	Estimated        int16
	Remaining        int16
	Belongcmp        int64
	Source           string
	Jobdate          pgtype.Timestamp
	Memo             pgtype.Text
	CreateDate       pgtype.Timestamp
	EndDate          pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	FinishedDate     pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}

type Logint struct {
	ID         int64
	Userid     pgtype.Int8
	CreateDate pgtype.Timestamp
}

type Usert struct {
	ID               int64
	Phonenum         interface{}
	Pwd              string
	Name             string
	Belongcmp        pgtype.Int8
	Role             int16
	Initpwdchanged   bool
	CreateDate       pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}
