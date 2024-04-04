package service

import (
	"context"
	"fmt"
	"log"
	db "main/sql"
)

type UserServ interface {
	HaveUser(queryParam db.GetUserParams) (db.GetUserRow, error)
	RegisterCmpAdmin(queryParam db.CreateCmpAdminParams)
	RegisterDriver(queryParam db.CreateDriverParams)
	DeleteUser(queryParam int64) error
}

type UserServImpl struct {
	q *db.Queries
}

func (u *UserServImpl) RegisterCmpAdmin(queryParam db.CreateCmpAdminParams) {
	res, err := u.q.CreateCmpAdmin(context.Background(), queryParam)
	if err != nil {

		log.Fatalf("Err!! %s", err.Error())

	}
	fmt.Println(res)
}

func (u *UserServImpl) RegisterDriver(queryParam db.CreateDriverParams) {
	res, err := u.q.CreateDriver(context.Background(), queryParam)
	if err != nil {
		log.Fatalf("Err!! %s", err.Error())
	}
	fmt.Println(res)
}

func (u *UserServImpl) HaveUser(queryParam db.GetUserParams) (db.GetUserRow, error) {
	res, err := u.q.GetUser(context.Background(), queryParam)
	fmt.Println(res)
	return res, err
}

func (u *UserServImpl) DeleteUser(queryParam int64) error {
	err := u.q.DeleteUser(context.Background(), queryParam)

	return err
}

func UserServInit(q *db.Queries) *UserServImpl {
	return &UserServImpl{
		q: q,
	}
}
