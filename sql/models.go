// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

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
	Userid           pgtype.Int8
	Percentage       int16
	Nationalidnumber interface{}
	CreateDate       pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}

type Logint struct {
	ID         int64
	Userid     pgtype.Int8
	CreateDate pgtype.Timestamp
}

type Usert struct {
	ID               int64
	Username         string
	Pwd              string
	Name             string
	Belongcmp        pgtype.Int8
	Phonenum         interface{}
	Role             int16
	Initpwdchanged   bool
	CreateDate       pgtype.Timestamp
	DeletedDate      pgtype.Timestamp
	LastModifiedDate pgtype.Timestamp
}
