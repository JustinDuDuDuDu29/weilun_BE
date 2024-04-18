package service

type AppService struct {
	UserServ UserServ
	CmpServ  CmpServ
}

func AppServiceInit(userServ UserServ, cmpServ CmpServ) *AppService {
	return &AppService{
		UserServ: userServ,
		CmpServ:  cmpServ,
	}
}
