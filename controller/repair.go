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

type RepairCtrl interface {
	CreateNewRepair(c *gin.Context)
	DeleteRepair(c *gin.Context)
	ApproveRepair(c *gin.Context)
	GetRepair(c *gin.Context)
}

type RepairCtrlImpl struct {
	svc *service.AppService
}

func (r *RepairCtrlImpl) GetRepair(c *gin.Context) {

	UserID := c.MustGet("UserID").(int)

	res, err := r.svc.UserServ.GetUserById(int64(UserID))

	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	userRole := res.Role

	var Id sql.NullInt64
	if !(c.Query("id") == "") {
		Id.Scan(c.Query("id"))
	} else {
		Id.Valid = false
	}

	var DriverId sql.NullInt64
	// driver
	if userRole >= 300 {
		DriverId.Scan(UserID)
	} else {
		if !(c.Query("driverid") == "") {
			DriverId.Scan(c.Query("driverid"))
		}
	}

	var Name sql.NullString
	if !(c.Query("name") == "") {
		Name.Scan(c.Query("name"))
	} else {
		Name.Valid = false
	}

	// cmp and up
	var BelongCmp sql.NullInt64
	if userRole > 100 {
		BelongCmp.Scan(res.Belongcmp)
	} else {
		if !(c.Query("belongCmp") == "") {
			BelongCmp.Scan(c.Query("belongCmp"))
		}
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

	param := db.GetRepairParams{
		ID:                    Id,
		DriverID:              DriverId,
		Name:                  Name,
		Belongcmp:             BelongCmp,
		CreateDateStart:       CreateDateStart,
		CreateDateEnd:         CreateDateEnd,
		DeletedDateStart:      DeletedDateStart,
		DeletedDateEnd:        DeletedDateEnd,
		LastModifiedDateStart: LastModifiedDateStart,
		LastModifiedDateEnd:   LastModifiedDateEnd,
	}
	repairRes, err := r.svc.RepairServ.GetRepair(param)

	if err != nil && err != sql.ErrNoRows {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": repairRes})

}

func (r *RepairCtrlImpl) ApproveRepair(c *gin.Context) {

	param_id := c.Param("id")

	if param_id == "" {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	id, err := strconv.Atoi(param_id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	err = r.svc.RepairServ.ApproveRepair(int64(id))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	c.Status(http.StatusOK)
	c.Abort()
}

func (r *RepairCtrlImpl) DeleteRepair(c *gin.Context) {

	param_id := c.Param("id")

	if param_id == "" {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	id, err := strconv.Atoi(param_id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	err = r.svc.RepairServ.DeleteRepair(int64(id))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	c.Status(http.StatusOK)
	c.Abort()
}

func (r *RepairCtrlImpl) CreateNewRepair(c *gin.Context) {
	repairType := c.Query("type")

	if repairType != "gas" && repairType != "maintain" {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	var reqBody apptypes.NewRepairBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		c.Abort()
		return
	}

	cuid := c.MustGet("UserID").(int)

	param := db.CreateNewRepairParams{
		Type:       repairType,
		Driverid:   int64(cuid),
		Repairinfo: reqBody.Repairinfo,
	}

	res, err := r.svc.RepairServ.NewRepair(param)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": res})
}

func RepairCtrlInit(svc *service.AppService) *RepairCtrlImpl {
	return &RepairCtrlImpl{
		svc: svc,
	}
}
