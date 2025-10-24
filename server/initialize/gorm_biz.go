package initialize

import (
	"gmserver/global"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate()
	if err != nil {
		return err
	}
	return nil
}
