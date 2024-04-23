package service

import (
	"database/sql"
	db "main/sql"
)

type RepairServ interface {
}

type RepairServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func RepairServInit(q *db.Queries, conn *sql.DB) *RepairServImpl {
	return &RepairServImpl{
		q:    q,
		conn: conn,
	}
}
