package service

import (
	"context"
	"fmt"
	"log"
	db "main/sql"
)

type UserServ interface {
	HaveUser(queryParam db.GetUserParams)
}

type UserServImpl struct {
	q *db.Queries
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
