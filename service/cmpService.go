package service

import (
	"context"
	"fmt"
	db "main/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CmpServ interface {
	GetCmp(queryParam db.GetUserParams) (db.GetUserRow, error)
	NewCmp(queryParam db.CreateUserParams) (int64, error)
	UpdateCmp(queryParam int64) error
	DeleteCmp(queryParam int64) error
}

type CmpServImpl struct {
	q    *db.Queries
	conn *pgx.Conn
}

func (u *CmpServImpl) GetCmp(cmpId int64) (db.GetCmpRow, error) {
	res, err := u.q.GetCmp(context.Background(), cmpId)
	return res, err
}

func (u *CmpServImpl) NewCmp(name string) (int64, error) {
	id, err := u.q.NewCmp(context.Background(), name)

	return id, err
}

func (u *CmpServImpl) HaveUser(queryParam db.GetUserParams) (db.GetUserRow, error) {
	res, err := u.q.GetUser(context.Background(), queryParam)
	return res, err
}

func (u *CmpServImpl) DeleteUser(queryParam int64) error {
	err := u.q.DeleteUser(context.Background(), queryParam)
	return err
}

func CmpServInit(q *db.Queries, conn *pgx.Conn) *CmpServImpl {
	return &CmpServImpl{
		q:    q,
		conn: conn,
	}
}
