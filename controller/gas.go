package controller

import (
	"database/sql"
	"fmt"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"main/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GasCtrl interface {
	CreateNewGas(c *gin.Context)
	DeleteRepair(c *gin.Context)
	ApproveRepair(c *gin.Context)
	GetRepair(c *gin.Context)
	GetRepairByID(c *gin.Context)
	GetRepairDate(c *gin.Context)
}

type GasCtrlImpl struct {
	svc *service.AppService
}

func (u *GasCtrlImpl) GetRepairDate(c *gin.Context) {
	// protect
	sid := c.Query("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}
	res, err := u.svc.RepairServ.GetRepairDate(int64(id))

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}
func (r *GasCtrlImpl) GetRepairByID(c *gin.Context) {

	// UserID := c.MustGet("UserID").(int)
	rid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	res, err := r.svc.RepairServ.GetRepairById(int64(rid))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, res)

}

func (r *GasCtrlImpl) GetRepair(c *gin.Context) {

	UserID := c.MustGet("UserID").(int)
	belongCmp := c.MustGet("belongCmp").(int64)
	// res, err := r.svc.UserServ.GetUserById(int64(UserID))

	// if err != nil {
	// 	c.Status(http.StatusBadRequest)
	// 	c.Abort()
	// 	return
	// }

	userRole := c.MustGet("Role").(int16)

	var Id sql.NullInt64
	if !(c.Query("id") == "") {
		Id.Scan(c.Query("id"))
	} else {
		Id.Valid = false
	}
	var Ym sql.NullTime
	if c.Query("ym") != "" {
		d := c.Query("ym") + "-01"
		dt, err := time.Parse(time.DateOnly, d)
		if err != nil {
			fmt.Println("err!! ", err)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		Ym.Scan(dt)
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
	if userRole >= 200 {
		BelongCmp.Scan(belongCmp)
	} else {
		if !(c.Query("belongCmp") == "") {
			BelongCmp.Scan(c.Query("belongCmp"))
		}
	}

	// var CreateDateStart sql.NullTime
	// if c.Query("CreateDateStart") == "" {
	// 	CreateDateStart.Valid = false
	// } else {
	// 	CreateDateStart.Scan(c.Query("CreateDateStart"))
	// }

	// var CreateDateEnd sql.NullTime
	// if c.Query("CreateDateEnd") == "" {
	// 	CreateDateEnd.Valid = false
	// } else {
	// 	CreateDateEnd.Scan(c.Query("CreateDateEnd"))
	// }

	// var DeletedDateStart sql.NullTime
	// if c.Query("DeletedDateStart") == "" {
	// 	DeletedDateStart.Valid = false
	// } else {
	// 	DeletedDateStart.Scan(c.Query("DeletedDateStart"))
	// }
	//
	// var DeletedDateEnd sql.NullTime
	// if c.Query("DeletedDateEnd") == "" {
	// 	DeletedDateEnd.Valid = false
	// } else {
	// 	DeletedDateEnd.Scan(c.Query("DeletedDateEnd"))
	// }

	// var LastModifiedDateStart sql.NullTime
	// if c.Query("LastModifiedDateStart") == "" {
	// 	LastModifiedDateStart.Valid = false
	// } else {
	// 	LastModifiedDateStart.Scan(c.Query("LastModifiedDateStart"))
	// }

	// var LastModifiedDateEnd sql.NullTime
	// if c.Query("LastModifiedDateEnd") == "" {
	// 	LastModifiedDateEnd.Valid = false
	// } else {
	// 	LastModifiedDateEnd.Scan(c.Query("LastModifiedDateEnd"))
	// }
	var Cat sql.NullString

	if c.Query("cat") != "" {
		Cat.Scan(c.Query("cat"))
	} else {
		Cat.Valid = false
	}
	param := db.GetRepairParams{
		ID:        Id,
		DriverID:  DriverId,
		Name:      Name,
		Belongcmp: BelongCmp,
		Cat:       Cat,
		Ym:        Ym,
	}
	repairRes, err := r.svc.RepairServ.GetRepair(param)

	if err != nil && err != sql.ErrNoRows {
		fmt.Print("err", err)
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

func (r *GasCtrlImpl) ApproveRepair(c *gin.Context) {

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

	// res, err := r.svc.RepairServ.GetRepairById(int64(id))
	_, err = r.svc.RepairServ.GetRepairById(int64(id))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	// SandMsg(int(res.Uid), 300, "Repair "+strconv.Itoa(id)+" is approved")

	c.Status(http.StatusOK)
	c.Abort()
}

func (r *GasCtrlImpl) DeleteRepair(c *gin.Context) {

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

func (r *GasCtrlImpl) CreateNewRepair(c *gin.Context) {
	// TODO: Repair Info In Body

	var reqBody apptypes.NewRepairBodyT
	err := c.Bind(&reqBody)

	if err != nil {
		c.Abort()
		return
	}

	cuid := c.MustGet("UserID").(int)

	var pic sql.NullString
	if reqBody.RepairPic != nil {
		path, uuid, err := utils.GenPicRoute(reqBody.RepairPic.Header["Content-Type"][0])
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = c.SaveUploadedFile(reqBody.RepairPic, path)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		pic.Scan(uuid)

	}

	param := db.CreateNewRepairParams{
		Driverid: int64(cuid),
		Pic:      pic,
		Place:    reqBody.Place,
	}

	rID, err := r.svc.RepairServ.NewRepair(param)

	if err != nil {
		fmt.Println("err: ", err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	for _, item := range reqBody.Repairinfo {
		_, err := r.svc.RepairServ.NewRepairInfo(item)
		if err != nil {
			fmt.Println("err: ", err)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"res": rID})
}

func GasCtrlInit(svc *service.AppService) *GasCtrlImpl {
	return &GasCtrlImpl{
		svc: svc,
	}
}
