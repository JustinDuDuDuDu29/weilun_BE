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

type GasCtrl interface {
	CreateNewGas(c *gin.Context)
	DeleteGas(c *gin.Context)
	ApproveGas(c *gin.Context)
	GetGas(c *gin.Context)
	GetGasByID(c *gin.Context)
	GetGasDate(c *gin.Context)
	GetGasCmpUser(c *gin.Context)
	UpdateGas(c *gin.Context)
}

type GasCtrlImpl struct {
	svc *service.AppService
}

func (r *GasCtrlImpl) UpdateGas(c *gin.Context) {
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
		price, err := strconv.Atoi(item.Totalprice)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusBadRequest)
			c.Abort()
			return
		}

		data := db.UpdateGasParams{
			ID:         int64(item.Id),
			Totalprice: int64(price),
		}
		// fmt.Println(data)
		err = r.svc.GasServ.UpdateGas(data)
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

func (r *GasCtrlImpl) GetGasCmpUser(c *gin.Context) {

	// UserID := c.MustGet("UserID").(int)
	belongCmp := c.MustGet("belongCmp").(int64)
	userRole := c.MustGet("Role").(int16)

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
	// param := db.GetGasCmpUserParams{
	// Belongcmp: BelongCmp,
	// 	Cat:       Cat,
	// }
	repairRes, err := r.svc.GasServ.GetGasCmpUser(BelongCmp)

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
	c.JSON(http.StatusOK, repairRes)

}
func (u *GasCtrlImpl) GetGasDate(c *gin.Context) {
	// protect
	sid := c.Query("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}
	res, err := u.svc.GasServ.GetGasDate(int64(id))

	if err != nil {
		// fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}
func (r *GasCtrlImpl) GetGasByID(c *gin.Context) {

	// UserID := c.MustGet("UserID").(int)
	rid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	param := db.GetGasParams{
		ID:        sql.NullInt64{Int64: int64(rid), Valid: true},
		DriverID:  sql.NullInt64{Int64: int64(-1), Valid: false},
		Name:      sql.NullString{String: "", Valid: false},
		Belongcmp: sql.NullInt64{Int64: int64(-1), Valid: false},
		Cat:       sql.NullString{String: "", Valid: false},
		Ym:        sql.NullString{String: "", Valid: false},
	}
	res, err := r.svc.GasServ.GetGas(param)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, res)

}

func (r *GasCtrlImpl) GetGas(c *gin.Context) {

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
	param := db.GetGasParams{
		ID:        Id,
		DriverID:  DriverId,
		Name:      Name,
		Belongcmp: BelongCmp,
		Cat:       Cat,
		Ym:        Ym,
	}
	repairRes, err := r.svc.GasServ.GetGas(param)

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

func (r *GasCtrlImpl) ApproveGas(c *gin.Context) {

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

	err = r.svc.GasServ.ApproveGas(int64(id))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	param := db.GetRepairParams{
		ID:        sql.NullInt64{Int64: int64(id), Valid: true},
		DriverID:  sql.NullInt64{Int64: int64(-1), Valid: false},
		Name:      sql.NullString{String: "", Valid: false},
		Belongcmp: sql.NullInt64{Int64: int64(-1), Valid: false},
		Cat:       sql.NullString{String: "", Valid: false},
		Ym:        sql.NullString{String: "", Valid: false},
	}
	_, err = r.svc.RepairServ.GetRepair(param)
	// res, err := r.svc.GasServ.GetGasById(int64(id))
	// _, err = r.svc.GasServ.GetGasById(int64(id))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	// SandMsg(int(res.Uid), 300, "Gas "+strconv.Itoa(id)+" is approved")

	c.Status(http.StatusOK)
	c.Abort()
}

func (r *GasCtrlImpl) DeleteGas(c *gin.Context) {

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

	err = r.svc.GasServ.DeleteGas(int64(id))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	c.Status(http.StatusOK)
	c.Abort()
}

func (r *GasCtrlImpl) CreateNewGas(c *gin.Context) {
	// TODO: Gas Info In Body

	var reqBody apptypes.NewGasBodyT
	err := c.Bind(&reqBody)

	if err != nil {
		fmt.Println(err)
		c.Abort()
		return
	}

	cuid := c.MustGet("UserID").(int)

	var pic sql.NullString
	if reqBody.GasPic != nil {
		path, uuid, err := utils.GenPicRoute(reqBody.GasPic.Header["Content-Type"][0])
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = c.SaveUploadedFile(reqBody.GasPic, path)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		pic.Scan(uuid)

	}

	param := db.CreateNewGasParams{
		Driverid: int64(cuid),
		Pic:      pic,
		// Place:    reqBody.Place,
	}

	rID, err := r.svc.GasServ.NewGas(param)

	if err != nil {
		// fmt.Println("err: ", err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	data := []db.CreateNewGasInfoParams{}
	err = json.Unmarshal([]byte(reqBody.Gasinfo), &data)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, item := range data {
		// subdata := db.CreateDriverInfoParams{}

		// err = json.Unmarshal(([]byte(item)), &subdata)

		item.Gasid = rID
		_, err := r.svc.GasServ.NewGasInfo(item)
		if err != nil {
			// fmt.Println("err: ", err)
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
