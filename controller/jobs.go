package controller

import (
	"database/sql"
	"fmt"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"main/utils"
	"net/http"
	"os"
	"strconv"
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
	GetAllClaimedJobs(c *gin.Context)
	GetCurrentClaimedJob(c *gin.Context)
	ApproveClaimedJob(c *gin.Context)
}

type JobsCtrlImpl struct {
	svc *service.AppService
}

func (u *JobsCtrlImpl) GetAllClaimedJobs(c *gin.Context) {
	res, err := u.svc.JobsServ.GetAllClaimedJobs()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *JobsCtrlImpl) GetCurrentClaimedJob(c *gin.Context) {
	cuid := c.MustGet("UserID").(int)

	res, err := u.svc.JobsServ.GetCurrentClaimedJob(int64(cuid))

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{})
			return
		}

		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *JobsCtrlImpl) GetAllJob(c *gin.Context) {

	jobList, err := u.svc.JobsServ.GetAllJobs()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, jobList)
}

func (u *JobsCtrlImpl) ClaimJob(c *gin.Context) {

	if c.Param("id") == "" {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	UserID := c.MustGet("UserID").(int)

	param := db.ClaimJobParams{
		Jobid:    int64(id),
		Driverid: int64(UserID),
	}
	res, err := u.svc.JobsServ.ClaimJob(param)

	if err != nil {
		fmt.Print(err)
		if err.Error() == "already have ongoing job" {
			fmt.Print("already have ongoing job")
			res, err := u.svc.JobsServ.GetClaimedJobByID(res)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				c.Abort()
				return
			}
			c.JSON(http.StatusConflict, res)
		}
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": res})

}

func (u *JobsCtrlImpl) FinishClaimJob(c *gin.Context) {
	sid := c.Param("id")
	if sid == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	UserID := c.MustGet("UserID").(int)

	var reqBody apptypes.FinishClaimJobBodyT
	if err := c.Bind(&reqBody); err != nil {
		c.Abort()
		return
	}

	cType := reqBody.File.Header["Content-Type"][0]

	path, uuid, err := utils.GenPicRoute(cType)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = c.SaveUploadedFile(reqBody.File, path)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var Fp sql.NullString
	Fp.Scan(uuid)

	param := db.FinishClaimedJobParams{
		ID:        int64(id),
		Driverid:  int64(UserID),
		Finishpic: Fp,
	}

	err = u.svc.JobsServ.FinishClaimedJob(param)
	if err != nil {
		os.Remove(path)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)

}

func (u *JobsCtrlImpl) ApproveClaimedJob(c *gin.Context) {
	UserID := c.MustGet("UserID").(int)
	if c.Param("id") == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var ApprovedBy sql.NullInt64
	ApprovedBy.Scan(UserID)

	param := db.ApproveFinishedJobParams{
		ID:         int64(id),
		ApprovedBy: ApprovedBy,
	}
	err = u.svc.JobsServ.ApproveFinishedJob(param)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)

}

func (u *JobsCtrlImpl) CancelClaimJob(c *gin.Context) {
	// call delete
	UserID := c.MustGet("UserID").(int)
	if c.Param("id") == "" {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	res, err := u.svc.UserServ.GetUserById(int64(UserID))
	if err != nil {
		fmt.Print(err)
		c.Abort()
		return
	}
	cJobRes, err := u.svc.JobsServ.GetClaimedJobByID(int64(id))

	if err != nil {
		fmt.Print(err)
		c.Abort()
		return
	}
	if !(cJobRes.CreateDate.Add(time.Minute*10).After(time.Now()) && cJobRes.Driverid == int64(UserID)) && !(res.Role <= int16(100)) {
		// reject has pass 5 min

		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	var uid sql.NullInt64
	uid.Scan(UserID)

	param := db.DeleteClaimedJobParams{
		ID:        int64(id),
		DeletedBy: uid,
	}

	err = u.svc.JobsServ.DeleteClaimedJob(param)

	if err != nil {
		fmt.Print(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.Status(http.StatusOK)
	c.Abort()
}

func (u *JobsCtrlImpl) CreateJob(c *gin.Context) {

	var reqBody apptypes.CreateJobBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		c.Abort()
		return
	}

	// cuid := c.MustGet("UserID").(int)

	var Mid sql.NullString
	if reqBody.Mid == "" {
		Mid.Valid = false
	} else {
		Mid.Scan(reqBody.Mid)
	}

	var Memo sql.NullString
	if reqBody.Memo == "" {
		Memo.Valid = false
	} else {
		Memo.Scan(reqBody.Memo)
	}

	Jobdate, err := time.Parse(time.DateOnly, reqBody.Jobdate)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	var CloseDate sql.NullTime
	if reqBody.CloseDate == "" {
		CloseDate.Valid = false
	} else {
		CloseDate.Scan(reqBody.CloseDate)
	}

	// var UserID sql.NullInt64
	// UserID.Scan(cuid)

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
		CloseDate: CloseDate,
	}
	res, err := u.svc.JobsServ.CreateJob(param)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": res})

}

func (u *JobsCtrlImpl) DeleteJob(c *gin.Context) {
	if c.Param("id") == "" {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	err = u.svc.JobsServ.DeleteJob(int64(id))

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
	c.Abort()

}

func (u *JobsCtrlImpl) UpdateJob(c *gin.Context) {

	var reqBody apptypes.UpdateJobBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		fmt.Print(err)
		c.Abort()
		return
	}

	cuid := c.MustGet("UserID").(int)

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

	var CloseDate sql.NullTime
	CloseDate.Scan(reqBody.CloseDate)

	var UserID sql.NullInt64
	UserID.Scan(cuid)

	param := db.UpdateJobParams{
		FromLoc:   reqBody.FromLoc,
		Mid:       Mid,
		ToLoc:     reqBody.ToLoc,
		Price:     int16(reqBody.Price),
		Belongcmp: int64(reqBody.Belongcmp),
		Source:    reqBody.Source,
		Jobdate:   Jobdate,
		Memo:      Memo,
		CloseDate: CloseDate,
		Remaining: int16(reqBody.Remaining),
	}
	res, err := u.svc.JobsServ.UpdateJob(param)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": res})
}

func JobsCtrlInit(svc *service.AppService) *JobsCtrlImpl {
	return &JobsCtrlImpl{
		svc: svc,
	}
}
