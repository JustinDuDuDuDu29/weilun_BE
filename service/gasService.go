package service

import (
	"context"
	"database/sql"
	db "main/sql"
)

type GasServ interface {
	NewGas(param db.CreateNewGasParams) (int64, error)
	NewGasInfo(param db.CreateNewGasInfoParams) (int64, error)
	GetGas(param db.GetGasParams) ([]db.GetGasRow, error)
	DeleteGas(param int64) error
	ApproveGas(param int64) error
	// GetRepairById(param int64) ([]db.Repairt, error)
	GetGasInfoById(param int64) ([]db.Gasinfot, error)
	GetGasDate(param int64) ([]string, error)
}

type GasServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (s *GasServImpl) GetGasDate(param int64) ([]string, error) {
	res, err := s.q.GetGasDate(context.Background(), param)
	return res, err
}

func (r *GasServImpl) NewGas(param db.CreateNewGasParams) (int64, error) {
	res, err := r.q.CreateNewGas(context.Background(), param)
	return res, err
}

func (r *GasServImpl) NewGasInfo(param db.CreateNewGasInfoParams) (int64, error) {
	res, err := r.q.CreateNewGasInfo(context.Background(), param)
	return res, err
}

func (r *GasServImpl) GetGas(param db.GetGasParams) ([]db.GetGasRow, error) {
	res, err := r.q.GetGas(context.Background(), param)
	return res, err
}

func (r *GasServImpl) DeleteGas(param int64) error {
	err := r.q.DeleteGasT(context.Background(), param)
	return err
}

func (r *GasServImpl) ApproveGas(param int64) error {
	err := r.q.ApproveGas(context.Background(), param)
	return err
}

func (r *GasServImpl) GetGasInfoById(param int64) ([]db.Gasinfot, error) {
	res, err := r.q.GetGasInfoById(context.Background(), param)
	if err == sql.ErrNoRows {

		var r []db.Gasinfot
		return r, nil
	}
	return res, err
}

func GasServInit(q *db.Queries, conn *sql.DB) *GasServImpl {
	return &GasServImpl{
		q:    q,
		conn: conn,
	}
}
