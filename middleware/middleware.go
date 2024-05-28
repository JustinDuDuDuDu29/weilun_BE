package middleware

type AppMiddlewareImpl struct {
	RoleMid RoleMid
}

func AppMiddlewareInit(roleMid RoleMid) *AppMiddlewareImpl {
	return &AppMiddlewareImpl{
		RoleMid: roleMid,
	}
}
