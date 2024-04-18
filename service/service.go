package service

type AppService struct {
	UserServ UserServ
	CmpServ  CmpServ
	JobsServ JobsServ
}

func AppServiceInit(userServ UserServ, cmpServ CmpServ, jobsServ JobsServ) *AppService {
	return &AppService{
		UserServ: userServ,
		CmpServ:  cmpServ,
		JobsServ: jobsServ,
	}
}
