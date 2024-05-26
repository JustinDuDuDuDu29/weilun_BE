package middleware

import (
	"errors"
	"fmt"
	"main/service"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RoleMid interface {
	SuperAdminOnly(c *gin.Context)
	IsLoggedIn(c *gin.Context)
	DriverOnly(c *gin.Context)
}

type RoleMidImpl struct {
	svc *service.AppService
}

func (m *RoleMidImpl) IsLoggedIn(c *gin.Context) {
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
			c.JSON(http.StatusUnauthorized, gin.H{"err": "grow up, K? get a real job or something..."})
			c.Abort()
			return
		}

	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		res, err := claims.GetAudience()

		if err != nil {

			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "grow up, K? get a real job or something..."})
			c.Abort()
			return
		}

		id, err := strconv.Atoi(res[0])
		if err != nil {

			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "grow up, K? get a real job or something..."})
			c.Abort()
			return
		}
		c.Set("UserID", id)
		info, err := m.svc.UserServ.GetSeed(int64(id))
		if err != nil {

			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "grow up, K? get a real job or something..."})
			c.Abort()
			return
		}

		issuer, err := claims.GetIssuer()
		if err != nil {

			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "grow up, K? get a real job or something..."})
			c.Abort()
			return
		}
		if info.String != issuer {
			c.JSON(http.StatusUnavailableForLegalReasons, gin.H{"err": "Revalid"})
			c.Abort()
			return
		}
		userInfo, err := m.svc.UserServ.GetUserById(int64(id))
		fmt.Println(err)

		if err != nil {

			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "grow up, K? get a real job or something..."})
			c.Abort()
			return
		}
		c.Set("Role", userInfo.Role)
		c.Set("belongCmp", userInfo.Belongcmp)
		c.Next()

		return

	} else {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"err": "grow up, K? get a real job or something..."})
		c.Abort()
		return
	}

}

func (m *RoleMidImpl) SuperAdminOnly(c *gin.Context) {
	res := c.MustGet("Role").(int16)
	if (res) > 100 {
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
	c.Next()
}

func (m *RoleMidImpl) DriverOnly(c *gin.Context) {
	res := c.MustGet("Role").(int16)
	if res < 300 {
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
	c.Next()
}

func RoleMidInit(svc *service.AppService) *RoleMidImpl {
	return &RoleMidImpl{
		svc: svc,
	}
}
