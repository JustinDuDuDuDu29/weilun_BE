package controller

type AppControllerImpl struct {
	UserCtrl UserCtrl
	AuthCtrl AuthCtrl
	CmpCtrl  CmpCtrl
	JobsCtrl JobsCtrl
}

func AppControllerInit(userCtrl UserCtrl, authCtrl AuthCtrl, cmpCtrl CmpCtrl, jobsCtrl JobsCtrl) *AppControllerImpl {
	return &AppControllerImpl{
		UserCtrl: userCtrl,
		AuthCtrl: authCtrl,
		CmpCtrl:  cmpCtrl,
		JobsCtrl: jobsCtrl,
	}
}
