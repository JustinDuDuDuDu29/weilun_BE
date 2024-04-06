package service

import (
	"context"
	"fmt"
	db "main/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserServ interface {
	HaveUser(queryParam db.GetUserParams) (db.GetUserRow, error)
	RegisterCmpAdmin(queryParam db.CreateUserParams) (int64, error)
	RegisterDriver(queryParam db.CreateUserParams, percentage int, nationalIdNumber string) (int64, error)
	DeleteUser(queryParam int64) error
}

type UserServImpl struct {
	q    *db.Queries
	conn *pgx.Conn
}

func (u *UserServImpl) RegisterCmpAdmin(queryParam db.CreateUserParams) (int64, error) {
	res, err := u.q.CreateUser(context.Background(), queryParam)
	return res, err
}

func (u *UserServImpl) RegisterDriver(queryParam db.CreateUserParams, percentage int, nationalIdNumber string) (int64, error) {

	tx, err := u.conn.Begin(context.Background())

	if err != nil {
		return -99, err
	}

	qtx := u.q.WithTx(tx)
	id, err := qtx.CreateUser(context.Background(), queryParam)

	if err != nil {
		fmt.Print("HERE1")
		tx.Rollback(context.Background())
		return -99, err
	}
	driverParam := db.CreateDriverInfoParams{
		Userid:           pgtype.Int8{Int64: id, Valid: true},
		Percentage:       int16(percentage),
		Nationalidnumber: nationalIdNumber,
	}

	userid, err := qtx.CreateDriverInfo(context.Background(), driverParam)

	if err != nil {
		fmt.Print("HERE2")
		tx.Rollback(context.Background())
		return -99, err
	}

	err = tx.Commit(context.Background())

	if err != nil {
		fmt.Print("HERE3")
		tx.Rollback(context.Background())
	}
	return int64(userid.Int64), err
}

func (u *UserServImpl) HaveUser(queryParam db.GetUserParams) (db.GetUserRow, error) {
	res, err := u.q.GetUser(context.Background(), queryParam)
	return res, err
}

func (u *UserServImpl) DeleteUser(queryParam int64) error {
	err := u.q.DeleteUser(context.Background(), queryParam)
	return err
}

func UserServInit(q *db.Queries, conn *pgx.Conn) *UserServImpl {
	return &UserServImpl{
		q:    q,
		conn: conn,
	}
}
