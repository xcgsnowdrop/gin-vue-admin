package v1

import (
	"gmserver/api/v1/example"
	"gmserver/api/v1/gm"
	"gmserver/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	GmApiGroup      gm.ApiGroup
}
