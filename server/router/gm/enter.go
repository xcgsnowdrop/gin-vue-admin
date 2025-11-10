package gm

import (
	api "gmserver/api/v1"
)

type RouterGroup struct {
	GameApiProxyRouter
	EmailAuditRouter
}

var (
	gmEmailAuditApi = api.ApiGroupApp.GmApiGroup.EmailAuditApi
)
