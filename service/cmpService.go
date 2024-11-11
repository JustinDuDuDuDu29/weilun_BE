package service

import (
	"context"
	"database/sql"
	db "main/sql"
)

type CmpServ interface {
	GetCmp(cmpId int64) (db.GetCmpRow, error)
	NewCmp(name string) (int64, error)
	GetAllCmp() ([]db.Cmpt, error)
	GetJobCmp(queryParam db.GetJobCmpParams) ([]db.GetJobCmpRow, error)
	UpdateCmp(queryParam db.UpdateCmpParams) error
	DeleteCmp(queryParam int64) error
}

type CmpServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (u *CmpServImpl) GetCmp(cmpId int64) (db.GetCmpRow, error) {
	res, err := u.q.GetCmp(context.Background(), cmpId)
	return res, err
}

func (u *CmpServImpl) GetJobCmp(queryParam db.GetJobCmpParams) ([]db.GetJobCmpRow, error) {
	res, err := u.q.GetJobCmp(context.Background(), queryParam)
	return res, err
}

func (u *CmpServImpl) GetAllCmp() ([]db.Cmpt, error) {
	res, err := u.q.GetAllCmp(context.Background())
	return res, err
}

func (u *CmpServImpl) NewCmp(name string) (int64, error) {
	id, err := u.q.NewCmp(context.Background(), name)
	return id, err
}

func (u *CmpServImpl) UpdateCmp(queryParam db.UpdateCmpParams) error {

	err := u.q.UpdateCmp(context.Background(), queryParam)
	return err
}

func (u *CmpServImpl) DeleteCmp(queryParam int64) error {
	err := u.q.DeleteCmp(context.Background(), queryParam)
	return err
}

func CmpServInit(q *db.Queries, conn *sql.DB) *CmpServImpl {
	return &CmpServImpl{
		q:    q,
		conn: conn,
	}
}
