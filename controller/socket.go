package controller

import (
	"main/service"

	"github.com/gin-gonic/gin"
)

type SocketCtrl interface {
}

type SocketCtrlImpl struct {
	svc *service.AppService
}

func (s *SocketCtrlImpl) TestSocket(c *gin.Context) {

}

func SocketCtrlInit(svc *service.AppService) *SocketCtrlImpl {
	return &SocketCtrlImpl{
		svc: svc,
	}
}
