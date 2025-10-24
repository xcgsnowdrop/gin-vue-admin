package request

import (
	"gmserver/model/common/request"
	"gmserver/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
