package service

import (
	"context"
	"database/sql"
	db "main/sql"
)

type RevenueServ interface {
	GetRevenue(param db.GetDriverRevenueParams) ([]db.GetDriverRevenueRow, error)
}

type RevenueServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (s *RevenueServImpl) GetRevenue(param db.GetDriverRevenueParams) ([]db.GetDriverRevenueRow, error) {
	res, err := s.q.GetDriverRevenue(context.Background(), param)
	if len(res) == 0 {
		res = []db.GetDriverRevenueRow{}
		tmp := db.GetDriverRevenueRow{
			Earn:  0,
			Count: 0,
		}
		res = append(res, tmp)

		return res, err
	}
	return res, err
}

func RevenueServInit(q *db.Queries, conn *sql.DB) *RevenueServImpl {
	return &RevenueServImpl{
		q:    q,
		conn: conn,
	}
}
