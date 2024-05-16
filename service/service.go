package service

type AppService struct {
	UserServ    UserServ
	CmpServ     CmpServ
	JobsServ    JobsServ
	RepairServ  RepairServ
	AlertServ   AlertServ
	RevenueServ RevenueServ
}

func AppServiceInit(userServ UserServ, cmpServ CmpServ, jobsServ JobsServ, repairServ RepairServ, alertServ AlertServ, revenueServ RevenueServ) *AppService {
	return &AppService{
		UserServ:    userServ,
		CmpServ:     cmpServ,
		JobsServ:    jobsServ,
		RepairServ:  repairServ,
		AlertServ:   alertServ,
		RevenueServ: revenueServ,
	}
}
