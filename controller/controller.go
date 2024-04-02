package controller

type AppControllerImpl struct {
	UserCtrl UserCtrl
}

func AppControllerInit(userCtrl UserCtrl) *AppControllerImpl {
	return &AppControllerImpl{
		UserCtrl: userCtrl,
	}
}
