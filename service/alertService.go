package service

import (
	"context"
	"database/sql"
	"fmt"
	db "main/sql"
)

type AlertServ interface {
	CreateAlert(param db.CreateAlertParams) (int64, error)
	GetAlert(param db.GetAlertParams) ([]db.Alertt, error)
	DeleteAlert(id int64) error
	UpdateAlert(param db.UpdateAlertParams) error
	UpdateLastAlert(param db.UpdateLastAlertParams) error
	HaveNewAlert(id int64) (bool, error)
}

type AlertServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (s *AlertServImpl) HaveNewAlert(id int64) (bool, error) {
	res, err := s.q.GetDriver(context.Background(), id)
	if err != nil {
		fmt.Print("1")
		return false, err
	}
	cmpLastalert, err := s.q.GetAlertByCmp(context.Background(), res.Belongcmp)
	if err != nil && err != sql.ErrNoRows {
		fmt.Print("2")
		return false, err
	}

	if res.Lastalert.Int64 < cmpLastalert[0].ID {
		fmt.Print("3")

		return true, nil
	}
	fmt.Println("4")
	fmt.Println(res.Lastalert.Int64)
	fmt.Println(cmpLastalert[0].ID)
	return false, nil
}

func (s *AlertServImpl) CreateAlert(param db.CreateAlertParams) (int64, error) {
	res, err := s.q.CreateAlert(context.Background(), param)
	return res, err
}

func (s *AlertServImpl) GetAlert(param db.GetAlertParams) ([]db.Alertt, error) {
	res, err := s.q.GetAlert(context.Background(), param)
	return res, err
}

func (s *AlertServImpl) DeleteAlert(id int64) error {
	err := s.q.DeleteAlert(context.Background(), id)
	return err
}

func (s *AlertServImpl) UpdateAlert(param db.UpdateAlertParams) error {
	err := s.q.UpdateAlert(context.Background(), param)
	return err
}

func (s *AlertServImpl) UpdateLastAlert(param db.UpdateLastAlertParams) error {
	tx, err := s.conn.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	qtx := s.q.WithTx(tx)
	res, err := qtx.GetLastAlert(context.Background(), param.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if !res.Valid || res.Int64 < param.Lastalert.Int64 {
		err = qtx.UpdateLastAlert(context.Background(), param)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func AlertServInit(q *db.Queries, conn *sql.DB) *AlertServImpl {
	return &AlertServImpl{
		q:    q,
		conn: conn,
	}
}
