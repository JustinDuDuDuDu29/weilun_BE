package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *gin.Context) {

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userName": "Justin",
		"exp":      time.Now().Add(time.Second * 20).Unix(),
	})
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("jwtSecret")))

	fmt.Print(err)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jwt": tokenString})
	return
}
