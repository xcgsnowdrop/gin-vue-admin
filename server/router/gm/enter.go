package gm

import (
	api "gmserver/api/v1"
)

type RouterGroup struct {
	UserRouter
	ItemRouter
	GameApiProxyRouter
	EmailAuditRouter
}

var (
	gmUserApi       = api.ApiGroupApp.GmApiGroup.UserApi
	gmItemApi       = api.ApiGroupApp.GmApiGroup.ItemApi
	gmEmailAuditApi = api.ApiGroupApp.GmApiGroup.EmailAuditApi
)
