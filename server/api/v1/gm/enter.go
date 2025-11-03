package gm

import "gmserver/service"

type ApiGroup struct {
	UserApi
	ItemApi
	EmailAuditApi
}

var (
	gmUserService = &service.ServiceGroupApp.GmServiceGroup.UserService
	gmItemService = &service.ServiceGroupApp.GmServiceGroup.ItemService
)
