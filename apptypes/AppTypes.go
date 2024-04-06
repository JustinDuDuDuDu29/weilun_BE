package apptypes

type DriverInfo struct {
	Percentage       int    `json:"percentage" binding:"omitempty"`
	NationalIdNumber string `json:"nationalIdNumber" binding:"omitempty"`
}

type RegisterUserBodyT struct {
	Name       string     `json:"name" binding:"required"`
	Role       string     `json:"role" binding:"required"`
	PhoneNum   string     `json:"phoneNum" binding:"required"`
	BelongCmp  int        `json:"belongCmp" binding:"required"`
	DriverInfo DriverInfo `json:"driverInfo" validate:"omitempty"`
}

type DeleteUserBodyT struct {
	ToDeleteUserId int `json:"id" binding:"required"`
}
