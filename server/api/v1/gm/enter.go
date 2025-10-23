package gm

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	UserApi
	ItemApi
}

var (
	gmUserService = &service.ServiceGroupApp.GmServiceGroup.UserService
	gmItemService = &service.ServiceGroupApp.GmServiceGroup.ItemService
)
