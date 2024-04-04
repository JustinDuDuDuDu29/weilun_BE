package controller

import (
	"fmt"
	"main/service"
	db "main/sql"

	// db "main/sql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthCtrl interface {
	Login(c *gin.Context)
}

type AuthCtrlImpl struct {
	svc *service.AppService
}

type loginBodyT struct {
	Phonenum string `json:"phoneNum" binding:"required"`
	Pwd      string `json:"pwd" binding:"required"`
}

func (a *AuthCtrlImpl) Login(c *gin.Context) {

	var reqBody loginBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		return
	}

	res, err := a.svc.UserServ.HaveUser(db.GetUserParams{
		Phonenum: reqBody.Phonenum,
		Pwd:      reqBody.Pwd,
	})

	if err != nil {
		c.Status(http.StatusNotFound)
	}

	fmt.Printf("%+v", res)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":  res.ID,
		"exp": time.Now().Add(time.Second * 20).Unix(),
	})

	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("accessToken")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": tokenString})

	// return
}

func AuthCtrlInit(svc *service.AppService) *AuthCtrlImpl {
	return &AuthCtrlImpl{
		svc: svc,
	}
}
