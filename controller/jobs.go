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
	"strings"
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

	role := c.MustGet("Role").(int16)
	belongCmp := c.MustGet("belongCmp").(int64)
	// UserID := c.MustGet("UserID").(int)

	if role <= 100 {
		var reqBody apptypes.GetJobsBodyT
		if err := c.Bind(&reqBody); err != nil {
			c.Abort()
			return
		}

		var ID sql.NullInt64
		if reqBody.ID != 0 {
			ID.Scan(reqBody.ID)
		}

		var FromLoc sql.NullString
		if reqBody.FromLoc != "" {
			FromLoc.Scan(reqBody.FromLoc)
		}

		var Mid sql.NullString
		if reqBody.Mid != "" {
			Mid.Scan(reqBody.Mid)
		}

		var ToLoc sql.NullString
		if reqBody.ToLoc != "" {
			ToLoc.Scan(reqBody.ToLoc)
		}

		var Belongcmp sql.NullInt64
		if reqBody.Belongcmp != 0 {
			Belongcmp.Scan(reqBody.Belongcmp)
		}

		var Remaining sql.NullInt16
		if reqBody.Remaining != 0 {
			Remaining.Scan(reqBody.Remaining)
		}

		var CloseDateStart sql.NullTime
		if reqBody.CloseDateStart != "" {
			dt, err := time.Parse(time.DateOnly, reqBody.CloseDateStart)
			if err != nil {
				c.Status(http.StatusBadRequest)
				c.Abort()
				return
			}
			CloseDateStart.Scan(dt)
		}

		var CloseDateEnd sql.NullTime
		if reqBody.CloseDateEnd != "" {
			dt, err := time.Parse(time.DateOnly, reqBody.CloseDateEnd)
			if err != nil {
				c.Status(http.StatusBadRequest)
				c.Abort()
				return
			}
			CloseDateEnd.Scan(dt)
		}

		var CreateDateStart sql.NullTime
		if reqBody.CreateDateStart != "" {
			dt, err := time.Parse(time.DateOnly, reqBody.CreateDateStart)
			if err != nil {
				c.Status(http.StatusBadRequest)
				c.Abort()
				return
			}
			CreateDateStart.Scan(dt)
		}

		var CreateDateEnd sql.NullTime
		if reqBody.CreateDateEnd != "" {
			dt, err := time.Parse(time.DateOnly, reqBody.CreateDateEnd)
			if err != nil {
				c.Status(http.StatusBadRequest)
				c.Abort()
				return
			}
			CreateDateEnd.Scan(dt)
		}

		var DeletedDateStart sql.NullTime
		if reqBody.DeletedDateStart != "" {
			dt, err := time.Parse(time.DateOnly, reqBody.DeletedDateStart)
			if err != nil {
				c.Status(http.StatusBadRequest)
				c.Abort()
				return
			}
			DeletedDateStart.Scan(dt)
		}

		var DeletedDateEnd sql.NullTime
		if reqBody.DeletedDateEnd != "" {
			dt, err := time.Parse(time.DateOnly, reqBody.DeletedDateEnd)
			if err != nil {
				c.Status(http.StatusBadRequest)
				c.Abort()
				return
			}
			DeletedDateEnd.Scan(dt)
		}

		var LastModifiedDateStart sql.NullTime
		if reqBody.LastModifiedDateStart != "" {
			dt, err := time.Parse(time.DateOnly, reqBody.LastModifiedDateStart)
			if err != nil {
				c.Status(http.StatusBadRequest)
				c.Abort()
				return
			}
			LastModifiedDateStart.Scan(dt)
		}

		var LastModifiedDateEnd sql.NullTime
		if reqBody.LastModifiedDateEnd != "" {
			dt, err := time.Parse(time.DateOnly, reqBody.LastModifiedDateEnd)
			if err != nil {
				c.Status(http.StatusBadRequest)
				c.Abort()
				return
			}
			LastModifiedDateEnd.Scan(dt)
		}

		param := db.GetAllJobsAdminParams{
			ID:                    ID,
			FromLoc:               FromLoc,
			Mid:                   Mid,
			ToLoc:                 ToLoc,
			Belongcmp:             Belongcmp,
			Remaining:             Remaining,
			CloseDateStart:        CloseDateStart,
			CloseDateEnd:          CloseDateEnd,
			CreateDateStart:       CreateDateStart,
			CreateDateEnd:         CreateDateEnd,
			DeletedDateStart:      DeletedDateStart,
			DeletedDateEnd:        DeletedDateEnd,
			LastModifiedDateStart: LastModifiedDateStart,
			LastModifiedDateEnd:   LastModifiedDateEnd,
		}
		res, err := u.svc.JobsServ.GetAllJobs(param)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, res)

	} else {

		var reqBody apptypes.GetJobsClientBodyT
		if err := c.Bind(&reqBody); err != nil {
			c.Abort()
			return
		}

		var ID sql.NullInt64
		if reqBody.ID != 0 {
			ID.Scan(reqBody.ID)
		}

		var FromLoc sql.NullString
		if reqBody.FromLoc != "" {
			FromLoc.Scan(reqBody.FromLoc)
		}

		var Mid sql.NullString
		if reqBody.Mid != "" {
			Mid.Scan(reqBody.Mid)
		}

		var ToLoc sql.NullString
		if reqBody.ToLoc != "" {
			ToLoc.Scan(reqBody.ToLoc)
		}

		var Belongcmp sql.NullInt64
		Belongcmp.Scan(belongCmp)

		param := db.GetAllJobsClientParams{
			ID:        ID,
			FromLoc:   FromLoc,
			Mid:       Mid,
			ToLoc:     ToLoc,
			Belongcmp: Belongcmp,
		}
		res, err := u.svc.JobsServ.GetAllJobsClient(param)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, res)
	}

}

