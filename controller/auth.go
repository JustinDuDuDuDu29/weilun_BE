package controller

import (
	"main/apptypes"
	"main/service"
	db "main/sql"
	"time"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thanhpk/randstr"
)

type AuthCtrl interface {
	Login(c *gin.Context)
}

type AuthCtrlImpl struct {
	svc *service.AppService
}

func (a *AuthCtrlImpl) Login(c *gin.Context) {

	var reqBody apptypes.LoginBodyT

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
		c.Abort()
		return
	}
	if res.DeletedDate.Valid {
		c.JSON(http.StatusOK, gin.H{"err": "Account is deleted"})
		return
	}

	var newClaim apptypes.CustomClaims
	newClaim.Audience = []string{"audience-example"}
	newClaim.IssuedAt = jwt.NewNumericDate(time.Now())
	newClaim.NotBefore = jwt.NewNumericDate(time.Now())
	newClaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * 2000))
	newClaim.ID = randstr.Hex(16)
	newClaim.Issuer = "Weilun"

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaim)

	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("accessToken")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Token": tokenString})
}

func AuthCtrlInit(svc *service.AppService) *AuthCtrlImpl {
	return &AuthCtrlImpl{
		svc: svc,
	}
}
