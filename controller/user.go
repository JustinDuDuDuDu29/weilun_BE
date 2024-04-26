package controller

import (
	"database/sql"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"main/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserCtrl interface {
	RegisterUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetUserById(c *gin.Context)
	GetUserList(c *gin.Context)
	UpdateUser(c *gin.Context)
	UpdatePassword(c *gin.Context)
	UpdateDriverPic(c *gin.Context)
}

type UserCtrlImpl struct {
	svc *service.AppService
}

func (u *UserCtrlImpl) UpdateDriverPic(c *gin.Context) {

	var reqBody apptypes.UpdateDriverPic
	if err := c.ShouldBind(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var DriverLicense sql.NullString
	if reqBody.DriverLicense != nil {
		path, uuid, err := utils.GenPicRoute(reqBody.DriverLicense.Header["Content-Type"][0])
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = c.SaveUploadedFile(reqBody.DriverLicense, path)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		DriverLicense.Scan(uuid)

	}

	var Insurances sql.NullString
	if reqBody.Insurances != nil {
		path, uuid, err := utils.GenPicRoute(reqBody.Insurances.Header["Content-Type"][0])
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = c.SaveUploadedFile(reqBody.Insurances, path)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		Insurances.Scan(uuid)

	}

	var Registration sql.NullString
	if reqBody.Registration != nil {
		path, uuid, err := utils.GenPicRoute(reqBody.Registration.Header["Content-Type"][0])
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = c.SaveUploadedFile(reqBody.Registration, path)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		Registration.Scan(uuid)

	}

	var TruckLicense sql.NullString
	if reqBody.TruckLicense != nil {
		path, uuid, err := utils.GenPicRoute(reqBody.TruckLicense.Header["Content-Type"][0])
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = c.SaveUploadedFile(reqBody.TruckLicense, path)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		TruckLicense.Scan(uuid)

	}

	param := db.UpdateDriverPicParams{
		ID:            int64(c.MustGet("UserID").(int)),
		Insurances:    Insurances,
		Registration:  Registration,
		Driverlicense: DriverLicense,
		Trucklicense:  TruckLicense,
	}
	err := u.svc.UserServ.UpdateDriverPic(param)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (u *UserCtrlImpl) UpdatePassword(c *gin.Context) {
	cuid := c.MustGet("UserID").(int)

	var reqBody apptypes.UpdatePasswordBodyT
	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	param := db.UpdateUserPasswordParams{
		ID:  int64(cuid),
		Pwd: reqBody.Pwd,
	}
	err := u.svc.UserServ.UpdatePassword(param)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (u *UserCtrlImpl) UpdateUser(c *gin.Context) {
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

	// cuid := c.MustGet("UserID").(int)
	// role := c.MustGet("Role").(int)
	// bcmp := c.MustGet("belongCmp").(int)

	res, err := u.svc.UserServ.GetUserById(int64(id))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var reqBody apptypes.RegisterDriverBodyT
	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var role int
	var param db.UpdateUserParams

	if reqBody.Role == "cmpAdmin" {
		role = 200
		param = db.UpdateUserParams{
			ID:        int64(id),
			Phonenum:  reqBody.PhoneNum,
			Name:      reqBody.Name,
			Belongcmp: res.Belongcmp,
			Role:      int16(role),
		}
		u.svc.UserServ.UpdateUser(param)
	} else {
		role = 300
		param = db.UpdateUserParams{
			ID:        int64(id),
			Phonenum:  reqBody.PhoneNum,
			Name:      reqBody.Name,
			Belongcmp: res.Belongcmp,
			Role:      int16(role),
		}
		driverParam := db.UpdateDriverParams{
			ID:         int64(id),
			Percentage: int16(reqBody.DriverInfo.Percentage),
		}

		u.svc.UserServ.UpdateDriver(driverParam, param)

	}
}

func (u *UserCtrlImpl) RegisterUser(c *gin.Context) {
	userType := c.Query("userType")

	var newid int64
	var err error

	switch userType {
	case "cmpAdmin":
		var reqBody apptypes.RegisterCmpAdminBodyT

		if err := c.BindJSON(&reqBody); err != nil {
			return
		}

		param := db.CreateUserParams{
			Pwd:       reqBody.PhoneNum,
			Name:      reqBody.Name,
			Belongcmp: int64(reqBody.BelongCmp),
			Phonenum:  reqBody.PhoneNum,
			Role:      200,
		}
		newid, err = u.svc.UserServ.RegisterCmpAdmin(param)

	case "driver":
		var reqBody apptypes.RegisterDriverBodyT

		if err := c.BindJSON(&reqBody); err != nil {
			return
		}
		param := db.CreateUserParams{
			Pwd:       reqBody.PhoneNum,
			Name:      reqBody.Name,
			Belongcmp: int64(reqBody.BelongCmp),
			Phonenum:  reqBody.PhoneNum,
			Role:      300,
		}

		newid, err = u.svc.UserServ.RegisterDriver(param, reqBody.DriverInfo.Percentage, reqBody.DriverInfo.NationalIdNumber)

	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err != nil {
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": newid})
}

func (u *UserCtrlImpl) DeleteUser(c *gin.Context) {
	var reqBody apptypes.DeleteUserBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		return
	}

	err = u.svc.UserServ.DeleteUser(int64(reqBody.ToDeleteUserId))

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.Status(http.StatusOK)
	c.Abort()
}

func (u *UserCtrlImpl) GetUserList(c *gin.Context) {

	var ID sql.NullInt64
	if c.Query("Id") == "" {
		ID.Valid = false
	} else {
		ID.Scan(c.Query("Id"))
	}

	var PhoneNum sql.NullString
	if c.Query("PhoneNum") == "" {
		PhoneNum.Valid = false
	} else {
		PhoneNum.Scan(c.Query("PhoneNum"))
	}

	var Name sql.NullString
	if c.Query("Name") == "" {
		Name.Valid = false
	} else {
		Name.Scan(c.Query("Name"))

	}

	var BelongCmp sql.NullInt64
	if c.Query("BelongCmp") == "" {
		BelongCmp.Valid = false
	} else {
		BelongCmp.Scan(c.Query("BelongCmp"))
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

	param := db.GetUserListParams{
		ID:                    ID,
		PhoneNum:              PhoneNum,
		Name:                  Name,
		Belongcmp:             BelongCmp,
		CreateDateStart:       CreateDateStart,
		CreateDateEnd:         CreateDateEnd,
		DeletedDateStart:      DeletedDateStart,
		DeletedDateEnd:        DeletedDateEnd,
		LastModifiedDateStart: LastModifiedDateStart,
		LastModifiedDateEnd:   LastModifiedDateEnd,
	}

	res, err := u.svc.UserServ.GetUserList(param)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (u *UserCtrlImpl) GetUserById(c *gin.Context) {
	var Id sql.NullInt64
	if c.Param("id") == "" {
		Id.Valid = false
	} else {
		Id.Scan(c.Param("id"))
	}

	res, err := u.svc.UserServ.GetUserById(Id.Int64)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": res})
}

func UserCtrlInit(svc *service.AppService) *UserCtrlImpl {
	return &UserCtrlImpl{
		svc: svc,
	}
}
