package middleware

import (
	"errors"
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
	err := c.BindJSON(&reqBody)

	if err != nil {
		c.Abort()
		return
	}

	token, err := jwt.Parse(reqBody.Token, func(token *jwt.Token) (interface{}, error) {
		// if _, ok := token.Method .(*jwt.SigningMethodHMAC); !ok {
		// 	c.Abort()
		// 	return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		// }
		return []byte(os.Getenv("accessToken")), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Token has expired!!!"})
			c.Abort()
			return

		default:
			c.JSON(http.StatusBadRequest, gin.H{"err": "grow up, K? get a real job or something..."})
			c.Abort()
			return
		}

	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		res, err := claims.GetAudience()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "grow up, K? get a real job or something..."})
			c.Abort()
			return
		}
		c.Set("UserID", res)
		c.Next()
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"err": "grow up, K? get a real job or something..."})
		c.Abort()
		return
	}

}
