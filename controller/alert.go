package controller

import (
	"database/sql"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlertCtrl interface {
	CreateAlert(c *gin.Context)
	DeleteAlert(c *gin.Context)
	CheckNewAlert(c *gin.Context)
	GetAlert(c *gin.Context)
}

type AlertCtrlImpl struct {
	svc *service.AppService
}

func (a *AlertCtrlImpl) CheckNewAlert(c *gin.Context) {
	cuid := c.MustGet("UserID").(int)

	res, err := a.svc.AlertServ.HaveNewAlert(int64(cuid))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (a *AlertCtrlImpl) GetAlert(c *gin.Context) {
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
	res, err := a.svc.AlertServ.GetAlert(param)

	if err != nil && err != sql.ErrNoRows {
		// fmt.Print(err)
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
		a.svc.AlertServ.UpdateLastAlert(updateParam)
	}

	c.JSON(http.StatusOK, gin.H{"res": res})
}

func (a *AlertCtrlImpl) DeleteAlert(c *gin.Context) {

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

	err = a.svc.AlertServ.DeleteAlert(int64(id))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (a *AlertCtrlImpl) CreateAlert(c *gin.Context) {
	role := c.MustGet("Role").(int16)
	bcmp := c.MustGet("belongCmp").(int64)

	var reqBody apptypes.CreateAlertBodyT
	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var cmpID int64
	cmpID = int64(reqBody.BelongCmp)

	if role >= 200 {
		// send cmp to own driver
		cmpID = bcmp
		param := db.CreateAlertParams{
			Alert:     reqBody.Alert,
			Belongcmp: cmpID,
		}
		res, err := a.svc.AlertServ.CreateAlert(param)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		SendMessageCmpToDriver(reqBody.BelongCmp, 100, reqBody.Alert)

		c.JSON(http.StatusOK, gin.H{"res": res})
		return
	}

	param := db.CreateAlertParams{
		Alert:     reqBody.Alert,
		Belongcmp: cmpID,
	}
	res, err := a.svc.AlertServ.CreateAlert(param)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	SendMessageByCmp(reqBody.BelongCmp, 100, reqBody.Alert)

	c.JSON(http.StatusOK, gin.H{"res": res})
}

func AlertCtrlInit(svc *service.AppService) *AlertCtrlImpl {
	return &AlertCtrlImpl{
		svc: svc,
	}
}
