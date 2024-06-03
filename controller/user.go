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

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserCtrl interface {
	RegisterUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetUserById(c *gin.Context)
	GetUserList(c *gin.Context)
	UpdateUser(c *gin.Context)
	UpdatePassword(c *gin.Context)
	UpdateDriverPic(c *gin.Context)
	ApproveUser(c *gin.Context)
	Me(c *gin.Context)
	ResetPassword(c *gin.Context)
}

type UserCtrlImpl struct {
	svc *service.AppService
}

func (u *UserCtrlImpl) Me(c *gin.Context) {
	id := c.MustGet("UserID").(int)

	res, err := u.svc.UserServ.GetUserById(int64(id))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if res.Role >= 300 {
		res, err := u.svc.UserServ.GetDriverInfo(int64(id))

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
func (u *UserCtrlImpl) ApproveUser(c *gin.Context) {
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

	err = u.svc.UserServ.ApproveDriver(int64(id))

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	SandMsg(id, 200, "driver Approved")
	c.AbortWithStatus(http.StatusOK)
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

func (u *UserCtrlImpl) ResetPassword(c *gin.Context) {
	role := c.MustGet("Role").(int16)

	var reqBody apptypes.ResetPasswordBodyT
	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if role != 100 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := u.svc.UserServ.GetUserById(int64(reqBody.Id))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// fmt.Println(user.Phonenum)
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Phonenum.(string)), bcrypt.MinCost)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	param := db.UpdateUserPasswordParams{
		ID:  int64(reqBody.Id),
		Pwd: string(hash),
	}
	err = u.svc.UserServ.UpdatePassword(param)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (u *UserCtrlImpl) UpdatePassword(c *gin.Context) {
	cuid := c.MustGet("UserID").(int)
	// role := c.MustGet("Role").(int16)

	var reqBody apptypes.UpdatePasswordBodyT
	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if cuid != reqBody.Id {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := u.svc.UserServ.GetUserById(int64(reqBody.Id))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res, err := u.svc.UserServ.HaveUser(user.Phonenum)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(res.Pwd), []byte(reqBody.OldPwd)); err != nil {
		c.Status(http.StatusNotAcceptable)
		c.Abort()
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Pwd), bcrypt.MinCost)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	param := db.UpdateUserPasswordParams{
		ID:  int64(cuid),
		Pwd: string(hash),
	}

	err = u.svc.UserServ.UpdatePassword(param)

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
			Belongcmp: int64(reqBody.BelongCmp),
			Role:      int16(role),
		}
		err = u.svc.UserServ.UpdateUser(param)
	} else if reqBody.Role == "superAdmin" {

		role = 100
		param = db.UpdateUserParams{
			ID:        int64(id),
			Phonenum:  reqBody.PhoneNum,
			Name:      reqBody.Name,
			Belongcmp: int64(reqBody.BelongCmp),
			Role:      int16(role),
		}
	} else {
		role = 300
		param = db.UpdateUserParams{
			ID:        int64(id),
			Phonenum:  reqBody.PhoneNum,
			Name:      reqBody.Name,
			Belongcmp: int64(reqBody.BelongCmp),
			Role:      int16(role),
		}
		driverParam := db.UpdateDriverParams{
			ID: int64(id),
			// Percentage:       int32(reqBody.DriverInfo.Percentage),
			Platenum:         reqBody.DriverInfo.PlateNum,
			Nationalidnumber: reqBody.DriverInfo.NationalIdNumber,
		}

		err = u.svc.UserServ.UpdateDriver(driverParam, param)

	}

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (u *UserCtrlImpl) RegisterUser(c *gin.Context) {

	userType := c.Query("userType")

	var newid int64

	switch userType {
	case "cmpAdmin":
		var reqBody apptypes.RegisterCmpAdminBodyT

		if err := c.BindJSON(&reqBody); err != nil {
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.PhoneNum), bcrypt.MinCost)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		param := db.CreateUserParams{
			Pwd:       string(hash),
			Name:      reqBody.Name,
			Belongcmp: int64(reqBody.BelongCmp),
			Phonenum:  reqBody.PhoneNum,
			Role:      200,
		}
		newid, err = u.svc.UserServ.RegisterCmpAdmin(param)
		if err != nil {
			c.Status(http.StatusConflict)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": newid})

	case "driver":

		var reqBody apptypes.RegisterDriverBodyT

		if err := c.BindJSON(&reqBody); err != nil {
			fmt.Println(err)
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.PhoneNum), bcrypt.MinCost)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		param := db.CreateUserParams{
			Pwd:       string(hash),
			Name:      reqBody.Name,
			Belongcmp: int64(reqBody.BelongCmp),
			Phonenum:  reqBody.PhoneNum,
			Role:      300,
		}

		newid, err = u.svc.UserServ.RegisterDriver(param, reqBody.DriverInfo.NationalIdNumber, reqBody.DriverInfo.PlateNum)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": newid})

	default:

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

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
		Name.Scan("%" + c.Query("Name") + "%")

	}

	var BelongCmp sql.NullInt64
	if c.Query("BelongCmp") == "" {
		BelongCmp.Valid = false
	} else {
		BelongCmp.Scan(c.Query("BelongCmp"))
	}

	var BelongCmpName sql.NullString
	if c.Query("BelongCmpName") == "" {
		BelongCmpName.Valid = false
	} else {
		BelongCmpName.Scan("%" + c.Query("BelongCmpName") + "%")
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
		BelongcmpName:         BelongCmpName,
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
	c.JSON(http.StatusOK, res)
}

func UserCtrlInit(svc *service.AppService) *UserCtrlImpl {
	return &UserCtrlImpl{
		svc: svc,
	}
}
