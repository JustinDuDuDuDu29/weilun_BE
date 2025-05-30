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
	GetClaimedJobByID(c *gin.Context)
	ClaimJob(c *gin.Context)
	FinishClaimJob(c *gin.Context)
	CancelClaimJob(c *gin.Context)
	GetAllClaimedJobs(c *gin.Context)
	GetUserWithPendingJob(c *gin.Context)
	GetCurrentClaimedJob(c *gin.Context)
	ApproveClaimedJob(c *gin.Context)
	ApproveClaimedJobs(c *gin.Context)
	GetClaimedJobByDriverID(c *gin.Context)
	GetCJDate(c *gin.Context)
}

type JobsCtrlImpl struct {
	svc *service.AppService
}

func (u *JobsCtrlImpl) GetClaimedJobByDriverID(c *gin.Context) {
	uid := c.MustGet("UserID").(int)
	role := c.MustGet("Role")
	bcmp := c.MustGet("belongCmp")

	qid := c.Query("id")
	qcmp := c.Query("cmp")

	if qcmp != "" {

		cmp, err := strconv.Atoi(qcmp)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if role.(int16) >= 300 || (role == 200 && bcmp != cmp) {
			// ???
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//start
		year, err := strconv.Atoi(c.Query("year"))

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		month, err := strconv.Atoi(c.Query("month"))

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		tDate := time.Now()

		var AppFrom sql.NullTime
		var AppEnd sql.NullTime
		fm, err := time.Parse(time.DateOnly, strings.Split(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		AppFrom.Scan(fm)

		me, err := time.Parse(time.DateOnly, strings.Split(time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
		if err != nil {
			// fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		// fmt.Println("param = ", me)

		AppEnd.Scan(me)

		param := db.GetClaimedJobByCmpParams{
			Belongcmp:      bcmp.(int64),
			ApprovedDate:   AppFrom,
			ApprovedDate_2: AppEnd,
		}

		//endstart
		res, err := u.svc.JobsServ.GetClaimedJobByCmp(param)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, res)
		return

	} else {
		var info db.GetDriverRow
		if qid != "" {
			id, err := strconv.Atoi(qid)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			info, err = u.svc.UserServ.GetDriverInfo(int64(id))

			if err != nil {
				// fmt.Println("it is: ", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
		if (role == 300) || (role == 200 && info.Belongcmp == bcmp) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//start
		year, err := strconv.Atoi(c.Query("year"))

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		month, err := strconv.Atoi(c.Query("month"))

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		tDate := time.Now()

		var AppFrom sql.NullTime
		var AppEnd sql.NullTime
		fm, err := time.Parse(time.DateOnly, strings.Split(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		AppFrom.Scan(fm)

		me, err := time.Parse(time.DateOnly, strings.Split(time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
		if err != nil {
			// fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		// fmt.Println("param = ", me)

		AppEnd.Scan(me)

		param := db.GetClaimedJobByDriverIDParams{
			ID:             int64(uid),
			ApprovedDate:   AppFrom,
			ApprovedDate_2: AppEnd,
		}
		//endstart

		res, err := u.svc.JobsServ.GetClaimedJobByDriverID(param)
		if err != nil {

			// fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, res)
		return

	}
}

func (u *JobsCtrlImpl) GetClaimedJobByID(c *gin.Context) {

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

	res, err := u.svc.JobsServ.GetClaimedJobByID(int64(id))

	if err != nil {
		// fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, res)

}

func (u *JobsCtrlImpl) GetUserWithPendingJob(c *gin.Context) {
	// TODO: Finish This Part
	role := c.MustGet("Role").(int16)
	// id := c.MustGet("UserID")
	belongCmp := c.MustGet("belongCmp")

	var CmpID sql.NullInt64
	if role >= 200 {
		CmpID.Scan(belongCmp)
	} else {
		if c.Query("cmpID") != "" {
			CmpID.Scan(c.Query("cmpID"))
		}
	}

	res, err := u.svc.JobsServ.GetUserWithPendingJob(CmpID)

	if err != nil {
		// fmt.Println("err is ", err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}
	// fmt.Println(res)

	c.JSON(http.StatusOK, res)
}

func (u *JobsCtrlImpl) GetAllClaimedJobs(c *gin.Context) {
	role := c.MustGet("Role").(int16)
	id := c.MustGet("UserID")
	belongCmp := c.MustGet("belongCmp")

	var Uid sql.NullInt64
	if role == 300 {
		Uid.Scan(id)
	} else {
		if c.Query("uid") != "" {
			Uid.Scan(c.Query("uid"))
		}
	}
	var Jobid sql.NullInt64
	if c.Query("jobid") != "" {
		Jobid.Scan(c.Query("jobid"))
	}
	var cjID sql.NullInt64
	if c.Query("cjID") != "" {
		cjID.Scan(c.Query("cjID"))
	}
	var CmpID sql.NullInt64
	if role >= 200 {
		CmpID.Scan(belongCmp)
	} else {
		if c.Query("cmpID") != "" {
			CmpID.Scan(c.Query("cmpID"))
		}
	}
	var Ym sql.NullTime
	if c.Query("ym") != "" {
		d := c.Query("ym") + "-01"
		dt, err := time.Parse(time.DateOnly, d)
		if err != nil {
			// fmt.Println("err!! ", err)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		Ym.Scan(dt)
	}
	var Cat sql.NullString

	if c.Query("cat") != "" {
		Cat.Scan(c.Query("cat"))
	} else {
		Cat.Valid = false
	}

	param := db.GetAllClaimedJobsParams{
		Uid:   Uid,
		Jobid: Jobid,
		CmpID: CmpID,
		CjID:  cjID,
		Ym:    Ym,
		Cat:   Cat,
	}

	res, err := u.svc.JobsServ.GetAllClaimedJobs(param)

	if err != nil {
		// fmt.Println("err is ", err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}
	// fmt.Println(res)

	c.JSON(http.StatusOK, res)
}

func (u *JobsCtrlImpl) GetCJDate(c *gin.Context) {
	// protect
	sid := c.Query("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}
	res, err := u.svc.JobsServ.GetCJDate(int64(id))

	if err != nil {
		// fmt.Println(err)
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
		if role >= 200 {
			Belongcmp.Scan(belongCmp)
		}

		var Remaining sql.NullInt32
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

		param := db.GetAllJobsSuperParams{
			ID:        ID,
			FromLoc:   FromLoc,
			Mid:       Mid,
			ToLoc:     ToLoc,
			Belongcmp: Belongcmp,
			Remaining: Remaining,
			// CloseDateStart:        CloseDateStart,
			// CloseDateEnd:          CloseDateEnd,
			CreateDateStart: CreateDateStart,
			CreateDateEnd:   CreateDateEnd,
			// DeletedDateStart:      DeletedDateStart,
			// DeletedDateEnd:        DeletedDateEnd,
			LastModifiedDateStart: LastModifiedDateStart,
			LastModifiedDateEnd:   LastModifiedDateEnd,
		}
		res, err := u.svc.JobsServ.GetAllJobsSuper(param)

		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, res)

	} else {
		fmt.Println("Here")
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
	res, err, _ := u.svc.JobsServ.ClaimJob(param)

	if err != nil {
		// fmt.Print(err)
		if err.Error() == "already have ongoing job" {
			fmt.Println(res)
			res, err := u.svc.JobsServ.GetClaimedJobByID(res)
			if err != nil {
				fmt.Println(err)
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
	// if num == 1 {
	// 	SandMsg(1, 400, "check job open")
	// }
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
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	UserID := c.MustGet("UserID").(int)
	Role := c.MustGet("Role")
	cres, err := u.svc.JobsServ.GetClaimedJobByID(int64(id))

	if err != nil {
		fmt.Println(int64(id))
		fmt.Println(err)

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if Role != 100 && UserID != int(cres.Userid) {
		fmt.Println(err)

		c.AbortWithStatus(http.StatusBadRequest)
	}

	var reqBody apptypes.FinishClaimJobBodyT
	if err := c.Bind(&reqBody); err != nil {
		// fmt.Print("NOT OK")
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

		// fmt.Print(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)

}

func (u *JobsCtrlImpl) ApproveClaimedJob(c *gin.Context) {

	var reqBody apptypes.ApproveJob
	if err := c.Bind(&reqBody); err != nil {
		// fmt.Print("NOT OK")
		c.Abort()
		return
	}
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

	var Memo sql.NullString
	Memo.Scan(reqBody.Memo)

	var ApprovedBy sql.NullInt64
	ApprovedBy.Scan(UserID)

	param := db.ApproveFinishedJobParams{
		ID:         int64(id),
		ApprovedBy: ApprovedBy,
		Memo:       Memo,
	}
	err = u.svc.JobsServ.ApproveFinishedJob(param)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = u.svc.JobsServ.GetClaimedJobByID(int64(id))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// SandMsg(int(cJobInfo.Userid), 300, "Job "+strconv.Itoa(id)+" is approved")

	c.AbortWithStatus(http.StatusOK)

}
func (u *JobsCtrlImpl) ApproveClaimedJobs(c *gin.Context) {
	var reqBody struct {
		IDs  []int64 `json:"ids"`
		Memo string  `json:"memo"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	UserID := c.MustGet("UserID").(int)

	if len(reqBody.IDs) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}

	param := db.ApproveMultipleJobsParams{
		Ids:        reqBody.IDs,
		Memo:       sql.NullString{String: reqBody.Memo, Valid: reqBody.Memo != ""},
		ApprovedBy: sql.NullInt64{Int64: int64(UserID), Valid: true},
	}

	err := u.svc.JobsServ.ApproveFinishedJobs(param)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
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
		// fmt.Print(err)
		c.Abort()
		return
	}
	cJobRes, err := u.svc.JobsServ.GetClaimedJobByID(int64(id))

	if err != nil {
		// fmt.Print(err)
		c.Abort()
		return
	}
	if !(cJobRes.CreateDate.Add(time.Minute*10).After(time.Now()) && cJobRes.Userid == int64(UserID)) && !(res.Role <= int16(200)) {
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
		// SandMsg(1, 400, "check job open")
	}

	if err != nil {
		// fmt.Print(err)
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
		fmt.Println(err)
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

	// Jobdate, err := time.Parse(time.DateOnly, reqBody.Jobdate)
	// if err != nil {
	// 	c.Status(http.StatusBadRequest)
	// 	c.Abort()
	// 	return
	// }

	// var CloseDate sql.NullTime
	// if reqBody.CloseDate == "" {
	// 	CloseDate.Valid = false
	// } else {
	// 	CloseDate.Scan(reqBody.CloseDate)
	// }

	// var UserID sql.NullInt64
	// UserID.Scan(cuid)
	var estimatedN int
	if reqBody.Estimated == 0 {
		estimatedN = 2147483647
	} else {
		estimatedN = reqBody.Estimated
	}

	param := db.CreateJobParams{
		Fromloc:   reqBody.FromLoc,
		Mid:       Mid,
		Toloc:     reqBody.ToLoc,
		Price:     int32(reqBody.Price),
		Estimated: int32(estimatedN),
		Belongcmp: int64(reqBody.Belongcmp),
		Source:    reqBody.Source,
		// Jobdate:   Jobdate,
		Memo: Memo,
		// CloseDate: CloseDate,
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

	// Jobdate, err := time.Parse(time.DateOnly, strings.Split(reqBody.Jobdate, "T")[0])
	// if err != nil {
	// 	fmt.Print(err)
	// 	c.Status(http.StatusBadRequest)
	// 	c.Abort()
	// 	return
	// }

	// var CloseDate sql.NullTime
	// if reqBody.CloseDate != "" {

	// 	ct, err := time.Parse(time.DateOnly, strings.Split(reqBody.CloseDate, "T")[0])
	// 	if err != nil {

	// 		fmt.Print(2)
	// 		c.Status(http.StatusBadRequest)
	// 		c.Abort()
	// 		return
	// 	}
	// 	fmt.Print("ct: ", ct)
	// 	CloseDate.Scan(ct)

	// } else {
	// 	CloseDate.Valid = false
	// }

	var UserID sql.NullInt64
	UserID.Scan(cuid)
	var Remaining int
	if reqBody.Remaining == 0 {
		Remaining = 2147483647
	} else {
		Remaining = reqBody.Remaining
	}

	param := db.UpdateJobParams{
		ID:        int64(reqBody.ID),
		Fromloc:   reqBody.FromLoc,
		Mid:       Mid,
		Toloc:     reqBody.ToLoc,
		Price:     int32(reqBody.Price),
		Belongcmp: int64(reqBody.Belongcmp),
		Source:    reqBody.Source,
		// Jobdate:   Jobdate,
		Memo: Memo,
		// CloseDate: CloseDate,
		Remaining: int32(Remaining),
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
