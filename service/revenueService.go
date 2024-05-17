package service

import (
	"context"
	"database/sql"
	db "main/sql"
)

type RevenueServ interface {
	GetRevenue(param db.GetDriverRevenueParams) ([]db.GetDriverRevenueRow, error)
	GetRevenueByCmp(param db.GetDriverRevenueByCmpParams) ([]db.GetDriverRevenueByCmpRow, error)
}

type RevenueServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (s *RevenueServImpl) GetRevenueByCmp(param db.GetDriverRevenueByCmpParams) ([]db.GetDriverRevenueByCmpRow, error) {
	res, err := s.q.GetDriverRevenueByCmp(context.Background(), param)
	if len(res) == 0 {
		res = []db.GetDriverRevenueByCmpRow{}
		tmp := db.GetDriverRevenueByCmpRow{
			Earn:  0,
			Count: 0,
		}
		res = append(res, tmp)

		return res, err
	}
	return res, err
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
