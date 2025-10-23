package gm

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	UserRouter
	ItemRouter
}

var (
	gmUserApi = api.ApiGroupApp.GmApiGroup.UserApi
	gmItemApi = api.ApiGroupApp.GmApiGroup.ItemApi
)
