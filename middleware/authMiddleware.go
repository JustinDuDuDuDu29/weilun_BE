package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type reqBodyT struct {
	// must start with uppercase or it wont export!!!
	// https://github.com/gin-gonic/gin/issues/424
	Token string `json:"Token" binding:"required"`
}

func IsLoggedIn(c *gin.Context) {

	var reqBody reqBodyT
	err := c.ShouldBindJSON(&reqBody)

	if err != nil {
		c.Abort()
		return
	}

	fmt.Print("res: ", reqBody)
	token, err := jwt.Parse(reqBody.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("jwtSecret")), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Token has expired!!!"})
			c.Abort()
			return

		default:
			fmt.Print("err:", err)
			c.JSON(http.StatusBadRequest, gin.H{"err": "grow up, K? get a real job or something..."})
			c.Abort()
			return
		}

	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("claims:", claims)
	} else {
		fmt.Println("err:", err)
	}
	c.Next()

}
