package service

import (
	"context"
	"database/sql"
	"fmt"
	db "main/sql"
	// "github.com/jackc/pgx/v5"
)

type UserServ interface {
	HaveUser(queryParam db.GetUserParams) (db.GetUserRow, error)
	GetUserById(id int64) (db.GetUserByIDRow, error)
	RegisterCmpAdmin(queryParam db.CreateUserParams) (int64, error)
	RegisterDriver(queryParam db.CreateUserParams, percentage int, nationalIdNumber string) (int64, error)
	DeleteUser(queryParam int64) error
	GetUserList(queryParam db.GetUserListParams) ([]db.GetUserListRow, error)
}

type UserServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (u *UserServImpl) GetUserList(queryParam db.GetUserListParams) ([]db.GetUserListRow, error) {
	res, err := u.q.GetUserList(context.Background(), queryParam)
	return res, err
}

func (u *UserServImpl) RegisterCmpAdmin(queryParam db.CreateUserParams) (int64, error) {
	res, err := u.q.CreateUser(context.Background(), queryParam)
	return res, err
}

func (u *UserServImpl) RegisterDriver(queryParam db.CreateUserParams, percentage int, nationalIdNumber string) (int64, error) {

	tx, err := u.conn.BeginTx(context.Background(), nil)

	if err != nil {
		return -99, err
	}

	qtx := u.q.WithTx(tx)
	id, err := qtx.CreateUser(context.Background(), queryParam)

	if err != nil {
		fmt.Print("HERE1")
		tx.Rollback()
		return -99, err
	}
	driverParam := db.CreateDriverInfoParams{
		ID:               id,
		Percentage:       int16(percentage),
		Nationalidnumber: nationalIdNumber,
	}

	userid, err := qtx.CreateDriverInfo(context.Background(), driverParam)

	if err != nil {
		tx.Rollback()
		return -99, err
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
	}
	return int64(userid), err
}

func (u *UserServImpl) GetUserById(id int64) (db.GetUserByIDRow, error) {
	res, err := u.q.GetUserByID(context.Background(), id)
	return res, err
}

func (u *UserServImpl) HaveUser(queryParam db.GetUserParams) (db.GetUserRow, error) {
	res, err := u.q.GetUser(context.Background(), queryParam)
	return res, err
}

func (u *UserServImpl) DeleteUser(queryParam int64) error {
	err := u.q.DeleteUser(context.Background(), queryParam)
	return err
}

func UserServInit(q *db.Queries, conn *sql.DB) *UserServImpl {
	return &UserServImpl{
		q:    q,
		conn: conn,
	}
}
