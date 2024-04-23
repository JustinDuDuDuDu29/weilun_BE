package middleware

import (
	"fmt"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleMid interface {
	SuperAdminOnly(c *gin.Context)
}

type RoleMidImpl struct {
	svc *service.AppService
}

func (m *RoleMidImpl) SuperAdminOnly(c *gin.Context) {
	// var uid
	fmt.Println(c.MustGet("UserID"))
	res, err := m.svc.UserServ.GetUserById(int64(c.MustGet("UserID").(int)))

	if err != nil || res.Role > 100 {
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
