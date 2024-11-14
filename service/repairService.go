package service

import (
	"context"
	"database/sql"
	db "main/sql"
)

type RepairServ interface {
	NewRepair(param db.CreateNewRepairParams) (int64, error)
	NewRepairInfo(param db.CreateNewRepairInfoParams) (int64, error)
	GetRepair(param db.GetRepairParams) ([]db.GetRepairRow, error)
	DeleteRepair(param int64) error
	ApproveRepair(param int64) error
	// GetRepairById(param int64) ([]db.Repairt, error)
	GetRepairInfoById(param int64) ([]db.Repairinfot, error)
	GetRepairDate(param int64) ([]string, error)
	UpdateItem(param db.UpdateItemParams) error
}

type RepairServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (s *RepairServImpl) GetRepairDate(param int64) ([]string, error) {
	res, err := s.q.GetRepairDate(context.Background(), param)
	return res, err
}

func (r *RepairServImpl) NewRepair(param db.CreateNewRepairParams) (int64, error) {
	res, err := r.q.CreateNewRepair(context.Background(), param)
	return res, err
}

func (r *RepairServImpl) NewRepairInfo(param db.CreateNewRepairInfoParams) (int64, error) {
	res, err := r.q.CreateNewRepairInfo(context.Background(), param)
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

func (r *RepairServImpl) GetRepairInfoById(param int64) ([]db.Repairinfot, error) {
	res, err := r.q.GetRepairInfoById(context.Background(), param)
	if err == sql.ErrNoRows {

		var r []db.Repairinfot
		return r, nil
	}
	return res, err
}

func (r *RepairServImpl) UpdateItem(param db.UpdateItemParams) error {
	err := r.q.UpdateItem(context.Background(), param)
	return err
}

func RepairServInit(q *db.Queries, conn *sql.DB) *RepairServImpl {
	return &RepairServImpl{
		q:    q,
		conn: conn,
	}
}
