package service

type AppService struct {
	UserServ UserServ
}

func AppServiceInit(userServ UserServ) *AppService {
	return &AppService{
		UserServ: userServ,
	}
}
