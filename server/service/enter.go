package service

import (
	"gmserver/service/example"
	"gmserver/service/gm"
	"gmserver/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	GmServiceGroup      gm.ServiceGroup
}
