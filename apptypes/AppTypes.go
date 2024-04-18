package apptypes

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
