package controller

import (
	"database/sql"
	"fmt"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CmpCtrl interface {
	RegisterCmp(c *gin.Context)
	DeleteCmp(c *gin.Context)
	UpdateCmp(c *gin.Context)
	GetCmp(c *gin.Context)
	GetAllCmp(c *gin.Context)
	GetJobCmp(c *gin.Context)
}

type CmpCtrlImpl struct {
	svc *service.AppService
}

func (u *CmpCtrlImpl) GetCmp(c *gin.Context) {

	var reqBody apptypes.GetCmpT

	err := c.BindJSON(&reqBody)

	if err != nil {
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	data, err := u.svc.CmpServ.GetCmp(int64(reqBody.Id))

	if err != nil {
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, data)
}

func (u *CmpCtrlImpl) GetAllCmp(c *gin.Context) {

	cmpList, err := u.svc.CmpServ.GetAllCmp()

	if err != nil {
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, cmpList)
}

func (u *CmpCtrlImpl) GetJobCmp(c *gin.Context) {
	year, err := strconv.Atoi(c.Query("year"))
	bcmp := c.MustGet("belongCmp")
	role := c.MustGet("Role").(int16)
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
	// TODO: test here
	var cmpid sql.NullInt64
	if role <= 100 {
		cmpid.Valid = false
	} else {

		cmpid.Scan(bcmp)

	}

	param := db.GetJobCmpParams{
		ApprovedDate:   AppFrom,
		ApprovedDate_2: AppEnd,
		CmpId:          cmpid,
	}
	// fmt.Println("param = ", param)

	res, err := u.svc.CmpServ.GetJobCmp(param)

	// cmpList, err := u.svc.CmpServ.GetAllCmp()

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *CmpCtrlImpl) UpdateCmp(c *gin.Context) {
	var reqBody apptypes.UpdateCmpT

	if err := c.BindJSON(&reqBody); err != nil {
		return
	}

	param := db.UpdateCmpParams{
		ID:   int64(reqBody.Id),
		Name: reqBody.CmpName,
	}

	err := u.svc.CmpServ.UpdateCmp(param)

	if err != nil {
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (u *CmpCtrlImpl) RegisterCmp(c *gin.Context) {
	var reqBody apptypes.RegisterCmpT

	if err := c.BindJSON(&reqBody); err != nil {
		return
	}

	newid, err := u.svc.CmpServ.NewCmp(reqBody.CmpName)

	if err != nil {
		// fmt.Print(err)
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": newid})
}

func (u *CmpCtrlImpl) DeleteCmp(c *gin.Context) {
	var reqBody apptypes.DeleteCmpBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		return
	}

	err = u.svc.CmpServ.DeleteCmp(int64(reqBody.ToDeleteCmpId))

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
	}

	c.Status(http.StatusOK)
	c.Abort()
}

func CmpCtrlInit(svc *service.AppService) *CmpCtrlImpl {
	return &CmpCtrlImpl{
		svc: svc,
	}
}
