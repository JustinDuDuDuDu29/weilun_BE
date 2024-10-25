package controller

import (
	"database/sql"
	"fmt"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"strconv"
	"time"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
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

	res, err := a.svc.UserServ.HaveUser(reqBody.Phonenum)
	fmt.Print("ININ")

	if err != nil {
		c.Status(http.StatusNotFound)
		c.Abort()
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(res.Pwd), []byte(reqBody.Pwd)); err != nil {

		c.Status(http.StatusNotFound)
		c.Abort()
		return
	}

	if res.DeletedDate.Valid {
		c.JSON(http.StatusOK, gin.H{"err": "Account is deleted"})
		return
	}

	rand := randstr.Hex(16)
	var seed sql.NullString
	seed.Scan(rand)

	param := db.NewSeedParams{
		ID:   res.ID,
		Seed: seed,
	}

	err = a.svc.UserServ.NewSeed(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	var newClaim apptypes.CustomClaims
	newClaim.Audience = []string{strconv.Itoa(int(res.ID))}
	newClaim.IssuedAt = jwt.NewNumericDate(time.Now())
	newClaim.NotBefore = jwt.NewNumericDate(time.Now())
	newClaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 1999999))
	newClaim.ID = rand
	newClaim.Issuer = rand

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
