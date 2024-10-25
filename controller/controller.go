package controller

type AppControllerImpl struct {
	UserCtrl    UserCtrl
	AuthCtrl    AuthCtrl
	CmpCtrl     CmpCtrl
	JobsCtrl    JobsCtrl
	RepairCtrl  RepairCtrl
	GasCtrl     GasCtrl
	AlertCtrl   AlertCtrl
	SocketCtrl  SocketCtrl
	RevenueCtrl RevenueCtrl
}

func AppControllerInit(userCtrl UserCtrl, authCtrl AuthCtrl, cmpCtrl CmpCtrl, jobsCtrl JobsCtrl, repairCtrl RepairCtrl, alertCtrl AlertCtrl, socketCtrl SocketCtrl, revenueCtrl RevenueCtrl, gasCtrl GasCtrl) *AppControllerImpl {
	return &AppControllerImpl{
		UserCtrl:    userCtrl,
		AuthCtrl:    authCtrl,
		CmpCtrl:     cmpCtrl,
		JobsCtrl:    jobsCtrl,
		RepairCtrl:  repairCtrl,
		AlertCtrl:   alertCtrl,
		SocketCtrl:  socketCtrl,
		RevenueCtrl: revenueCtrl,
		GasCtrl:     gasCtrl,
	}
}
