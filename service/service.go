package service

type AppService struct {
	UserServ    UserServ
	CmpServ     CmpServ
	JobsServ    JobsServ
	RepairServ  RepairServ
	GasServ     GasServ
	AlertServ   AlertServ
	RevenueServ RevenueServ
}

func AppServiceInit(userServ UserServ, cmpServ CmpServ, jobsServ JobsServ, repairServ RepairServ, alertServ AlertServ, revenueServ RevenueServ, gasServ GasServ) *AppService {
	return &AppService{
		UserServ:    userServ,
		CmpServ:     cmpServ,
		JobsServ:    jobsServ,
		RepairServ:  repairServ,
		GasServ:     gasServ,
		AlertServ:   alertServ,
		RevenueServ: revenueServ,
	}
}
