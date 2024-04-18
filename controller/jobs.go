package controller

import (
	"main/service"

	"github.com/gin-gonic/gin"
)

type JobsCtrl interface {
	CreateJob(c *gin.Context)
	GetJob(c *gin.Context)
	DeleteJob(c *gin.Context)
	UpdateJob(c *gin.Context)
	ClaimJob(c *gin.Context)
	FinishClaimJob(c *gin.Context)
	CancelClaimJob(c *gin.Context)
}

type JobsCtrlImpl struct {
	svc *service.AppService
}

func (u *JobsCtrlImpl) GetJob(c *gin.Context) {
	// UserID := c.MustGet("UserID")

}

func (u *JobsCtrlImpl) ClaimJob(c *gin.Context) {
	// UserID := c.MustGet("UserID")

}

func (u *JobsCtrlImpl) FinishClaimJob(c *gin.Context) {
	// UserID := c.MustGet("UserID")

}

func (u *JobsCtrlImpl) CancelClaimJob(c *gin.Context) {
	// UserID := c.MustGet("UserID")

}

func (u *JobsCtrlImpl) CreateJob(c *gin.Context) {
	// UserID := c.MustGet("UserID")

}

func (u *JobsCtrlImpl) DeleteJob(c *gin.Context) {
	// UserID := c.MustGet("UserID")

}

func (u *JobsCtrlImpl) UpdateJob(c *gin.Context) {
	// UserID := c.MustGet("UserID")

}

func JobsCtrlInit(svc *service.AppService) *JobsCtrlImpl {
	return &JobsCtrlImpl{
		svc: svc,
	}
}
