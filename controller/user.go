package controller

import (
	"fmt"
	"main/service"
	db "main/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserCtrl interface {
	RegisterNewUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type UserCtrlImpl struct {
	svc *service.AppService
}

type registerNewUserBodyT struct {
	Name       string `json:"name" binding:"required"`
	Role       string `json:"role" binding:"required"`
	PhoneNum   string `json:"phoneNum" binding:"required"`
	BelongCmp  int    `json:"belongCmp" binding:"required"`
	DriverInfo string `json:"driverInfo"`
}

func (u *UserCtrlImpl) RegisterNewUser(c *gin.Context) {
	var reqBody registerNewUserBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		return
	}

	switch reqBody.Role {
	case "admin":
		// return dont 87
		return
	case "cmpAdmin":
		fmt.Print("11")
		param := db.CreateCmpAdminParams{
			Pwd:       reqBody.PhoneNum,
			Name:      reqBody.Name,
			Belongcmp: pgtype.Int8{Int64: int64(reqBody.BelongCmp), Valid: true},
			Phonenum:  reqBody.PhoneNum,
			Role:      200,
		}
		u.svc.UserServ.RegisterCmpAdmin(param)
	case "driver":

	default:

	}

	c.Status(http.StatusOK)
}

type deleteUserBodyT struct {
	ToDeleteUserId int `json:"id" binding:"required"`
}

func (u *UserCtrlImpl) DeleteUser(c *gin.Context) {
	var reqBody deleteUserBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		return
	}

	u.svc.UserServ.DeleteUser(int64(reqBody.ToDeleteUserId))

	c.Status(http.StatusOK)
}

func UserCtrlInit(svc *service.AppService) *UserCtrlImpl {
	return &UserCtrlImpl{
		svc: svc,
	}
}
