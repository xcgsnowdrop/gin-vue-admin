package request

import (
	"gmserver/global"
	"gmserver/model/system"
)

// AddMenuAuthorityInfo Add menu authority info structure
type AddMenuAuthorityInfo struct {
	Menus       []system.SysBaseMenu `json:"menus"`
	AuthorityId uint                 `json:"authorityId"` // 角色ID
}

// DefaultMenu 返回默认菜单列表（dashboard）
func DefaultMenu() []system.SysBaseMenu {
	return []system.SysBaseMenu{{
		GVA_MODEL: global.GVA_MODEL{ID: 1}, // 注意：dashboard的在sys_base_menus表中的ID未必是1，后续需要通过Name或Path从数据库中查询实际的菜单记录
		ParentId:  0,
		Path:      "dashboard",
		Name:      "dashboard",
		Component: "view/dashboard/index.vue",
		Sort:      1,
		Meta: system.Meta{
			Title: "仪表盘",
			Icon:  "setting",
		},
	}}
}
