package service

type AppService struct {
	UserServ   UserServ
	CmpServ    CmpServ
	JobsServ   JobsServ
	RepairServ RepairServ
}

func AppServiceInit(userServ UserServ, cmpServ CmpServ, jobsServ JobsServ, repairServ RepairServ) *AppService {
	return &AppService{
		UserServ:   userServ,
		CmpServ:    cmpServ,
		JobsServ:   jobsServ,
		RepairServ: repairServ,
	}
}
