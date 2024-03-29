package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type reqBodyT struct {
	Name       string `json:"name" binding:"required"`
	Role       string `json:"role" binding:"required"`
	PhoneNum   string `json:"phoneNum" binding:"required"`
	DriverInfo string `json:"driverInfo"`
}

func RegisterNewUser(c *gin.Context) {
	var reqBody reqBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		return
	}

	fmt.Print(reqBody.Name)
	c.Status(http.StatusOK)
}
