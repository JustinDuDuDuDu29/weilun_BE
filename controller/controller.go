package controller

type AppControllerImpl struct {
	UserCtrl UserCtrl
	AuthCtrl AuthCtrl
	CmpCtrl  CmpCtrl
}

func AppControllerInit(userCtrl UserCtrl, authCtrl AuthCtrl, cmpCtrl CmpCtrl) *AppControllerImpl {
	return &AppControllerImpl{
		UserCtrl: userCtrl,
		AuthCtrl: authCtrl,
		CmpCtrl:  cmpCtrl,
	}
}
