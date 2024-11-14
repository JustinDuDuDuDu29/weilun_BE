package controller

import (
	"database/sql"
	"encoding/json"
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

type RepairCtrl interface {
	CreateNewRepair(c *gin.Context)
	DeleteRepair(c *gin.Context)
	ApproveRepair(c *gin.Context)
	GetRepair(c *gin.Context)
	GetRepairByID(c *gin.Context)
	GetRepairDate(c *gin.Context)
	UpdateItem(c *gin.Context)
}

type RepairCtrlImpl struct {
	svc *service.AppService
}

func (u *RepairCtrlImpl) GetRepairDate(c *gin.Context) {
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
		// fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *RepairCtrlImpl) GetRepairByID(c *gin.Context) {

	// UserID := c.MustGet("UserID").(int)
	rid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	param := db.GetRepairParams{
		ID:        sql.NullInt64{Int64: int64(rid), Valid: true},
		DriverID:  sql.NullInt64{Int64: int64(-1), Valid: false},
		Name:      sql.NullString{String: "", Valid: false},
		Belongcmp: sql.NullInt64{Int64: int64(-1), Valid: false},
		Cat:       sql.NullString{String: "", Valid: false},
		Ym:        sql.NullString{String: "", Valid: false},
	}
	res, err := r.svc.RepairServ.GetRepair(param)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, res)

}

func (r *RepairCtrlImpl) GetRepair(c *gin.Context) {

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
			// fmt.Println("err!! ", err)
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

	// res, err := r.svc.RepairServ.GetRepairById(int64(id))
	param := db.GetRepairParams{
		ID:        sql.NullInt64{Int64: int64(id), Valid: true},
		DriverID:  sql.NullInt64{Int64: int64(-1), Valid: false},
		Name:      sql.NullString{String: "", Valid: false},
		Belongcmp: sql.NullInt64{Int64: int64(-1), Valid: false},
		Cat:       sql.NullString{String: "", Valid: false},
		Ym:        sql.NullString{String: "", Valid: false},
	}
	_, err = r.svc.RepairServ.GetRepair(param)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	// SandMsg(int(res.Uid), 300, "Repair "+strconv.Itoa(id)+" is approved")

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

func (r *RepairCtrlImpl) UpdateItem(c *gin.Context) {
	// bodyBytes, _ := io.ReadAll(c.Request.Body)
	// fmt.Printf("Body: %s\n", string(bodyBytes))
	var reqBody apptypes.UpdatedItems
	err := c.BindJSON(&reqBody)

	if err != nil {
		fmt.Println("out here: ", err)
		c.Abort()
		return
	}
	// err = json.Unmarshal([]byte(reqBody.UpdatedItems), &data)

	//danger???
	//TODO: check this part
	for _, item := range reqBody.UpdatedItems {
		price, err := strconv.Atoi(item.Price)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusBadRequest)
			c.Abort()
			return
		}

		data := db.UpdateItemParams{
			ID:         int64(item.Id),
			Totalprice: int64(price),
		}
		// fmt.Println(data)
		err = r.svc.RepairServ.UpdateItem(data)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusBadRequest)
			c.Abort()
			return
		}
	}

	c.Status(http.StatusOK)
	c.Abort()
}

func (r *RepairCtrlImpl) CreateNewRepair(c *gin.Context) {
	// TODO: Repair Info In Body

	var reqBody apptypes.NewRepairBodyT
	err := c.Bind(&reqBody)

	if err != nil {
		// fmt.Print("out here")
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
		// fmt.Println("err: ", err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}
	data := []db.CreateNewRepairInfoParams{}
	err = json.Unmarshal([]byte(reqBody.Repairinfo), &data)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, item := range data {
		// subdata := db.CreateDriverInfoParams{}

		// err = json.Unmarshal(([]byte(item)), &subdata)

		item.Repairid = rID
		fmt.Println(item)
		_, err := r.svc.RepairServ.NewRepairInfo(item)
		if err != nil {
			fmt.Println("1err: ", err)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"res": rID})
}

func RepairCtrlInit(svc *service.AppService) *RepairCtrlImpl {
	return &RepairCtrlImpl{
		svc: svc,
	}
}
