package controller

import (
	"main/service"
)

type SocketCtrl interface {
}

type SocketCtrlImpl struct {
	svc *service.AppService
}

func SocketCtrlInit(svc *service.AppService) *SocketCtrlImpl {
	return &SocketCtrlImpl{
		svc: svc,
	}
}
