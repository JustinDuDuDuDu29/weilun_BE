package apptypes

import (
	// "encoding/json"
	db "main/sql"
	"mime/multipart"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
}

type Excel struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	List     []struct {
		Gas  int `json:"gas"`
		Data []struct {
			Jp       int    `json:"jp"`
			Ss       string `json:"ss"`
			Mid      string `json:"mid"`
			Count    int    `json:"count"`
			Toloc    string `json:"toloc"`
			Total    int    `json:"total"`
			CmpName  string `json:"cmpName"`
			FromLoc  string `json:"fromLoc"`
			Platenum string `json:"platenum"`
		}
		Date string `json:"date"`
	} `json:"list"`
}

type ApproveJob struct {
	Memo string `json:"memo"`
}
type GetJobsClientBodyT struct {
	ID      int    `json:"id"`
	FromLoc string `json:"fromLoc"`
	Mid     string `json:"mid"`
	ToLoc   string `json:"toLoc"`
}
type GetJobsBodyT struct {
	ID                    int    `json:"id"`
	FromLoc               string `json:"fromLoc"`
	Mid                   string `json:"mid"`
	ToLoc                 string `json:"toLoc"`
	Belongcmp             int    `json:"belongCmp"`
	Remaining             int    `json:"remaining"`
	CloseDateStart        string `json:"closeDateStart"`
	CloseDateEnd          string `json:"closeDateEnd"`
	CreateDateStart       string `json:"createDateStart"`
	CreateDateEnd         string `json:"createDateEnd"`
	DeletedDateStart      string `json:"deletedDateStart"`
	DeletedDateEnd        string `json:"deletedDateEnd"`
	LastModifiedDateStart string `json:"lastModifiedDateStart"`
	LastModifiedDateEnd   string `json:"lastModifiedDateEnd"`
}

type UpdateDriverPic struct {
	Insurances    *multipart.FileHeader `form:"Insurances"`
	Registration  *multipart.FileHeader `form:"Registration"`
	DriverLicense *multipart.FileHeader `form:"DriverLicense"`
	TruckLicense  *multipart.FileHeader `form:"TruckLicense"`
}

type FinishClaimJobBodyT struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UpdatePasswordBodyT struct {
	Id     int    `json:"id" binding:"required"`
	Pwd    string `json:"pwd" binding:"required"`
	OldPwd string `json:"oldPwd"`
}
type ResetPasswordBodyT struct {
	Id int `json:"id" binding:"required"`
}

type UpdateJobBodyT struct {
	FromLoc   string `json:"fromLoc" binding:"required"`
	Mid       string `json:"mid"`
	ToLoc     string `json:"toLoc" binding:"required"`
	Price     int    `json:"price" binding:"required"`
	Remaining int    `json:"remaining" binding:""`
	Belongcmp int    `json:"belongCmp" binding:"required"`
	Source    string `json:"source" binding:"required"`
	// Jobdate   string `json:"jobDate" binding:"required"`
	Memo string `json:"memo"`
	// CloseDate string `json:"closeDate"`
	ID int `json:"id" binding:"required"`
}

type NewRepairBodyT struct {
	Place      string                         `form:"place" binding:"required"`
	Repairinfo []db.CreateNewRepairInfoParams `form:"repairInfo" binding:"required"`
	RepairPic  *multipart.FileHeader          `form:"repairPic"`
}

type NewGasBodyT struct {
	Place   string                      `form:"place" binding:"required"`
	Gasinfo []db.CreateNewGasInfoParams `form:"gasInfo" binding:"required"`
	GasPic  *multipart.FileHeader       `form:"gasPic"`
}

type CreateAlertBodyT struct {
	Alert     string `json:"alert" binding:"required"`
	BelongCmp int    `json:"belongCmp" binding:"required"`
}

type CreateJobBodyT struct {
	FromLoc string `json:"fromLoc" binding:"required"`
	Mid     string `json:"mid"`
	ToLoc   string `json:"toLoc" binding:"required"`
	Price   int    `json:"price" binding:"required"`
	// Estimated int    `json:"estimated" binding:"required"`
	Estimated int    `json:"estimated"`
	Belongcmp int    `json:"belongCmp" binding:"required"`
	Source    string `json:"source" binding:"required"`
	// Jobdate   string `json:"jobDate" binding:"required"`
	Memo string `json:"memo"`
	// CloseDate string `json:"closeDate"`
}

type DriverInfo struct {
	PlateNum string `json:"plateNum" binding:"required"`
	// Percentage       int    `json:"percentage" binding:"required"`
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
