package service

import (
	"context"
	"fmt"
	"log"
	db "main/sql"
)

type UserServ interface {
	HaveUser(queryParam db.GetUserParams)
	RegisterCmpAdmin(queryParam db.CreateCmpAdminParams)
	RegisterDriver(queryParam db.CreateDriverParams)
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

func (u *UserServImpl) HaveUser(queryParam db.GetUserParams) {
	res, err := u.q.GetUser(context.Background(), queryParam)
	if err != nil {

		log.Fatalf("Err!! %s", err.Error())

	}
	fmt.Println(res)
}

func UserServInit(q *db.Queries) *UserServImpl {
	return &UserServImpl{
		q: q,
	}
}
