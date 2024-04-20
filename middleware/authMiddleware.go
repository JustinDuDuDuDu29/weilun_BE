package middleware

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// type reqBodyT struct {
// 	// must start with uppercase or it wont export!!!
// 	// https://github.com/gin-gonic/gin/issues/424
// 	Token string `json:"Token" binding:"required"`
// }

func IsLoggedIn(c *gin.Context) {

	// var reqBody reqBodyT
	// err := c.BindJSON(&reqBody)

	// if err != nil {
	// c.Abort()
	// return
	// }
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Authorization is null in Header",
		})
		c.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Format of Authorization is wrong",
		})
		c.Abort()
		return
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
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
		id, err := strconv.Atoi(res[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "grow up, K? get a real job or something..."})
			c.Abort()
			return
		}
		c.Set("UserID", id)
		c.Next()
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"err": "grow up, K? get a real job or something..."})
		c.Abort()
		return
	}

}
