package gm

import "gmserver/service"

type ApiGroup struct {
	UserApi
	ItemApi
}

var (
	gmUserService = &service.ServiceGroupApp.GmServiceGroup.UserService
	gmItemService = &service.ServiceGroupApp.GmServiceGroup.ItemService
)
