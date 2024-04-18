package service

import (
	"context"
	"database/sql"
	db "main/sql"
)

type JobsServ interface {
	GetAllJobs() ([]db.Jobst, error)
	GetAllJobsByCmp(belongCmp int64) ([]db.Jobst, error)
	DeleteJob(id int64) error
}

type JobsServImpl struct {
	q    *db.Queries
	conn *sql.DB
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

func JobsServInit(q *db.Queries, conn *sql.DB) *JobsServImpl {
	return &JobsServImpl{
		q:    q,
		conn: conn,
	}
}
