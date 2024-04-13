package controller

import (
	"fmt"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserCtrl interface {
	RegisterUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type UserCtrlImpl struct {
	svc *service.AppService
}

func (u *UserCtrlImpl) RegisterUser(c *gin.Context) {
	var reqBody apptypes.RegisterUserBodyT

	if err := c.BindJSON(&reqBody); err != nil {
		return
	}

	var newid int64
	var err error

	switch reqBody.Role {
	case "cmpAdmin":
		param := db.CreateUserParams{
			Pwd:       reqBody.PhoneNum,
			Name:      reqBody.Name,
			Belongcmp: pgtype.Int8{Int64: int64(reqBody.BelongCmp), Valid: true},
			Phonenum:  reqBody.PhoneNum,
			Role:      200,
		}
		newid, err = u.svc.UserServ.RegisterCmpAdmin(param)

	case "driver":
		param := db.CreateUserParams{
			Pwd:       reqBody.PhoneNum,
			Name:      reqBody.Name,
			Belongcmp: pgtype.Int8{Int64: int64(reqBody.BelongCmp), Valid: true},
			Phonenum:  reqBody.PhoneNum,
			Role:      300,
		}

		newid, err = u.svc.UserServ.RegisterDriver(param, reqBody.DriverInfo.Percentage, reqBody.DriverInfo.NationalIdNumber)

	default:
		fmt.Print("out")
	}

	if err != nil {
		fmt.Print("HERE")
		fmt.Printf("\n%s", err)

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

func UserCtrlInit(svc *service.AppService) *UserCtrlImpl {
	return &UserCtrlImpl{
		svc: svc,
	}
}
