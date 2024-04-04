package controller

type AppControllerImpl struct {
	UserCtrl UserCtrl
	AuthCtrl AuthCtrl
}

func AppControllerInit(userCtrl UserCtrl, authCtrl AuthCtrl) *AppControllerImpl {
	return &AppControllerImpl{
		UserCtrl: userCtrl,
		AuthCtrl: authCtrl,
	}
}
