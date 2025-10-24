package initialize

import (
	_ "gmserver/source/example"
	_ "gmserver/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
