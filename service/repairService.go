package service

import (
	"context"
	"database/sql"
	db "main/sql"
)

type RepairServ interface {
	NewRepair(param db.CreateNewRepairParams) (int64, error)
	GetRepair(param db.GetRepairParams) ([]db.GetRepairRow, error)
	DeleteRepair(param int64) error
	ApproveRepair(param int64) error
	GetRepairById(param int64) (db.GetRepairByIdRow, error)
}

type RepairServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (r *RepairServImpl) NewRepair(param db.CreateNewRepairParams) (int64, error) {
	res, err := r.q.CreateNewRepair(context.Background(), param)
	return res, err
}

func (r *RepairServImpl) GetRepair(param db.GetRepairParams) ([]db.GetRepairRow, error) {
	res, err := r.q.GetRepair(context.Background(), param)
	return res, err
}

func (r *RepairServImpl) DeleteRepair(param int64) error {
	err := r.q.DeleteRepair(context.Background(), param)
	return err
}

func (r *RepairServImpl) ApproveRepair(param int64) error {
	err := r.q.ApproveRepair(context.Background(), param)
	return err
}

func (r *RepairServImpl) GetRepairById(param int64) (db.GetRepairByIdRow, error) {
	res, err := r.q.GetRepairById(context.Background(), param)
	return res, err
}

func RepairServInit(q *db.Queries, conn *sql.DB) *RepairServImpl {
	return &RepairServImpl{
		q:    q,
		conn: conn,
	}
}
