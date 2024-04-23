package controller

import (
	"main/service"

	"github.com/gin-gonic/gin"
)

type RepairCtrl interface {
	CreateNewRepair(c *gin.Context)
}

type RepairCtrlImpl struct {
	svc *service.AppService
}

func (r *RepairCtrlImpl) GetAllRepair(c *gin.Context) {

}

func (r *RepairCtrlImpl) GetRepairByID(c *gin.Context) {

}

func (r *RepairCtrlImpl) GetAllRepairByDriverID(c *gin.Context) {

}

func (r *RepairCtrlImpl) GetAllRepairByCMP(c *gin.Context) {

}

func (r *RepairCtrlImpl) CreateNewRepair(c *gin.Context) {

}

func RepairCtrlInit(svc *service.AppService) *RepairCtrlImpl {
	return &RepairCtrlImpl{
		svc: svc,
	}
}
