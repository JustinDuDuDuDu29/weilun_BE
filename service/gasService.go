package service

import (
	"context"
	"database/sql"
	"encoding/json"
	db "main/sql"
)

type GasServ interface {
	NewGas(param db.CreateNewGasParams) (int64, error)
	NewGasInfo(param db.CreateNewGasInfoParams) (int64, error)
	GetGas(param db.GetGasParams) ([]db.GetGasRow, error)
	GetGasCmpUser(param sql.NullInt64) ([]json.RawMessage, error)

	DeleteGas(param int64) error
	ApproveGas(param int64) error
	// GetRepairById(param int64) ([]db.Repairt, error)
	GetGasInfoById(param int64) ([]db.Gasinfot, error)
	GetGasDate(param int64) ([]string, error)
	UpdateGas(param db.UpdateGasParams) error
}

type GasServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (r *GasServImpl) UpdateGas(param db.UpdateGasParams) error {
	err := r.q.UpdateGas(context.Background(), param)
	return err
}

func (r *GasServImpl) GetGasCmpUser(param sql.NullInt64) ([]json.RawMessage, error) {
	res, err := r.q.GetGasCmpUser(context.Background(), param)
	return res, err
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
	// Fetch data from the database
	res, err := r.q.GetGas(context.Background(), param)
	if err != nil {
		return nil, err
	}

	// Iterate over the result and convert Repairinfo from interface{} to json.RawMessage
	for i := range res {
		if raw, ok := res[i].Repairinfo.([]byte); ok {
			res[i].Repairinfo = json.RawMessage(raw)
		}
	}

	return res, nil
}

// func (r *GasServImpl) GetGas(param db.GetGasParams) ([]db.GetGasRow, error) {
// 	res, err := r.q.GetGas(context.Background(), param)
// 	return res, err
// }

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
