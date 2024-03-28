// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sql

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Cmpincharget struct {
	ID               int64
	Useri            pgtype.Int8
	Cmpid            int64
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
	Name             string
	Phonenum         interface{}
	Belongcmp        pgtype.Int8
	Percentage       int16
	CreateDate       pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}

type Joblistt struct {
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
	DeletedDate      pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}

type Jobtakent struct {
	ID               int64
	Jobid            pgtype.Int8
	Driverid         pgtype.Int8
	Percentage       pgtype.Int2
	IsFinished       bool
	FinishedDate     pgtype.Timestamp
	CreateDate       pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}

type Revenuet struct {
	ID               int64
	Jobtakenid       pgtype.Int8
	Driverearn       int16
	CreateDate       pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}

type Usert struct {
	ID               int64
	Username         string
	Pwd              string
	Role             int16
	CreateDate       pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}
