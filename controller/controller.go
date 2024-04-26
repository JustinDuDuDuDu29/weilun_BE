package controller

type AppControllerImpl struct {
	UserCtrl   UserCtrl
	AuthCtrl   AuthCtrl
	CmpCtrl    CmpCtrl
	JobsCtrl   JobsCtrl
	RepairCtrl RepairCtrl
	AlertCtrl  AlertCtrl
	SocketCtrl SocketCtrl
}

func AppControllerInit(userCtrl UserCtrl, authCtrl AuthCtrl, cmpCtrl CmpCtrl, jobsCtrl JobsCtrl, repairCtrl RepairCtrl, alertCtrl AlertCtrl, socketCtrl SocketCtrl) *AppControllerImpl {
	return &AppControllerImpl{
		UserCtrl:   userCtrl,
		AuthCtrl:   authCtrl,
		CmpCtrl:    cmpCtrl,
		JobsCtrl:   jobsCtrl,
		RepairCtrl: repairCtrl,
		AlertCtrl:  alertCtrl,
		SocketCtrl: socketCtrl,
	}
}
