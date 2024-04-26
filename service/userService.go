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
	UpdateUser(param db.UpdateUserParams) error
	UpdateDriver(param db.UpdateDriverParams, userParam db.UpdateUserParams) error
	ApproveDriver(id int64) error
	UpdateDriverPic(param db.UpdateDriverPicParams) error
	UpdatePassword(param db.UpdateUserPasswordParams) error
}

type UserServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (u *UserServImpl) UpdatePassword(param db.UpdateUserPasswordParams) error {
	err := u.q.UpdateUserPassword(context.Background(), param)
	return err
}

func (u *UserServImpl) ApproveDriver(id int64) error {
	err := u.q.ApproveDriver(context.Background(), id)
	return err
}

func (u *UserServImpl) UpdateDriverPic(param db.UpdateDriverPicParams) error {
	err := u.q.UpdateDriverPic(context.Background(), param)
	return err
}

func (u *UserServImpl) UpdateDriver(param db.UpdateDriverParams, userParam db.UpdateUserParams) error {
	tx, err := u.conn.BeginTx(context.Background(), nil)

	if err != nil {
		return err
	}

	qtx := u.q.WithTx(tx)
	err = qtx.UpdateDriver(context.Background(), param)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = qtx.UpdateUser(context.Background(), userParam)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (u *UserServImpl) UpdateUser(param db.UpdateUserParams) error {
	err := u.q.UpdateUser(context.Background(), param)
	return err
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
		return -99, err
	}
	return int64(userid), nil
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