func (u *JobsCtrlImpl) GetAllJobN(c *gin.Context) {

	role := c.MustGet("Role").(int16)
	cuid := c.MustGet("UserID").(int)
	belongCmp := c.MustGet("belongCmp").(int64)
	var ID sql.NullInt64
	if role >= 200 {
		ID.Valid = false
	} else {
		if c.Query("Id") == "" {
			ID.Valid = false
		} else {
			ID.Scan(c.Query("Id"))
		}
	}

	var BelongCMP sql.NullInt64
	if role >= 200 {
		BelongCMP.Scan(belongCmp)
	} else {
		if c.Query("belongCMP") == "" {
			BelongCMP.Valid = false
		} else {
			BelongCMP.Scan(c.Query("belongCMP"))
		}
	}

	var Alert sql.NullString
	if c.Query("alert") == "" {
		Alert.Valid = false
	} else {
		Alert.Scan(c.Query("alert"))
	}

	var CreateDateStart sql.NullTime
	if c.Query("CreateDateStart") == "" {
		CreateDateStart.Valid = false
	} else {
		CreateDateStart.Scan(c.Query("CreateDateStart"))
	}

	var CreateDateEnd sql.NullTime
	if c.Query("CreateDateEnd") == "" {
		CreateDateEnd.Valid = false
	} else {
		CreateDateEnd.Scan(c.Query("CreateDateEnd"))
	}

	var DeletedDateStart sql.NullTime
	if c.Query("DeletedDateStart") == "" {
		DeletedDateStart.Valid = false
	} else {
		DeletedDateStart.Scan(c.Query("DeletedDateStart"))
	}

	var DeletedDateEnd sql.NullTime
	if c.Query("DeletedDateEnd") == "" {
		DeletedDateEnd.Valid = false
	} else {
		DeletedDateEnd.Scan(c.Query("DeletedDateEnd"))
	}

	var LastModifiedDateStart sql.NullTime
	if c.Query("LastModifiedDateStart") == "" {
		LastModifiedDateStart.Valid = false
	} else {
		LastModifiedDateStart.Scan(c.Query("LastModifiedDateStart"))
	}

	var LastModifiedDateEnd sql.NullTime
	if c.Query("LastModifiedDateEnd") == "" {
		LastModifiedDateEnd.Valid = false
	} else {
		LastModifiedDateEnd.Scan(c.Query("LastModifiedDateEnd"))
	}

	param := db.GetAlertParams{
		ID:                    ID,
		BelongCMP:             BelongCMP,
		Alert:                 Alert,
		CreateDateStart:       CreateDateStart,
		CreateDateEnd:         CreateDateEnd,
		DeletedDateStart:      DeletedDateStart,
		DeletedDateEnd:        DeletedDateEnd,
		LastModifiedDateStart: LastModifiedDateStart,
		LastModifiedDateEnd:   LastModifiedDateEnd,
	}
	res, err := u.svc.AlertServ.GetAlert(param)

	if err != nil && err != sql.ErrNoRows {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err == sql.ErrNoRows || len(res) == 0 {
		c.JSON(http.StatusOK, gin.H{})
		c.Abort()
		return
	}

	var Lastalert sql.NullInt64
	Lastalert.Scan(res[0].ID)

	updateParam := db.UpdateLastAlertParams{
		ID:        int64(cuid),
		Lastalert: Lastalert,
	}

	if role >= 300 {
		u.svc.AlertServ.UpdateLastAlert(updateParam)
	}

	c.JSON(http.StatusOK, gin.H{"res": res})
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
	res, err, num := u.svc.JobsServ.ClaimJob(param)

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
	if num == 1 {
		SandMsg(1, 400, "check job open")
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
		return
	}

	UserID := c.MustGet("UserID").(int)

	var reqBody apptypes.FinishClaimJobBodyT
	if err := c.Bind(&reqBody); err != nil {
		fmt.Print("NOT OK")
		c.Abort()
		return
	}

	cType := reqBody.File.Header["Content-Type"][0]

	path, uuid, err := utils.GenPicRoute(cType)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Print(reqBody.File)

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

		fmt.Print(err)
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

	cJobInfo, err := u.svc.JobsServ.GetClaimedJobByID(int64(id))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	SandMsg(int(cJobInfo.Driverid), 300, "Job "+strconv.Itoa(id)+" is approved")

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

	num, err := u.svc.JobsServ.DeleteClaimedJob(param)
	if num == 1 {
		SandMsg(1, 400, "check job open")
	}

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

	if !(reqBody.Mid == "") {
		Mid.Scan(reqBody.Mid)
	} else {
		Mid.Valid = false
	}

	var Memo sql.NullString
	if !(reqBody.Memo == "") {
		Memo.Scan(reqBody.Memo)
	} else {
		Memo.Valid = false
	}

	Jobdate, err := time.Parse(time.DateOnly, strings.Split(reqBody.Jobdate, "T")[0])
	if err != nil {
		fmt.Print(err)
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	var CloseDate sql.NullTime
	if reqBody.CloseDate != "" {

		ct, err := time.Parse(time.DateOnly, strings.Split(reqBody.CloseDate, "T")[0])
		if err != nil {

			fmt.Print(2)
			c.Status(http.StatusBadRequest)
			c.Abort()
			return
		}
		fmt.Print("ct: ", ct)
		CloseDate.Scan(ct)

	} else {
		CloseDate.Valid = false
	}

	var UserID sql.NullInt64
	UserID.Scan(cuid)

	param := db.UpdateJobParams{
		ID:        int64(reqBody.ID),
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
