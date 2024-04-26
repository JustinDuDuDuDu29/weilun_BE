package controller

import (
	"fmt"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CmpCtrl interface {
	RegisterCmp(c *gin.Context)
	DeleteCmp(c *gin.Context)
	UpdateCmp(c *gin.Context)
	GetCmp(c *gin.Context)
	GetAllCmp(c *gin.Context)
}

type CmpCtrlImpl struct {
	svc *service.AppService
}

func (u *CmpCtrlImpl) GetCmp(c *gin.Context) {

	var reqBody apptypes.GetCmpT

	err := c.BindJSON(&reqBody)

	if err != nil {
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	data, err := u.svc.CmpServ.GetCmp(int64(reqBody.Id))

	if err != nil {
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, data)
}

func (u *CmpCtrlImpl) GetAllCmp(c *gin.Context) {

	cmpList, err := u.svc.CmpServ.GetAllCmp()

	if err != nil {
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, cmpList)
}

func (u *CmpCtrlImpl) UpdateCmp(c *gin.Context) {
	var reqBody apptypes.UpdateCmpT

	if err := c.BindJSON(&reqBody); err != nil {
		return
	}

	param := db.UpdateCmpParams{
		ID:   int64(reqBody.Id),
		Name: reqBody.CmpName,
	}

	err := u.svc.CmpServ.UpdateCmp(param)

	if err != nil {
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (u *CmpCtrlImpl) RegisterCmp(c *gin.Context) {
	var reqBody apptypes.RegisterCmpT

	if err := c.BindJSON(&reqBody); err != nil {
		return
	}

	newid, err := u.svc.CmpServ.NewCmp(reqBody.CmpName)

	if err != nil {
		fmt.Print(err)
		c.Status(http.StatusConflict)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": newid})
}

func (u *CmpCtrlImpl) DeleteCmp(c *gin.Context) {
	var reqBody apptypes.DeleteCmpBodyT
	err := c.BindJSON(&reqBody)

	if err != nil {
		return
	}

	err = u.svc.CmpServ.DeleteCmp(int64(reqBody.ToDeleteCmpId))

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
	}

	c.Status(http.StatusOK)
	c.Abort()
}

func CmpCtrlInit(svc *service.AppService) *CmpCtrlImpl {
	return &CmpCtrlImpl{
		svc: svc,
	}
}
