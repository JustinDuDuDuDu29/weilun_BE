package controller

import (
	"fmt"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserCtrl interface {
	RegisterUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetUserById(c *gin.Context)
	GetUserList(c *gin.Context)
}

type UserCtrlImpl struct {
	svc *service.AppService
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
			Belongcmp: pgtype.Int8{Int64: int64(reqBody.BelongCmp), Valid: true},
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
			Belongcmp: pgtype.Int8{Int64: int64(reqBody.BelongCmp), Valid: true},
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
	return
}

func (u *UserCtrlImpl) GetUserList(c *gin.Context) {
	userId := c.Query("id")
	id := &pgtype.Numeric{}
	err := id.Scan(userId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	phoneNum := c.Query("phoneNum")
	name := c.Query("name")

	cmp := c.Query("belongCmp")
	belongCmpInt64 := &pgtype.Numeric{}
	err = belongCmpInt64.Scan(cmp)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	belongCmp, err := belongCmpInt64.Int64Value()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	create_date_StartStr := c.Query("create_date_Start")
	create_date_Start := &pgtype.Timestamp{}
	err = create_date_Start.Scan(create_date_StartStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	create_date_EndStr := c.Query("create_date_end")
	create_date_End := &pgtype.Timestamp{}
	err = create_date_End.Scan(create_date_EndStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	deleted_date_StartStr := c.Query("deleted_date_Start")
	deleted_date_Start := &pgtype.Timestamp{}
	err = create_date_End.Scan(deleted_date_StartStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	deleted_date_EndStr := c.Query("deleted_date_End")
	deleted_date_End := &pgtype.Timestamp{}
	err = create_date_End.Scan(deleted_date_EndStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	last_date_StartStr := c.Query("last_date_Start")
	last_date_Start := &pgtype.Timestamp{}
	err = create_date_End.Scan(last_date_StartStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	last_date_EndStr := c.Query("last_date_End")
	last_date_End := &pgtype.Timestamp{}
	err = create_date_End.Scan(last_date_EndStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	param := db.GetUserListParams{
		ID:                 id.Int.Int64(),
		Phonenum:           phoneNum,
		Name:               name,
		Belongcmp:          belongCmp,
		CreateDate:         *create_date_Start,
		CreateDate_2:       *create_date_End,
		DeletedDate:        *deleted_date_Start,
		DeletedDate_2:      *deleted_date_End,
		LastModifiedDate:   *last_date_Start,
		LastModifiedDate_2: *last_date_End,
	}
	fmt.Printf("%+v", param)
}

func (u *UserCtrlImpl) GetUserById(c *gin.Context) {

	id := c.Param("id")

	if id != "" {

		id, err := strconv.Atoi(id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		res, err := u.svc.UserServ.GetUserById(int64(id))

		c.JSON(http.StatusOK, gin.H{"res": res})
		return
	}

	c.Status(http.StatusInternalServerError)
	c.Abort()
	return
}

func UserCtrlInit(svc *service.AppService) *UserCtrlImpl {
	return &UserCtrlImpl{
		svc: svc,
	}
}
