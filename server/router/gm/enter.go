package gm

import (
	api "gmserver/api/v1"
)

type RouterGroup struct {
	UserRouter
	ItemRouter
}

var (
	gmUserApi = api.ApiGroupApp.GmApiGroup.UserApi
	gmItemApi = api.ApiGroupApp.GmApiGroup.ItemApi
)
