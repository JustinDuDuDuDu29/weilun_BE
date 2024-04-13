package controller

import (
	"fmt"
	"main/apptypes"
	"main/service"
	db "main/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type CmpCtrl interface {
	RegisterCmp(c *gin.Context)
	DeleteCmp(c *gin.Context)
}

type CmpCtrlImpl struct {
	svc *service.AppService
}

func (u *CmpCtrlImpl) RegisterCmp(c *gin.Context) {
	var reqBody apptypes.RegisterCmpT

	if err := c.BindJSON(&reqBody); err != nil {
		return
	}

	newid, err = u.svc.CmpServ.RegisterCmp(param)

	if err != nil {
		fmt.Printf("\n%s", err)

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

	err = u.svc.CmpServ.DeleteCmp(int64(reqBody.ToDeleteUserId))

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
	}

	c.Status(http.StatusOK)
	c.Abort()
	return
}

func CmpCtrlInit(svc *service.AppService) *CmpCtrlImpl {
	return &CmpCtrlImpl{
		svc: svc,
	}
}
