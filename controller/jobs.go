package controller

import (
	"database/sql"
	"fmt"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type JobsCtrl interface {
	CreateJob(c *gin.Context)
	GetAllJob(c *gin.Context)
	DeleteJob(c *gin.Context)
	UpdateJob(c *gin.Context)
	ClaimJob(c *gin.Context)
	FinishClaimJob(c *gin.Context)
	CancelClaimJob(c *gin.Context)
}

type JobsCtrlImpl struct {
	svc *service.AppService
}

func (u *JobsCtrlImpl) GetAllJob(c *gin.Context) {
	cuid := c.MustGet("UserID")

	var UserID sql.NullInt64
	UserID.Scan(cuid)

	jobList, err := u.svc.JobsServ.GetAllJobs()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, jobList)
}

func (u *JobsCtrlImpl) ClaimJob(c *gin.Context) {

	var reqBody apptypes.ClaimJobBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		fmt.Print(err)
		c.Abort()
		return
	}

	UserID := c.MustGet("UserID").(int64)

	param := db.ClaimJobParams{

		Jobid:    int64(reqBody.JobID),
		Driverid: UserID,
	}
	res, err := u.svc.JobsServ.ClaimJob(param)

	if err != nil {
		fmt.Print(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": res})
	return

}

func (u *JobsCtrlImpl) FinishClaimJob(c *gin.Context) {
	// UserID := c.MustGet("UserID")

}

func (u *JobsCtrlImpl) CancelClaimJob(c *gin.Context) {
	// call delete
	UserID := c.MustGet("UserID").(int64)
	var reqBody apptypes.CreateJobBodyT
	err := c.BindJSON(&reqBody)

	res, err := u.svc.UserServ.GetUserById(UserID)
    cJobRes, err := 

	if err != nil {
		fmt.Print(err)
		c.Abort()
		return
	}
    
    if !(res.Role <= int16(100)) || UserID


}

func (u *JobsCtrlImpl) CreateJob(c *gin.Context) {

	var reqBody apptypes.CreateJobBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		fmt.Print(err)
		c.Abort()
		return
	}

	cuid := c.MustGet("UserID")

	var Mid sql.NullString
	Mid.Scan(reqBody.Mid)

	var Memo sql.NullString
	Memo.Scan(reqBody.Memo)

	Jobdate, err := time.Parse(time.DateOnly, reqBody.Jobdate)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	var EndDate sql.NullTime
	EndDate.Scan(reqBody.EndDate)

	var UserID sql.NullInt64
	UserID.Scan(cuid)

	param := db.CreateJobParams{
		FromLoc:   reqBody.FromLoc,
		Mid:       Mid,
		ToLoc:     reqBody.ToLoc,
		Price:     int16(reqBody.Price),
		Estimated: int16(reqBody.Estimated),
		Belongcmp: int64(reqBody.Belongcmp),
		Source:    reqBody.Source,
		Jobdate:   Jobdate,
		Memo:      Memo,
		EndDate:   EndDate,
	}
	res, err := u.svc.JobsServ.CreateJob(param)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": res})

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
