package controller

import (
	"fmt"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserCtrl interface {
	RegisterNewUser(c *gin.Context)
}

type UserCtrlImpl struct {
	svc *service.AppService
}

type registerNewUserBodyT struct {
	Name       string `json:"name" binding:"required"`
	Role       string `json:"role" binding:"required"`
	PhoneNum   string `json:"phoneNum" binding:"required"`
	DriverInfo string `json:"driverInfo"`
}

func (u *UserCtrlImpl) RegisterNewUser(c *gin.Context) {
	// var reqBody registerNewUserBodyT
	// err := c.BindJSON(&reqBody)

	fmt.Println("Hello!")
	// if err != nil {
	// 	return
	// }

	c.Status(http.StatusOK)
}

func UserCtrlInit(svc *service.AppService) *UserCtrlImpl {
	return &UserCtrlImpl{
		svc: svc,
	}
}
