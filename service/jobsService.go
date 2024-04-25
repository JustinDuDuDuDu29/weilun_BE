package service

import (
	"context"
	"database/sql"
	"errors"
	db "main/sql"
)

type JobsServ interface {
	GetClaimedJobByID(id int64) (db.Claimjobt, error)
	CreateJob(param db.CreateJobParams) (int64, error)
	IncreaseRemaining(id int64) error
	DecreaseRemaining(id int64) error
	FinishClaimedJob(param db.FinishClaimedJobParams) error
	ClaimJob(param db.ClaimJobParams) (int64, error)
	GetAllJobs() ([]db.Jobst, error)
	GetAllJobsByCmp(belongCmp int64) ([]db.Jobst, error)
	DeleteJob(id int64) error
	UpdateJob(param db.UpdateJobParams) (int64, error)
	GetCurrentClaimedJob(id int64) (db.GetCurrentClaimedJobRow, error)
	DeleteClaimedJob(param db.DeleteClaimedJobParams) error
	ApproveFinishedJob(param db.ApproveFinishedJobParams) error
	SetJobNoMore(id int64) error
	GetAllClaimedJobs() ([]db.Claimjobt, error)
}

type JobsServImpl struct {
	q    *db.Queries
	conn *sql.DB
}

func (s *JobsServImpl) CreateJob(param db.CreateJobParams) (int64, error) {
	res, err := s.q.CreateJob(context.Background(), param)
	return res, err
}

func (s *JobsServImpl) SetJobNoMore(id int64) error {
	err := s.q.SetJobNoMore(context.Background(), id)
	return err
}

func (s *JobsServImpl) ApproveFinishedJob(param db.ApproveFinishedJobParams) error {
	err := s.q.ApproveFinishedJob(context.Background(), param)
	return err
}

func (s *JobsServImpl) DeleteClaimedJob(param db.DeleteClaimedJobParams) error {

	tx, err := s.conn.BeginTx(context.Background(), nil)

	if err != nil {
		return err
	}

	qtx := s.q.WithTx(tx)
	err = qtx.DeleteClaimedJob(context.Background(), param)

	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := qtx.GetClaimedJobByID(context.Background(), param.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	jres, err := qtx.GetJobById(context.Background(), res.Jobid)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !(jres.CloseDate.Valid) {
		err = qtx.IncreaseRemaining(context.Background(), res.Jobid)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return err
}

func (s *JobsServImpl) GetCurrentClaimedJob(id int64) (db.GetCurrentClaimedJobRow, error) {
	res, err := s.q.GetCurrentClaimedJob(context.Background(), id)

	return res, err
}

func (s *JobsServImpl) IncreaseRemaining(id int64) error {
	err := s.q.IncreaseRemaining(context.Background(), id)

	return err
}

func (s *JobsServImpl) DecreaseRemaining(id int64) error {
	err := s.q.DecreaseRemaining(context.Background(), id)
	return err
}

func (s *JobsServImpl) GetClaimedJobByID(id int64) (db.Claimjobt, error) {
	res, err := s.q.GetClaimedJobByID(context.Background(), id)
	return res, err
}

func (s *JobsServImpl) GetAllClaimedJobs() ([]db.Claimjobt, error) {
	res, err := s.q.GetAllClaimedJobs(context.Background())
	return res, err
}

func (s *JobsServImpl) FinishClaimedJob(param db.FinishClaimedJobParams) error {
	err := s.q.FinishClaimedJob(context.Background(), param)
	return err
}

func (s *JobsServImpl) ClaimJob(param db.ClaimJobParams) (int64, error) {
	tx, err := s.conn.BeginTx(context.Background(), nil)

	if err != nil {
		return -99, err
	}

	qtx := s.q.WithTx(tx)

	cres, err := qtx.GetCurrentClaimedJob(context.Background(), param.Driverid)

	if err == nil {
		return cres.ID, errors.New("already have ongoing job")
	}

	// if err != nil && err != sql.ErrNoRows {
	if err != sql.ErrNoRows {
		tx.Rollback()
		return -99, err
	}

	res, err := qtx.ClaimJob(context.Background(), param)
	if err != nil {
		tx.Rollback()
		return -99, err
	}

	err = qtx.DecreaseRemaining(context.Background(), param.Jobid)
	if err != nil {
		tx.Rollback()
		return -99, err
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
	}

	return res, err
}

func (s *JobsServImpl) GetAllJobs() ([]db.Jobst, error) {
	res, err := s.q.GetAllJobs(context.Background())
	return res, err
}

func (s *JobsServImpl) GetAllJobsByCmp(belongCmp int64) ([]db.Jobst, error) {
	res, err := s.q.GetAllJobsByCmp(context.Background(), belongCmp)
	return res, err
}

func (s *JobsServImpl) DeleteJob(id int64) error {
	err := s.q.DeleteJob(context.Background(), id)
	return err
}

func (s *JobsServImpl) UpdateJob(param db.UpdateJobParams) (int64, error) {
	res, err := s.q.UpdateJob(context.Background(), param)
	return res, err
}

func JobsServInit(q *db.Queries, conn *sql.DB) *JobsServImpl {
	return &JobsServImpl{
		q:    q,
		conn: conn,
	}
}
