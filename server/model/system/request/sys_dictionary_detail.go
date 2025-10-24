package request

import (
	"gmserver/model/common/request"
	"gmserver/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
