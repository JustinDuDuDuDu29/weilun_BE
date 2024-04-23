package controller

type AppControllerImpl struct {
	UserCtrl   UserCtrl
	AuthCtrl   AuthCtrl
	CmpCtrl    CmpCtrl
	JobsCtrl   JobsCtrl
	RepairCtrl RepairCtrl
}

func AppControllerInit(userCtrl UserCtrl, authCtrl AuthCtrl, cmpCtrl CmpCtrl, jobsCtrl JobsCtrl, repairCtrl RepairCtrl) *AppControllerImpl {
	return &AppControllerImpl{
		UserCtrl:   userCtrl,
		AuthCtrl:   authCtrl,
		CmpCtrl:    cmpCtrl,
		JobsCtrl:   jobsCtrl,
		RepairCtrl: repairCtrl,
	}
}
