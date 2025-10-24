package router

import (
	"gmserver/router/example"
	"gmserver/router/gm"
	"gmserver/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	GM      gm.RouterGroup
}
