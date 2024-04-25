package apptypes

import (
	"encoding/json"
	"mime/multipart"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
}

type FinishClaimJobBodyT struct {
	File *multipart.FileHeader `form:"file"`
}

type UpdateJobBodyT struct {
	FromLoc   string `json:"fromLoc" binding:"required"`
	Mid       string `json:"mid"`
	ToLoc     string `json:"toLoc" binding:"required"`
	Price     int    `json:"price" binding:"required"`
	Remaining int    `json:"estimated" binding:"required"`
	Belongcmp int    `json:"belongCmp" binding:"required"`
	Source    string `json:"source" binding:"required"`
	// Jobdate   time.Time `json:"jobDate" binding:"required"`
	Jobdate   string `json:"jobDate" binding:"required"`
	Memo      string `json:"memo"`
	CloseDate string `json:"closeDate" binding:"required"`
	ID        int    `json:"id" binding:"required"`
}

type NewRepairBodyT struct {
	Repairinfo json.RawMessage `json:"repairInfo" binding:"required"`
}

type CreateAlertBodyT struct {
	Alert     string `json:"alert" binding:"required"`
	BelongCmp int    `json:"belongCmp" binding:"required"`
}

type CreateJobBodyT struct {
	FromLoc   string `json:"fromLoc" binding:"required"`
	Mid       string `json:"mid"`
	ToLoc     string `json:"toLoc" binding:"required"`
	Price     int    `json:"price" binding:"required"`
	Estimated int    `json:"estimated" binding:"required"`
	Belongcmp int    `json:"belongCmp" binding:"required"`
	Source    string `json:"source" binding:"required"`
	Jobdate   string `json:"jobDate" binding:"required"`
	Memo      string `json:"memo"`
	CloseDate string `json:"closeDate"`
}

type DriverInfo struct {
	Percentage       int    `json:"percentage" binding:"required"`
	NationalIdNumber string `json:"nationalIdNumber" binding:"required"`
}

type LoginBodyT struct {
	Phonenum string `json:"phoneNum" binding:"required"`
	Pwd      string `json:"pwd" binding:"required"`
}

type RegisterCmpAdminBodyT struct {
	Name      string `json:"name" binding:"required"`
	Role      string `json:"role" binding:"required"`
	PhoneNum  string `json:"phoneNum" binding:"required"`
	BelongCmp int    `json:"belongCmp" binding:"required"`
}

type RegisterDriverBodyT struct {
	Name       string     `json:"name" binding:"required"`
	Role       string     `json:"role" binding:"required"`
	PhoneNum   string     `json:"phoneNum" binding:"required"`
	BelongCmp  int        `json:"belongCmp" binding:"required"`
	DriverInfo DriverInfo `json:"driverInfo" validate:"required"`
}

type DeleteUserBodyT struct {
	ToDeleteUserId int `json:"id" binding:"required"`
}

type GetUserT struct {
	Id int `json:"id" binding:"required"`
}

type GetCmpT struct {
	Id int `json:"id" binding:"required"`
}

type RegisterCmpT struct {
	CmpName string `json:"cmpName" binding:"required"`
}

type UpdateCmpT struct {
	Id      int    `json:"id" binding:"required"`
	CmpName string `json:"cmpName" binding:"required"`
}

type DeleteCmpBodyT struct {
	ToDeleteCmpId int `json:"id" binding:"required"`
}
